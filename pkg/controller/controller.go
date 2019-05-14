package controller

import (
	"craftli.co/reload/pkg/handler"
	"craftli.co/reload/pkg/resource"
	"craftli.co/reload/pkg/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stakater/Reloader/pkg/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type Controller struct {
	client    *kubernetes.Clientset
	informer  cache.Controller
	indexer   cache.Indexer
	queue     workqueue.RateLimitingInterface
	namespace string
}

func NewController(resource string) (*Controller, error) {

	c := Controller{
		client: util.KubeClient,
	}

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	listWatcher := cache.NewListWatchFromClient(util.KubeClient.CoreV1().RESTClient(), resource, v1.NamespaceAll, fields.Everything())
	indexer, informer := cache.NewIndexerInformer(listWatcher, kube.ResourceMap[resource], 0, cache.ResourceEventHandlerFuncs{
		UpdateFunc: c.Update,
		//TODO
		//DeleteFunc:
	}, cache.Indexers{})
	c.indexer = indexer
	c.informer = informer
	c.queue = queue
	return &c, nil
}

func (c *Controller) Update(oldObj, newObj interface{}) {
	c.queue.Add(resource.UpdateHandler{
		NewObj: newObj,
		OldObj: oldObj,
	})
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	logrus.Infof("Stopping Controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	resourceHandler, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(resourceHandler)
	err := resourceHandler.(handler.ObjectHandler).Handler()
	c.handleErr(err, resourceHandler)
	return true
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		logrus.Errorf("Error syncing events %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	logrus.Infof("Dropping the key %q out of the queue: %v", key, err)
}
