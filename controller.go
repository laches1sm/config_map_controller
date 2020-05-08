package controller

import (
	"github.com/kubernetes/client-go/kubernetes"
	"github.com/kubernetes/client-go/tools/cache"
	"github.com/kubernetes/client-go/util/workqueue"
	"log"
)

// ClientSetController is a representation of the controller used to modify k8s ClientSets
// with a particular annotation.
type ClientSetController struct {
	logger     log.Logger
	kInterface kubernetes.Interface
	queue      workqueue.RateLimitingInterface
	informer   cache.SharedIndexInformer
}

type clientSetController interface {
	UpdateConfigMap()
	RecordEvent()
}

// NewClientSetController returns an instance of a ClientSetController
func NewClientSetController(log log.Logger, kInterface kubernetes.Interface, queue workqueue.RateLimitingInterface, informer cache.SharedIndexInformer) *ClientSetController {
	return &ClientSetController{logger: log, kInterface: kInterface, queue: queue, informer: informer}
}

func (*ClientSetController) UpdateConfigMap() {

}

func (*ClientSetController) RecordEvent() {

}

func getURLContents() {

}
