package domain

import (
	output "github.com/jtonynet/autogo/peripherals/output"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/keyboard"
)

type Servos struct {
	Kit  *output.Servos
	Pan  *gpio.ServoDriver
	Tilt *gpio.ServoDriver
}

func NewServos(servoKit *output.Servos) *Servos {
	this := &Servos{
		Kit:  servoKit,
		Pan:  servoKit.GetByName("pan"),
		Tilt: servoKit.GetByName("tilt"),
	}

	servoKit.Init()
	servoKit.SetCenter(this.Pan)
	servoKit.SetAngle(this.Tilt, uint8(servoKit.TiltPos["horizon"]))

	return this
}

func (this *Servos) ControllPanAndTilt(k int) {
	cfg := this.Kit.Cfg
	servoPan := this.Pan
	servoTilt := this.Tilt

	panAngle := int(servoPan.CurrentAngle)
	tiltAngle := int(servoTilt.CurrentAngle)

	switch k {
	case keyboard.W:
		newTilt := tiltAngle - cfg.PanTiltFactor
		if newTilt < this.Kit.TiltPos["top"] {
			newTilt = this.Kit.TiltPos["top"]
		}
		this.Kit.SetAngle(servoTilt, uint8(newTilt))

	case keyboard.S:
		newTilt := tiltAngle + cfg.PanTiltFactor
		if newTilt > this.Kit.TiltPos["down"] {
			newTilt = this.Kit.TiltPos["down"]
		}
		this.Kit.SetAngle(servoTilt, uint8(newTilt))

	case keyboard.A:
		newPan := panAngle + cfg.PanTiltFactor
		if newPan > this.Kit.PanPos["left"] {
			newPan = this.Kit.PanPos["left"]
		}
		this.Kit.SetAngle(servoPan, uint8(newPan))

	case keyboard.D:
		newPan := panAngle - cfg.PanTiltFactor
		if newPan < this.Kit.PanPos["right"] {
			newPan = this.Kit.PanPos["right"]
		}
		this.Kit.SetAngle(servoPan, uint8(newPan))

	case keyboard.X:
		this.Kit.SetCenter(servoPan)
		this.Kit.SetAngle(servoTilt, uint8(this.Kit.TiltPos["horizon"]))
	}
}
