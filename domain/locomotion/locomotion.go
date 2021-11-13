package domain

import (
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	output "github.com/jtonynet/autogo/peripherals/output"
	"gobot.io/x/gobot/platforms/keyboard"
)

type Locomotion struct {
	Motors *output.Motors
	Status *StatusDomain.Status
	LCD    *LcdDomain.LCD
}

func NewLocomotion(motors *output.Motors, lcd *LcdDomain.LCD, status *StatusDomain.Status) *Locomotion {
	this := &Locomotion{Motors: motors, LCD: lcd, Status: status}
	return this
}

func (this *Locomotion) Forward(speed byte) {
	this.Motors.Forward(speed)
}

func (this *Locomotion) Backward(speed byte) {
	this.Motors.Backward(speed)
}

func (this *Locomotion) Right(speed byte) {
	this.Motors.Right(speed)
}

func (this *Locomotion) Left(speed byte) {
	this.Motors.Left(speed)
}

func (this *Locomotion) Stop() {
	this.Motors.Stop()
}

func (this *Locomotion) ControllMoviment(k int) {
	oldDirection := this.Status.Direction
	cfg := this.Motors.Cfg

	switch k {
	case keyboard.ArrowUp:
		if !this.Status.ColissionDetected {
			this.Forward(cfg.MaxSpeed)
			this.Status.Direction = "Front"
			this.Status.LCDMsg = this.Status.Direction
		}

	case keyboard.ArrowDown:
		this.Backward(cfg.MaxSpeed)
		this.Status.Direction = "Back"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.ArrowRight:
		this.Left(cfg.MaxSpeed)
		this.Status.Direction = "Right"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.ArrowLeft:
		this.Right(cfg.MaxSpeed)
		this.Status.Direction = "Left"
		this.Status.LCDMsg = this.Status.Direction

	case keyboard.Q:
		this.Stop()
		this.Status.Direction = ""
		this.Status.LCDMsg = this.Status.Version + " Arrow key"
	}

	if this.LCD != nil && oldDirection != this.Status.Direction {
		this.LCD.ShowMessage(this.Status.LCDMsg, 2)
	}
}
