package main

import (
	"flag"
	"fmt"
	"github.com/bendahl/uinput"
	"github.com/rmt/dualshock3"
	"os"
)

func main() {
	wasd := false

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-wasd]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
This program converts a PS3 controller's direction & button presses into
key strokes under Linux.  You can then use your PS3 controller to play
keyboard-only games on Linux, or to control your editor.. whatever.

The key bindings are:
  * LeftStick -> arrow keys [or WASD]
  * RightStick -> keypad number keys (8, 4, 2, 6)
  * Cross -> Space
  * Circle -> Enter
  * Square -> Slash (/)
  * Triangle -> Backslash (\)
  * L1 -> LeftShift
  * L2 -> RightShift
  * L2 -> LeftCtrl
  * R2 -> RightCtrl
  * Start -> 1
  * Select -> i
`)
	}
	flag.BoolVar(&wasd, "wasd", false, "use WASD layout")
	flag.Parse()

	var keyUp = uinput.KeyUp
	var keyDown = uinput.KeyDown
	var keyLeft = uinput.KeyLeft
	var keyRight = uinput.KeyRight
	var keyL3 = uinput.KeyX
	var keyCross = uinput.KeySpace
	var keyCircle = uinput.KeyEnter
	var keySquare = uinput.KeySlash
	var keyTriangle = uinput.KeyBackslash
	var keyL1 = uinput.KeyLeftshift
	var keyR1 = uinput.KeyRightshift
	var keyL2 = uinput.KeyLeftctrl
	var keyR2 = uinput.KeyRightctrl
	var kpLeft = uinput.KeyKp4
	var kpRight = uinput.KeyKp6
	var kpUp = uinput.KeyKp8
	var kpDown = uinput.KeyKp2
	var keyR3 = uinput.KeyKp5
	var keyStart = uinput.Key1
	var keySelect = uinput.KeyI

	if wasd {
		keyUp = uinput.KeyW
		keyDown = uinput.KeyS
		keyLeft = uinput.KeyA
		keyRight = uinput.KeyD
	}

	kb, err := uinput.CreateKeyboard("/dev/uinput", []byte("ps3kb"))
	if err != nil {
		panic(err)
	}

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

	// doKey presses & releases a key in response to controller state changes
	// buttonPressed should be true if the controller button is down/on
	// keyDown is managed by doKey (tracks whether the kbd key is currently down)
	// keyCode is the key code of the key to press/release
	doKey := func(buttonPressed bool, keyDown *bool, keycode int) {
		if buttonPressed {
			if !*keyDown {
				*keyDown = true
				kb.KeyDown(keycode)
			}
		} else if *keyDown {
			*keyDown = false
			kb.KeyUp(keycode)
		}
	}

	var left, right, up, down, cross, circle, triangle, square bool
	var kpleft, kpright, kpup, kpdown, l1, l2, l3, r1, r2, r3, start, sel bool
	callback := func() {
		//fmt.Println(ctrl)
		if ctrl.Start && ctrl.Select {
			ctrl.Quit()
		}
		doKey(ctrl.LeftStick.X < -0.3 || ctrl.DPadLeft, &left, keyLeft)
		doKey(ctrl.LeftStick.X > 0.3 || ctrl.DPadRight, &right, keyRight)
		doKey(ctrl.LeftStick.Y < -0.3 || ctrl.DPadUp, &up, keyUp)
		doKey(ctrl.LeftStick.Y > 0.3 || ctrl.DPadDown, &down, keyDown)
		doKey(ctrl.L3, &l3, keyL3)
		doKey(ctrl.RightStick.X < -0.3, &kpleft, kpLeft)
		doKey(ctrl.RightStick.X > 0.3, &kpright, kpRight)
		doKey(ctrl.RightStick.Y < -0.3, &kpup, kpUp)
		doKey(ctrl.RightStick.Y > 0.3, &kpdown, kpDown)
		doKey(ctrl.R3, &r3, keyR3)
		doKey(ctrl.Cross, &cross, keyCross)
		doKey(ctrl.Circle, &circle, keyCircle)
		doKey(ctrl.Square, &square, keySquare)
		doKey(ctrl.Triangle, &triangle, keyTriangle)
		doKey(ctrl.L1, &l1, keyL1)
		doKey(ctrl.L2 > 0, &l2, keyL2)
		doKey(ctrl.R1, &r1, keyR1)
		doKey(ctrl.R2 > 0, &r2, keyR2)
		doKey(ctrl.Select, &sel, keySelect)
		doKey(ctrl.Start, &start, keyStart)
	}
	go ctrl.RunMotion(callback)
	ctrl.Run(callback)
}
