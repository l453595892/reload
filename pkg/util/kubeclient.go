package util

import (
	"flag"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	KubeClient *kubernetes.Clientset
)

func init() {
	flag.Parse()
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Error("kubernetes client config error")
	}
	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error("kubernetes client error")
	}
}
