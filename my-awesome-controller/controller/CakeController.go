package controller

import (
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/db"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type CakeController struct {
	db        db.DB
	informer  cache.SharedIndexInformer
	workQueue workqueue.RateLimitingInterface
}
