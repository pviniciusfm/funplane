package kube

import (
	"fmt"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	core "k8s.io/client-go/testing"
	kubeCache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	"github.frg.tech/cloud/fanplane/pkg/apis/client/clientset/versioned/fake"
	informers "github.frg.tech/cloud/fanplane/pkg/apis/client/informers/externalversions"
	fanplanecontroller "github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
)

type fixture struct {
	t *testing.T

	client     *fake.Clientset

	// Objects to put in the store.
	envoyBoostrapLister []*fanplanecontroller.EnvoyBootstrap
	gatewayListener     []*fanplanecontroller.Gateway

	// Actions expected to happen on the client.
	kubeactions []core.Action
	actions     []core.Action

	// Objects from here preloaded into NewSimpleFake.
	kubeobjects []runtime.Object
	objects     []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}
	f.kubeobjects = []runtime.Object{}
	return f
}

func newEnvoyBootsrap(name string, replicas *int32) *fanplanecontroller.EnvoyBootstrap {

	return &fanplanecontroller.EnvoyBootstrap{
		TypeMeta: metav1.TypeMeta{APIVersion: fanplanecontroller.SchemeGroupVersion.String(), Kind: "EnvoyBootstrap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
	}
}

func (f *fixture) newController() (*Controller, informers.SharedInformerFactory) {
	f.client = fake.NewSimpleClientset(f.objects...)

	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc())
	snapCache := cache.NewCache()

	c := NewController(snapCache, f.client,
		i.Fanplane().V1alpha1().Gateways(),
		i.Fanplane().V1alpha1().EnvoyBootstraps())

	c.envoySynced = alwaysReady

	c.recorder = &record.FakeRecorder{}

	for _, fanplaneObject := range f.envoyBoostrapLister {
		i.Fanplane().V1alpha1().EnvoyBootstraps().Informer().GetIndexer().Add(fanplaneObject)
	}

	for _, fanplaneObject := range f.gatewayListener {
		i.Fanplane().V1alpha1().Gateways().Informer().GetIndexer().Add(fanplaneObject)
	}

	test, _, _ := i.Fanplane().V1alpha1().EnvoyBootstraps().Informer().GetIndexer().Get(f.envoyBoostrapLister[0])
	fmt.Println(test)

	return c, i
}

func (f *fixture) run(envoyBootstrapName string) {
	f.runController(envoyBootstrapName, true, false)
}

func (f *fixture) runExpectError(envoyBootstrapName string) {
	f.runController(envoyBootstrapName, true, true)
}

func (f *fixture) runController(envoyName string, startInformers bool, expectError bool) {
	controller, informerFactory := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		informerFactory.Start(stopCh)
	}

	err := controller.syncHandlerEnvoyBootstrap(envoyName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing EnvoyBootstrap: %v", err)
	} else if expectError && err == nil {
		f.t.Error("expected error syncing EnvoyBootstrap, got nil")
	}

	actions := filterInformerActions(f.client.Actions())
	for i, action := range actions {
		if len(f.actions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(actions)-len(f.actions), actions[i:])
			break
		}

		expectedAction := f.actions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}

}

// checkAction verifies that expected and actual actions are equal and both have
// same attached resources
func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.CreateAction:
		e, _ := expected.(core.CreateAction)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expObject, object))
		}
	case core.UpdateAction:
		e, _ := expected.(core.UpdateAction)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expObject, object))
		}
	case core.PatchAction:
		e, _ := expected.(core.PatchAction)
		expPatch := e.GetPatch()
		patch := a.GetPatch()

		if !reflect.DeepEqual(expPatch, expPatch) {
			t.Errorf("Action %s %s has wrong patch\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintDiff(expPatch, patch))
		}
	}
}

// filterInformerActions filters list and watch actions for testing resources.
// Since list and watch don't change resource state we can filter it to lower
// nose level in our tests.
func filterInformerActions(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			(action.Matches("list", "envoybootstraps") ||
				action.Matches("watch", "envoybootstraps") ||
				action.Matches("list", "gateways") ||
				action.Matches("watch", "gateways")) {
			continue
		}
		ret = append(ret, action)
	}

	return ret
}

func (f *fixture) expectUpdateEnvoyBootstrapAction(envoybootstraps *fanplanecontroller.EnvoyBootstrap) {
	action := core.NewUpdateAction(schema.GroupVersionResource{Resource: "envoybootstraps"}, envoybootstraps.Namespace, envoybootstraps)
	f.actions = append(f.actions, action)
}

func (f *fixture) expectCreateEnvoyBootstrapAction(envoybootstraps *fanplanecontroller.EnvoyBootstrap) {
	action := core.NewCreateAction(schema.GroupVersionResource{Resource: "envoybootstraps"}, envoybootstraps.Namespace, envoybootstraps)
	f.actions = append(f.actions, action)
}

func getKey(envoyBootstrap *fanplanecontroller.EnvoyBootstrap, t *testing.T) string {
	key, err := kubeCache.DeletionHandlingMetaNamespaceKeyFunc(envoyBootstrap)
	if err != nil {
		t.Errorf("Unexpected error getting key for EnvoyBootstrap %v: %v", envoyBootstrap.Name, err)
		return ""
	}
	return key
}

func TestCreatesEnvoyBootstrap(t *testing.T) {
	f := newFixture(t)
	envoyBootstrap, err := fanplanecontroller.LoadEnvoyBootstrap("../testdata/envoy.yaml")
	envoyBootstrap.Status.Processed = true
	if err != nil {
		t.Fatal("couldn't parse testdata")
	}

	f.envoyBoostrapLister = append(f.envoyBoostrapLister, envoyBootstrap)
	f.objects = append(f.objects, envoyBootstrap)

	f.expectUpdateEnvoyBootstrapAction(envoyBootstrap)

	f.run(getKey(envoyBootstrap, t))
}

func TestDoNothing(t *testing.T) {
	f := newFixture(t)
	envoyBootstrapInvalid := newEnvoyBootsrap("test", int32Ptr(1))

	f.envoyBoostrapLister = append(f.envoyBoostrapLister, envoyBootstrapInvalid)
	f.objects = append(f.objects, envoyBootstrapInvalid)

	f.runExpectError(getKey(envoyBootstrapInvalid, t))
}

func TestUpdateDeployment(t *testing.T) {
	f := newFixture(t)
	envoyBootstrap, err := fanplanecontroller.LoadEnvoyBootstrap("../testdata/envoy.yaml")
	envoyBootstrap.Status.Processed = true
	if err != nil {
		t.Fatal("couldn't parse testdata")
	}

	f.envoyBoostrapLister = append(f.envoyBoostrapLister, envoyBootstrap)
	f.objects = append(f.objects, envoyBootstrap)

	f.expectUpdateEnvoyBootstrapAction(envoyBootstrap)

	f.run(getKey(envoyBootstrap, t))
}

func int32Ptr(i int32) *int32 { return &i }
