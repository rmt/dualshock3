package dualshock3

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
)

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
	return fmt.Sprintf("LeftStick=(%3d[%0.2f],%3d[%0.2f],%s) RightStick=(%3d[%0.2f],%3d[%0.2f],%s) L1=%s L2=%02x R1=%s R2=%02x DPad(%s,%s,%s,%s) Square=%s Triangle=%s Cross=%s Circle=%s Select=%s Start=%s",
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

func (ctrl *GamePadControls) Run(callback OnChangeFunc) error {
	if ctrl == nil {
		return nil
	}
	for {
		events, err := ctrl.Device.Read()
		if err != nil {
			return err
		}
		keyEvent := false
		leftChanged := false
		rightChanged := false

		for _, ev := range events {
			switch ev.Type {
			case evdev.EV_SYN:
				//
			case evdev.EV_ABS:
				//
				switch ev.Code {
				case 0:
					if ctrl.LeftStick.setX(ev.Value) {
						leftChanged = true
					}
				case 1:
					if ctrl.LeftStick.setY(ev.Value) {
						leftChanged = true
					}
				case 2:
					ctrl.L2 = ev.Value
					leftChanged = true
				case 3:
					if ctrl.RightStick.setX(ev.Value) {
						rightChanged = true
					}
				case 4:
					if ctrl.RightStick.setY(ev.Value) {
						rightChanged = true
					}
				case 5:
					ctrl.R2 = ev.Value
					rightChanged = true
				default:
					//
				}
			case evdev.EV_KEY:
				keyEvent = true
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
		if keyEvent || leftChanged || rightChanged {
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
