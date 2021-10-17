package peripherals

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	TiltPos = map[string]int{
		"top":     0,
		"horizon": 130,
		"down":    180,
	}

	PanPos = map[string]int{
		"left":  180,
		"right": 0,
	}
)

func NewKitDriver(a *raspi.Adaptor, addr int, bus int) *i2c.PCA9685Driver {
	servoKitDriver := i2c.NewPCA9685Driver(a,
		i2c.WithBus(bus),
		i2c.WithAddress(addr))

	return servoKitDriver
}

func NewServo(kitDriver *i2c.PCA9685Driver, servoId string, servoName string) *gpio.ServoDriver {
	s := gpio.NewServoDriver(kitDriver, servoId)
	s.SetName(servoName)

	return s
}
