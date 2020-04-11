package dualshock3

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"sync"
)

type Tilt struct {
	Xraw	int32
	Yraw	int32
	Zraw	int32
}

func (t *Tilt) String() string {
	if t == nil {
		return "-,-,-"
	}
	return fmt.Sprintf("%d,%d,%d", t.Xraw, t.Yraw, t.Zraw)
}

type AnalogStick struct {
	X    float64 // -1.0 to 1.0 (<0 is left, >0 is right)
	Y    float64 // -1.0 to 1.0 (<0 is up, >0 is right)
	Xraw int32
	Yraw int32
}

func analogCalc(v int32) (int32, float64) {
	// 0 - 255, with a center error margin of 8
	if v >= 119 && v <= 135 {
		return v, 0.0
	}
	if v > 127 {
		return v, (float64(v-127) / 128.0)
	}
	return v, (-float64(128-v) / 128.0)
}

func (g *AnalogStick) setX(v int32) bool {
	v, vpc := analogCalc(v)
	if g.X == vpc {
		return false
	}
	g.Xraw = v
	g.X = vpc
	return true
}

func (g *AnalogStick) setY(v int32) bool {
	v, vpc := analogCalc(v)
	if g.Y == vpc {
		return false
	}
	g.Yraw = v
	g.Y = vpc
	return true
}

type GamePadControls struct {
	sync.Mutex
	Device     *evdev.InputDevice
	Motion     *evdev.InputDevice
	Bluetooth  bool
	Tilt	   *Tilt
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

func toBool(v int32) bool {
	return v != 0
}

func pb(b bool) string {
	if b {
		return "X"
	}
	return "O"
}

func (g GamePadControls) String() string {
	return fmt.Sprintf("Tilt=(%s) LeftStick=(%3d[%0.2f],%3d[%0.2f],%s) RightStick=(%3d[%0.2f],%3d[%0.2f],%s) L1=%s L2=%02x R1=%s R2=%02x DPad(%s,%s,%s,%s) Square=%s Triangle=%s Cross=%s Circle=%s Select=%s Start=%s",
		g.Tilt.String(),
		g.LeftStick.Xraw, g.LeftStick.X,
		g.LeftStick.Yraw, g.LeftStick.Y, pb(g.L3),
		g.RightStick.Xraw, g.RightStick.X,
		g.RightStick.Yraw, g.RightStick.Y, pb(g.R3),
		pb(g.L1), g.L2, pb(g.R1), g.R2,
		pb(g.DPadLeft), pb(g.DPadRight), pb(g.DPadUp), pb(g.DPadDown),
		pb(g.Square), pb(g.Triangle), pb(g.Cross), pb(g.Circle),
		pb(g.Select), pb(g.Start),
	)
}

func (ctrl *GamePadControls) Quit() {
	if ctrl == nil {
		return
	}
	ctrl.quit = true
}

// optional callback if anything happens on the GamePadControls
type OnChangeFunc func()

func (ctrl *GamePadControls) RunMotion(callback OnChangeFunc) error {
	if ctrl.Motion == nil {
		return nil
	}
	if ctrl.Tilt == nil {
		ctrl.Tilt = &Tilt{}
	}
	for {
		events, err := ctrl.Motion.Read()
		if err != nil {
			return err
		}
		changed := false

		ctrl.Lock()
		for _, ev := range events {
			if ev.Type != evdev.EV_ABS {
				continue
			}
			switch ev.Code {
			case 0:
				ctrl.Tilt.Xraw = ev.Value
				changed = true
			case 1:
				ctrl.Tilt.Yraw = ev.Value
				changed = true
			case 2:
				ctrl.Tilt.Zraw = ev.Value
				changed = true
			}
		}
		ctrl.Unlock()
		if changed && callback != nil {
			callback()
		}
		if ctrl.quit {
			return nil
		}
	}
}

func (ctrl *GamePadControls) Run(callback OnChangeFunc) error {
	if ctrl == nil {
		return nil
	}
	for {
		events, err := ctrl.Device.Read()
		if err != nil {
			return err
		}
		changed := false

		ctrl.Lock()
		for _, ev := range events {
			switch ev.Type {
			case evdev.EV_ABS:
				//
				switch ev.Code {
				case 0:
					if ctrl.LeftStick.setX(ev.Value) {
						changed = true
					}
				case 1:
					if ctrl.LeftStick.setY(ev.Value) {
						changed = true
					}
				case 2:
					ctrl.L2 = ev.Value
					changed = true
				case 3:
					if ctrl.RightStick.setX(ev.Value) {
						changed = true
					}
				case 4:
					if ctrl.RightStick.setY(ev.Value) {
						changed = true
					}
				case 5:
					ctrl.R2 = ev.Value
					changed = true
				default:
					//
				}
			case evdev.EV_KEY:
				changed = true
				boolValue := ev.Value != 0
				switch ev.Code {
				case 304:
					ctrl.Cross = boolValue
				case 305:
					ctrl.Circle = boolValue
				case 307:
					ctrl.Triangle = boolValue
				case 308:
					ctrl.Square = boolValue
				case 310:
					ctrl.L1 = boolValue
				case 311:
					ctrl.R1 = boolValue
				case 314:
					ctrl.Select = boolValue
				case 315:
					ctrl.Start = boolValue
				case 317:
					ctrl.L3 = boolValue
				case 318:
					ctrl.R3 = boolValue
				case 544:
					ctrl.DPadUp = boolValue
				case 545:
					ctrl.DPadDown = boolValue
				case 546:
					ctrl.DPadLeft = boolValue
				case 547:
					ctrl.DPadRight = boolValue
				case 312: // L2 as button
				case 313: // R2 as button
				case 4: // for any button click, it seems
					//default:
					//	fmt.Println(ev.String())
				}
			}
		}
		ctrl.Unlock()
		if changed {
			if callback != nil {
				callback()
			}
		}
		if ctrl.quit {
			return nil
		}
	}
}

func Open(devpath string) (*GamePadControls, error) {
	device, err := evdev.Open(devpath)
	if err != nil {
		return nil, err
	}
	return &GamePadControls{Device: device}, nil
}

func OpenFirst() (*GamePadControls, error) {
	devices, err := evdev.ListInputDevicePaths("/dev/input/event*")
	if err != nil {
		return nil, err
	}
	var motionDev, inputDev *evdev.InputDevice

	bluetooth := false
	for _, devpath := range(devices) {
		device, err := evdev.Open(devpath)
		if err != nil {
			continue
		}

		// 0x054c, product 0x0268,
		// SHANWAN PS3 GamePad Motion Sensors usb-0000:00:14.0-9/input0 3 1356 616 33040
        // SHANWAN PS3 GamePad usb-0000:00:14.0-9/input0 3 1356 616 33040
		if device.Vendor == 0x054c && device.Product == 0x0268 {
			hasAbs := false
			hasKey := false
			for k, _ := range(device.Capabilities) {
				if k.Name == "EV_ABS" {
					hasAbs = true
				} else if k.Name == "EV_KEY" {
					hasKey = true
				}
			}
			if hasKey {
				inputDev = device
			} else if hasAbs {
				motionDev = device
			} else {
				fmt.Println("Ignoring input device with unexpected capabilities:", device.Name, device.Capabilities)
				continue
			}
		// FIXME: switch to Vendor/Product/Capabilities logic for bluetooth devices too
		} else if device.Name == "Sony Computer Entertainment Wireless Controller Motion Sensors" {
			motionDev = device
			bluetooth = true
		} else if device.Name == "Sony Computer Entertainment Wireless Controller" {
			inputDev = device
			bluetooth = true
		} else if device.Name == "Motion Controller" { // Move controller
			fmt.Println("Found the Move Controller... But it's not yet supported..")
		} else {
			fmt.Println(device.Name, device.Phys, device.Bustype, device.Vendor, device.Product, device.Version)
			device.File.Close()
			continue
		}

		if motionDev != nil && inputDev != nil {
			return &GamePadControls{
				Device: inputDev,
				Motion: motionDev,
				Bluetooth: bluetooth,
			}, nil
		}
	}
	// 
	if inputDev != nil {
		fmt.Println("No motion sensors detected?!?!?")
		return &GamePadControls{
			Device: inputDev,
		}, nil
	}
	return nil, nil
}
