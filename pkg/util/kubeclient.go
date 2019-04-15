package util

import (
	"flag"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeClient *kubernetes.Clientset
	//kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func init() {
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/liweijie/go/src/craftli.co/reload/cmd/config")
	if err != nil {
		logrus.Error("kubernetes client config error")
	}
	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error("kubernetes client error")
	}
}
