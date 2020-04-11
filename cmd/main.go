package main

import (
	"fmt"
	"github.com/rmt/dualshock3"
)

func main() {
	ctrl, err := dualshock3.OpenFirst()
	if err != nil {
		panic(err)
	}
	fmt.Println(ctrl.Device)
	fmt.Println(ctrl.Device.Capabilities)
	if ctrl.Motion != nil {
		fmt.Println(ctrl.Motion)
		fmt.Println(ctrl.Motion.Capabilities)
	}
	callback := func() {
		fmt.Println(ctrl)
		if ctrl.Start && ctrl.Select {
			ctrl.Quit()
		}
	}
	go ctrl.RunMotion(callback)
	ctrl.Run(callback)
}
