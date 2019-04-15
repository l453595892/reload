package main

import (
	"craftli.co/reload/pkg/controller"
	_ "craftli.co/reload/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

var (
	ResourceMap = map[string]runtime.Object{
		//string(v1.ResourceConfigMaps): &v1.ConfigMap{},
		//string(v1.ResourceSecrets):    &v1.Secret{},
		string(v1.ResourceServices): &v1.Service{},
	}
)

func main() {
	stop := make(chan struct{})
	defer close(stop)
	for k := range ResourceMap {
		c, err := controller.NewController(k, "hyper")
		if err != nil {
			logrus.Fatalf("%s", err)
		}

		logrus.Infof("Starting Controller to watch resource type: %s", k)
		go c.Run(1, stop)
	}
	for {
		time.Sleep(time.Second)
	}
}
