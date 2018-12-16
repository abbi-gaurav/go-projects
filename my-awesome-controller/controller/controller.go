package controller

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/db"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	listers "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/listers/awesome.controller.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type CakeController struct {
	db         db.DB
	informer   v1.CakeInformer
	workQueue  workqueue.RateLimitingInterface
	cakeSynced cache.InformerSynced
	lister     listers.CakeLister
}

func New(opts *opts.Options, informer v1.CakeInformer) *CakeController {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	database := db.New(opts.DbType)
	cakeController := &CakeController{
		db:         database,
		informer:   informer,
		workQueue:  queue,
		cakeSynced: informer.Informer().HasSynced,
		lister:     informer.Lister(),
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: cakeController.enqueue,
	})

	return cakeController

}

func (c *CakeController) enqueue(obj interface{}) {
	var key string
	var err error

	println("enqueue called", obj)

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilRuntime.HandleError(err)
		return
	}
	c.workQueue.AddRateLimited(key)
}

func (c *CakeController) Run(parallelism int, stopCh <-chan struct{}) error {
	defer utilRuntime.HandleCrash()
	defer c.workQueue.ShutDown()

	go c.informer.Informer().Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.cakeSynced) {
		utilRuntime.HandleError(fmt.Errorf("failed to wait for cache syncs"))
		return nil
	}

	for i := 0; i < parallelism; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	return nil
}

func (c *CakeController) Run2(parallelism int, stopCh <-chan struct{}) error {
	go c.informer.Informer().Run(stopCh)
	return nil
}

func (c *CakeController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *CakeController) processNextItem() bool {
	obj, shutDown := c.workQueue.Get()
	if shutDown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workQueue.Forget(obj)
			utilRuntime.HandleError(fmt.Errorf("expected string in workqueue, but got #%v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workQueue.Forget(obj)
		fmt.Printf("Successfully synced #%s", key)
		return nil
	}(obj)

	if err != nil {
		utilRuntime.HandleError(err)
		return true
	}
	return true
}

func (c *CakeController) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilRuntime.HandleError(fmt.Errorf("invalid resource key #%s", key))
		return nil
	}

	cake, err := c.lister.Cakes(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf("'cake %s' in workqueue no longer exists", key))
			return nil
		}
		return err
	}

	if c.db.Get(key) == nil {
		c.db.Add(key, cake)
	}

	return nil
}