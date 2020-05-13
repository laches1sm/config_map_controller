package controller

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestUpdateConfigMap(t *testing.T) {
	Convey("Given I have a controller", func() {
		c := NewConfigMapController()
		Convey("And I have a ConfigMap with the right annotation for the controller", func() {
			ch := make(chan struct{})
			c.Run(ch)
			var ClientSet kubernetes.Clientset

			testNamespaceName := "default"
			testConfigMapName := "test-configmap"

			testConfigMap, err := ClientSet.CoreV1().ConfigMaps(testNamespaceName).Create(&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: testConfigMapName,
					Labels: map[string]string{
						"x-k8s.io/curl-me-that": "",
					},
				},
				Data: map[string]string{
					"valueName": "value",
				},
			})
			if err != nil {
				fmt.Println(err)
			}
			Convey("And I am able to successfully get a joke", func() {
				joke := "Why did the chicken cross the road? To get to the other side!"
				Convey("Then I should see my ConfigMap updated with the joke")
				So(testConfigMap.ObjectMeta.Labels, ShouldAlmostEqual, joke)
			})

		})

	})
}
