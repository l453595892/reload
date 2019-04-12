package reload

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	"time"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/fields"
)

var UploadObejectList []runtime.Object

func NewConfigmapController(clientset *kubernetes.Clientset) cache.Controller{
	wl := cache.NewListWatchFromClient(clientset.Core().RESTClient(), string(v1.ResourceServices), v1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(
		wl,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("service added: %s \n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("service deleted: %s \n", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("service changed \n")
			},
		},
	)
	return controller
}

