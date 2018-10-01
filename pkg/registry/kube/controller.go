package kube

import (
	"fmt"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"time"

	log "github.com/sirupsen/logrus"
	clientset "github.frg.tech/cloud/fanplane/pkg/apis/client/clientset/versioned"
	fanplaneScheme "github.frg.tech/cloud/fanplane/pkg/apis/client/clientset/versioned/scheme"
	informers "github.frg.tech/cloud/fanplane/pkg/apis/client/informers/externalversions/fanplane/v1alpha1"
	listers "github.frg.tech/cloud/fanplane/pkg/apis/client/listers/fanplane/v1alpha1"
	fanplane "github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	kubeCache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "fanplane-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a EnvoyBootstrap is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a EnvoyBootstrap fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by EnvoyBootstrap"
	// MessageResourceSynced is the message used for an Event fired when a EnvoyBootstrap
	// is synced successfully
	MessageResourceSynced = "EnvoyBootstrap synced successfully"
)

// Controller is the controller implementation for Fanplane resources
type Controller struct {
	//cache is the go-control-plane cache used to hydrate envoy sidecars
	cache cache.Cache

	//version is an atomic uint32 used to generate unique version identifiers for snapCache
	version uint32

	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface

	// fanplaneClientSet is a clientset for our own API group
	fanplaneClientSet clientset.Interface

	deploymentsLister    listers.GatewayLister
	deploymentsSynced    kubeCache.InformerSynced
	gatewayLister        listers.GatewayLister
	gatewaySynced        kubeCache.InformerSynced
	envoyBootstrapLister listers.EnvoyBootstrapLister
	envoySynced          kubeCache.InformerSynced

	// these are rate limited work queues. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	envoyBtWorkqueue workqueue.RateLimitingInterface
	gatewayWorkqueue workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new Fanplane crd controller
func NewController(
	snapCache cache.Cache,
	fanplaneClientSet clientset.Interface,
	gatewayInformer informers.GatewayInformer,
	envoyBootstrapInformer informers.EnvoyBootstrapInformer) *Controller {

	// Create event broadcaster
	// Add fanplane-controller types to the default Kubernetes Scheme so Events can be
	// logged for fanplane-controller types.
	fanplaneScheme.AddToScheme(scheme.Scheme)

	log.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(log.Infof)
	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		cache:                snapCache,
		fanplaneClientSet:    fanplaneClientSet,
		gatewayLister:        gatewayInformer.Lister(),
		gatewaySynced:        gatewayInformer.Informer().HasSynced,
		envoyBootstrapLister: envoyBootstrapInformer.Lister(),
		envoySynced:          envoyBootstrapInformer.Informer().HasSynced,
		envoyBtWorkqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EnvoyBootstrap"),
		gatewayWorkqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Gateway"),
		recorder:             recorder,
	}

	log.Info("Setting up event handlers")

	// Set up an event handler for when EnvoyBootstrap resources change
	envoyBootstrapInformer.Informer().AddEventHandler(kubeCache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueEnvoyBootstrapObj,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueEnvoyBootstrapObj(new)
		},
		DeleteFunc: controller.removeConfigurationFromCache,
	})

	// Set up an event handler for when Gateway resources change
	gatewayInformer.Informer().AddEventHandler(kubeCache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueGatewayObj,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueGatewayObj(new)
		},
		DeleteFunc: controller.removeConfigurationFromCache,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the envoyBtWorkqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.envoyBtWorkqueue.ShutDown()
	defer c.gatewayWorkqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting EnvoyBootstrap controller")

	// Wait for the caches to be synced before starting workers
	log.Info("Waiting for informer caches to sync")
	if ok := kubeCache.WaitForCacheSync(stopCh, c.gatewaySynced, c.envoySynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Info("Starting workers")

	// Launch two workers to process Fanplane resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	// Launch two workers to process Fanplane resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runGatewayWorker, time.Second, stopCh)
	}

	log.Info("Started workers")
	<-stopCh
	log.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// envoyBtWorkqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) runGatewayWorker() {
	for c.processNextGatewayWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the envoyBtWorkqueue and
// attempt to process it, by calling the syncHandlerEnvoyBootstrap.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.envoyBtWorkqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.envoyBtWorkqueue.Done.
	err := func(obj interface{}) error {
		// call Done here so the envoyBtWorkqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the envoyBtWorkqueue and attempted again after a back-off
		// period.
		defer c.envoyBtWorkqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the envoyBtWorkqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// envoyBtWorkqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// envoyBtWorkqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the envoyBtWorkqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.envoyBtWorkqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in envoyBtWorkqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandlerEnvoyBootstrap, passing it the namespace/name string of the
		// EnvoyBootstrap resource to be synced.
		if err := c.syncHandlerEnvoyBootstrap(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.envoyBtWorkqueue.Forget(obj)
		log.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
	}

	return true
}

// processNextWorkItem will read a single work item off the envoyBtWorkqueue and
// attempt to process it, by calling the syncHandlerEnvoyBootstrap.
func (c *Controller) processNextGatewayWorkItem() bool {
	obj, shutdown := c.gatewayWorkqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.envoyBtWorkqueue.Done.
	err := func(obj interface{}) error {
		// call Done here so the envoyBtWorkqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the envoyBtWorkqueue and attempted again after a back-off
		// period.
		defer c.gatewayWorkqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the envoyBtWorkqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// envoyBtWorkqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// envoyBtWorkqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the envoyBtWorkqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.gatewayWorkqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in gatewayWorkqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandlerEnvoyBootstrap, passing it the namespace/name string of the
		// EnvoyBootstrap resource to be synced.
		if err := c.syncHandlerGateway(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.gatewayWorkqueue.Forget(obj)
		log.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
	}

	return true
}

// syncHandlerEnvoyBootstrap compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the EnvoyBootstrap resource
// with the current status of the resource.
func (c *Controller) syncHandlerEnvoyBootstrap(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := kubeCache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the EnvoyBootstrap resource with this namespace/name
	envoyBootstrap, err := c.envoyBootstrapLister.EnvoyBootstraps(namespace).Get(name)
	if err != nil {
		// The EnvoyBootstrap resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("envoyBootstrap '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	err = c.cache.Add(envoyBootstrap)

	if err != nil {
		return err
	}

	// Finally, we update the status block of the EnvoyBootstrap resource to reflect the
	// current state of the world
	err = c.updateEnvoyBootstrapStatus(envoyBootstrap)
	if err != nil {
		return err
	}

	c.recorder.Event(envoyBootstrap, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// syncHandlerEnvoyBootstrap compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the EnvoyBootstrap resource
// with the current status of the resource.
func (c *Controller) syncHandlerGateway(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := kubeCache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Gateway resource with this namespace/name
	gateway, err := c.gatewayLister.Gateways(namespace).Get(name)
	if err != nil {
		// The EnvoyBootstrap resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("gateway '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	err = c.cache.Add(gateway)

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the EnvoyBootstrap resource to reflect the
	// current state of the world
	err = c.updateGatewayStatus(gateway)
	if err != nil {
		return err
	}

	c.recorder.Event(gateway, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateEnvoyBootstrapStatus(envoyBootstrap *fanplane.EnvoyBootstrap) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	envoyBootstrapCopy := envoyBootstrap.DeepCopy()
	envoyBootstrapCopy.Status.Processed = true

	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the EnvoyBootstrap resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.fanplaneClientSet.FanplaneV1alpha1().EnvoyBootstraps(envoyBootstrap.Namespace).Update(envoyBootstrapCopy)
	return err
}

func (c *Controller) updateGatewayStatus(gateway *fanplane.Gateway) error {
	gatewayCopy := gateway.DeepCopy()
	gatewayCopy.Status.Processed = true

	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the EnvoyBootstrap resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.fanplaneClientSet.FanplaneV1alpha1().Gateways(gateway.Namespace).Update(gatewayCopy)
	return err
}

// enqueueEnvoyBootstrapObj takes a EnvoyBootstrap resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than EnvoyBootstrap.
func (c *Controller) enqueueEnvoyBootstrapObj(obj interface{}) {
	var key string
	var err error

	if key, err = kubeCache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.envoyBtWorkqueue.AddRateLimited(key)
}

// enqueueGatewayObj takes a Gateway resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Gateway.
func (c *Controller) enqueueGatewayObj(obj interface{}) {
	var key string
	var err error

	if key, err = kubeCache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.gatewayWorkqueue.AddRateLimited(key)
}

//removeConfigurationFromCache parses object header and removes it from server cache
func (c *Controller) removeConfigurationFromCache(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); ok {
		c.cache.RemoveById(object.GetName())
	} else {
		log.Errorf("couldn't remove obj %s from cache", obj)
	}
}
