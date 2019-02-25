package main

import (
	"fmt"
	"github.com/rmt/dualshock3"
	"os"
)

func main() {
	ctrl, err := dualshock3.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(ctrl.Device)
	fmt.Println(ctrl.Device.Capabilities)
	ctrl.Run(func() {
		fmt.Println(ctrl)
		if ctrl.Start && ctrl.Select {
			ctrl.Quit()
		}
	})
}
