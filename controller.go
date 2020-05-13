package controller

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kubernetes/client-go/kubernetes"
	"github.com/kubernetes/client-go/kubernetes/typed/core/v1"
	"github.com/kubernetes/client-go/tools/cache"
	"github.com/kubernetes/client-go/util/workqueue"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
)

// ConfigMapController is a representation of the controller used to modify k8s ClientSets
// with a particular annotation.
type ConfigMapController struct {
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	recorder  record.EventRecorder
}

type clientSetController interface {
	UpdateConfigMap()
	Run()
}

const (
	controllerAgentName = "config-map-controller"
)

// NewConfigMapController returns an instance of a ConfigMapController
func NewConfigMapController() *ConfigMapController {
	var clientset kubernetes.Interface
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = "x-k8s.io/curl-me-that"
				return clientset.CoreV1().ConfigMaps(metav1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = "x-k8s.io/curl-me-that"
				return clientset.CoreV1().ConfigMaps(metav1.NamespaceAll).Watch(options)
			},
		},
		&corev1.ConfigMap{},
		0, //Skip resync
		cache.Indexers{},
	)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&v1.EventSinkImpl{Interface: clientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	controller := &ConfigMapController{clientset: clientset, queue: queue, informer: informer, recorder: recorder}
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			controller.UpdateConfigMap(obj)
		},
		
		DeleteFunc: func(obj interface{}) {
			controller.UpdateConfigMap(obj)
		},
	})
	return controller

}

// UpdateConfigMap ...
func (c *ConfigMapController) UpdateConfigMap(obj interface{}) {
	joke, err := getURLContents("curl-a-joke.herokuapp.com")
	if err != nil{
		c.recorder.Event(nil, "error", err.Error(), "failed to get joke")
	}
	configMap := obj.(*corev1.ConfigMap)

			// Copy pod and update the annotation.
			newConfigMap := configMap.DeepCopy()
			ann := newConfigMap.ObjectMeta.Annotations
			if ann == nil {
				ann = make(map[string]string)
			}
			ann["x-k8s.io/curl-me-that"] = joke
			newConfigMap.ObjectMeta.Annotations = ann

			_, err = c.clientset.CoreV1().ConfigMaps(newConfigMap.ObjectMeta.Namespace).Update(newConfigMap)
			if err != nil {
				c.recorder.Event(nil, "error", err.Error(), "error while getting joke")
			}
}

// Run ...
func (c *ConfigMapController) Run(stopCh <-chan struct{}) {
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)
}

func getURLContents(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unsucessfull response from curl-a-joke, unable to place joke")

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	jokeResp := string(body)
	return jokeResp, nil

}
