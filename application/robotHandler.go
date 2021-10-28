package application

import (
	"gobot.io/x/gobot/platforms/keyboard"

	"github.com/jtonynet/autogo/config"
	domain "github.com/jtonynet/autogo/domain"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func Init(keys *keyboard.Driver, motors *output.Motors, servoKit *output.Servos, lcd *output.Display, sonarSet *input.SonarSet, cfg *config.Config) {

	if cfg.ServoKit.Enabled {
		servoPan := servoKit.GetByName("pan")
		servoTilt := servoKit.GetByName("tilt")

		servoKit.Init()
		servoKit.SetCenter(servoPan)
		servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
	}

	if cfg.ArduinoSonar.Enabled && cfg.Motors.Enabled {
		go domain.SonarWorker(sonarSet, motors, lcd, cfg)
	}

	keys.On(keyboard.Key, func(data interface{}) {
		key := data.(keyboard.KeyEvent)
		domain.ControllByKeyboard(key, motors, servoKit, lcd, sonarSet, cfg)
	})
}
