package application

import (
	"github.com/jtonynet/autogo/config"
	domain "github.com/jtonynet/autogo/domain"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
	"gobot.io/x/gobot/platforms/keyboard"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func Init(kbd *input.Keyboard, motors *output.Motors, servoKit *output.Servos, lcd *output.Display, sonarSet *input.SonarSet, cfg *config.Config) {
	keys := kbd.Driver
	robotDomain := domain.NewRobot(motors, servoKit, lcd, sonarSet, cfg)

	keys.On(kbd.Key, func(data interface{}) {
		//key := data.(input.KeyEvent)
		key := data.(keyboard.KeyEvent)
		robotDomain.ControllByKeyboard(key)
	})
}
