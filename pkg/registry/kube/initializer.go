package kube

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/server"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.frg.tech/cloud/fanplane/pkg/apis/client/clientset/versioned"
	informers "github.frg.tech/cloud/fanplane/pkg/apis/client/informers/externalversions"
)

func Initialize(config *server.FanplaneConfig, snapCache cache.SnapshotCache) {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(config.MasterURL, config.KubeCfgFile)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err)
	}

	kubeClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building example clientset: %s", err)
	}

	fanplaneInformerFactory := informers.NewSharedInformerFactory(kubeClient, 0)

	controller := NewController(snapCache, kubeClient,
		fanplaneInformerFactory.Fanplane().V1alpha1().Gateways(),
		fanplaneInformerFactory.Fanplane().V1alpha1().EnvoyBootstraps())

	go fanplaneInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		log.Fatalf("Error running controller: %s", err)
	}
}
