package handlers

import (
	"fmt"
	"log"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/keyboard"

	ArduinoSonarSet "github.com/matrixreality/autogo/peripherals/input"
	LCD "github.com/matrixreality/autogo/peripherals/output"
	Motors "github.com/matrixreality/autogo/peripherals/output"
	Servos "github.com/matrixreality/autogo/peripherals/output"
)

//TODO env vars on viper
const (
	VERSION         = "v0.0.5"
	PAN_TILT_FACTOR = 30
)

func InitKeyboard(servoKit *Servos.Servos, arduinoConn i2c.Connection, lcd *LCD.Display, motors *Motors.Motors, keys *keyboard.Driver) {
	firstRun := 1
	servoPan := servoKit.GetByName("pan")
	servoTilt := servoKit.GetByName("tilt")

	if firstRun == 1 {
		firstRun = 0
		servoKit.Init()
		Servos.SetCenter(servoPan)
		Servos.SetAngle(servoTilt, uint8(Servos.TiltPos["horizon"]))
	}

	keys.On(keyboard.Key, func(data interface{}) {
		key := data.(keyboard.KeyEvent)

		if key.Key == keyboard.B {
			sonarData, err := ArduinoSonarSet.GetData(arduinoConn)
			if err == nil {
				log.Println("///*********")
				log.Println("///Print arduino sonar data::")
				log.Println(sonarData)
				log.Println("///*********")
			}

		}

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)
		if key.Key == keyboard.W {
			newTilt := tiltAngle - PAN_TILT_FACTOR
			if newTilt < Servos.TiltPos["top"] {
				newTilt = Servos.TiltPos["top"]
			}
			Servos.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + PAN_TILT_FACTOR
			if newTilt > Servos.TiltPos["down"] {
				newTilt = Servos.TiltPos["down"]
			}
			Servos.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + PAN_TILT_FACTOR
			if newPan > Servos.PanPos["left"] {
				newPan = Servos.PanPos["left"]
			}
			Servos.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - PAN_TILT_FACTOR
			if newPan < Servos.PanPos["right"] {
				newPan = Servos.PanPos["right"]
			}
			Servos.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.X {
			Servos.SetCenter(servoPan)
			Servos.SetAngle(servoTilt, uint8(Servos.TiltPos["horizon"]))
		}

		if key.Key == keyboard.ArrowUp {
			motors.Forward(255)
			lcd.ShowMessage("Front", LCD.LINE_2)
		} else if key.Key == keyboard.ArrowDown {
			motors.Backward(255)
			lcd.ShowMessage("Back", LCD.LINE_2)
		} else if key.Key == keyboard.ArrowRight {
			motors.Left(255)
			lcd.ShowMessage("Left", LCD.LINE_2)
		} else if key.Key == keyboard.ArrowLeft {
			motors.Right(255)
			lcd.ShowMessage("Right", LCD.LINE_2)
		} else if key.Key == keyboard.Q {
			motors.Stop()
			lcd.ShowMessage(VERSION+" Arrow key", LCD.LINE_2)
		} else {
			fmt.Println("keyboard event!", key, key.Char)
		}
	})
}
