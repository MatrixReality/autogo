package peripherals

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

var (
	Kit = map[string]*gpio.ServoDriver{}

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

func NewDriver(a *raspi.Adaptor, bus int, addr int) *i2c.PCA9685Driver {
	driver := i2c.NewPCA9685Driver(a,
		i2c.WithBus(bus),
		i2c.WithAddress(addr))

	return driver
}

func Add(kitDriver *i2c.PCA9685Driver, servoId string, servoName string) *gpio.ServoDriver {
	s := gpio.NewServoDriver(kitDriver, servoId)
	s.SetName(servoName)

	Kit[servoName] = s

	return s
}

func GetByName(kitDriver *i2c.PCA9685Driver, servoName string) *gpio.ServoDriver {
	return Kit[servoName]
}

func SetAngle(s *gpio.ServoDriver, angle uint8) {
	s.Move(angle)
}

func SetCenter(s *gpio.ServoDriver) {
	s.Center()
}

func SetMin(s *gpio.ServoDriver) {
	s.Min()
}

func SetMax(s *gpio.ServoDriver) {
	s.Max()
}
