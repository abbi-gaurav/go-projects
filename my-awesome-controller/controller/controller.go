package controller

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/db"
	types "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/apis/awesome.controller.io/v1"
	clientSet "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/clientset/versioned"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	listers "github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/listers/awesome.controller.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

type CakeController struct {
	db         db.DB
	informer   v1.CakeInformer
	workQueue  workqueue.RateLimitingInterface
	cakeSynced cache.InformerSynced
	lister     listers.CakeLister
	clientSet  clientSet.Interface
}

func New(informer v1.CakeInformer, database db.DB, clientSet clientSet.Interface) *CakeController {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	cakeController := &CakeController{
		db:         database,
		informer:   informer,
		workQueue:  queue,
		cakeSynced: informer.Informer().HasSynced,
		lister:     informer.Lister(),
		clientSet:  clientSet,
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: cakeController.enqueue,
		UpdateFunc: func(old, new interface{}) {
			cakeController.enqueue(new)
		},
	})

	return cakeController

}

func (c *CakeController) enqueue(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilRuntime.HandleError(err)
		return
	}

	c.workQueue.AddRateLimited(key)
}

func (c *CakeController) Run(parallelism int, stopCh <-chan struct{}) error {
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

func (c *CakeController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *CakeController) processNextItem() bool {
	obj, shutDown := c.workQueue.Get()
	if shutDown {
		return false
	}

	err := c.process(obj)

	if err != nil {
		utilRuntime.HandleError(err)
		return true
	}
	return true
}

func (c *CakeController) process(obj interface{}) error {
	defer c.workQueue.Done(obj)
	var key string
	var ok bool

	if key, ok = obj.(string); !ok {
		c.workQueue.Forget(obj)
		utilRuntime.HandleError(fmt.Errorf("expected string in workqueue, but got #%v", obj))
		return nil
	}

	log.Println("calling handle")
	if err := c.handle(key); err != nil {
		c.workQueue.AddRateLimited(key)
		return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
	}
	c.workQueue.Forget(obj)
	fmt.Printf("Successfully synced #%s", key)
	return nil
}

func (c *CakeController) handle(key string) error {
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

	dbValue := c.db.Get(key)

	if dbValue == nil {
		c.db.Add(key, cake)
		err = c.updateStatus(cake, types.ADDED)
		return err
	} else {
		if dbValue.Spec != cake.Spec {
			c.db.Update(key, cake)
			err = c.updateStatus(cake, types.UPDATED)
			return err
		}
		return nil
	}
}

func (c *CakeController) updateStatus(cake *types.Cake, state string) error {
	cakeCopy := cake.DeepCopy()
	cakeCopy.Status.State = state
	_, err := c.clientSet.AwesomeV1().Cakes(cake.Namespace).Update(cakeCopy)
	return err
}

func (c *CakeController) ShutDown() {
	c.workQueue.ShutDown()
}
