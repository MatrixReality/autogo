package application

import (
	config "github.com/jtonynet/autogo/config"
	domain "github.com/jtonynet/autogo/domain"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	output "github.com/jtonynet/autogo/peripherals/actuators"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func Init(messageBroker *infrastructure.MessageBroker, kbd *sensors.Keyboard, motors *output.Motors, servoKit *output.Servos, lcd *output.Display, sonarSet *sensors.SonarSet, imu *sensors.IMU, cfg *config.Config) {
	keys := kbd.Driver
	robotDomain := domain.NewRobot(messageBroker, motors, servoKit, lcd, sonarSet, imu, cfg)

	keys.On(kbd.Key, func(data interface{}) {
		robotDomain.ControllByKeyboard(data)
	})
}
