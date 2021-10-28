package domain

import (
	"fmt"

	"gobot.io/x/gobot/platforms/keyboard"

	"github.com/jtonynet/autogo/config"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

var (
	direction         string = ""
	lcdMsg            string = ""
	colissionDetected bool   = false
)

func ControllByKeyboard(key keyboard.KeyEvent, motors *output.Motors, servoKit *output.Servos, lcd *output.Display, sonarSet *input.SonarSet, cfg *config.Config) {
	oldDirection := direction

	if cfg.ServoKit.Enabled {
		servoPan := servoKit.GetByName("pan")
		servoTilt := servoKit.GetByName("tilt")

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)

		if key.Key == keyboard.W {
			newTilt := tiltAngle - cfg.ServoKit.PanTiltFactor
			if newTilt < servoKit.TiltPos["top"] {
				newTilt = servoKit.TiltPos["top"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + cfg.ServoKit.PanTiltFactor
			if newTilt > servoKit.TiltPos["down"] {
				newTilt = servoKit.TiltPos["down"]
			}
			servoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + cfg.ServoKit.PanTiltFactor
			if newPan > servoKit.PanPos["left"] {
				newPan = servoKit.PanPos["left"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - cfg.ServoKit.PanTiltFactor
			if newPan < servoKit.PanPos["right"] {
				newPan = servoKit.PanPos["right"]
			}
			servoKit.SetAngle(servoPan, uint8(newPan))
		} else if key.Key == keyboard.X {
			servoKit.SetCenter(servoPan)
			servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
		}
	}

	if cfg.Motors.Enabled {
		if key.Key == keyboard.ArrowUp && colissionDetected == false {
			motors.Forward(cfg.Motors.MaxSpeed)
			direction = "Front"
			lcdMsg = direction
		} else if key.Key == keyboard.ArrowDown {
			motors.Backward(cfg.Motors.MaxSpeed)
			direction = "Back"
			lcdMsg = direction
		} else if key.Key == keyboard.ArrowRight {
			motors.Left(cfg.Motors.MaxSpeed)
			direction = "Right"
			lcdMsg = direction
		} else if key.Key == keyboard.ArrowLeft {
			motors.Right(cfg.Motors.MaxSpeed)
			direction = "Left"
			lcdMsg = direction
		} else if key.Key == keyboard.Q {
			motors.Stop()
			direction = ""
			lcdMsg = cfg.Version + " Arrow key"
		} else {
			fmt.Println(lcdMsg, key, key.Char)
		}
	}

	if cfg.LCD.Enabled && oldDirection != direction {
		lcd.ShowMessage(lcdMsg, output.LINE_2)
	}

}

func SonarWorker(sonarSet *input.SonarSet, motors *output.Motors, lcd *output.Display, cfg *config.Config) {
	for true {
		sonarData, err := sonarSet.GetData()
		if err == nil {
			if sonarData["center"] <= cfg.ArduinoSonar.MinStopValue && direction == "Front" && colissionDetected == false {
				colissionDetected = true
				motors.Stop()

				if cfg.LCD.Enabled {
					s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
					lcd.ShowMessage(s, output.LINE_2)
				}

			} else if colissionDetected && direction != "Front" {
				colissionDetected = false
			}

		}
	}
}
