package main

import (
	controller "github.com/laches1sm/config_map_controller"
)

func main() {
	ctrl := controller.NewConfigMapController()
	c := make(chan struct{})
	ctrl.Run(c)
}
