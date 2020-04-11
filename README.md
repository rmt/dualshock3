Golang library for Playstation Dualshock 3 controller on Linux
==============================================================

This is a Linux-only package, relying on the /dev/input/ devices.

When you call ctrl.Run(), it launches a Go routine to monitor the device.

Whenever there is a key press (or an analog stick change), it will update the
struct with the changes.  There's an optional callback function whenever a
change happens.

It's been tested with the DualShock 3 connected via USB.

Additionally, there are some example programs to monitor the device and to
simulate keystrokes (using uinput).

Structs
-------

```
type GamePadControls struct {
	Device     *evdev.InputDevice
	LeftStick  AnalogStick
	L3         bool // analog button
	RightStick AnalogStick
	R3         bool // analog button
	L1         bool
	L2         int32
	R1         bool
	R2         int32
	DPadLeft   bool
	DPadRight  bool
	DPadUp     bool
	DPadDown   bool
	Square     bool
	Triangle   bool
	Cross      bool
	Circle     bool
	Select     bool
	Start      bool
	quit       bool
}

type AnalogStick struct {
	X    float64 // -1.0 to 1.0 (<0 is left, >0 is right)
	Y    float64 // -1.0 to 1.0 (<0 is up, >0 is right)
	Xraw int32
	Yraw int32
}
```

Example
-------

```
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
```
