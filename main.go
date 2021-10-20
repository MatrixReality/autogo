package main

import (
	"fmt"
	"log"
	"net"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	ArduinoSonarSet "github.com/matrixreality/autogo/peripherals/input"
	LCD "github.com/matrixreality/autogo/peripherals/output"
	Motors "github.com/matrixreality/autogo/peripherals/output"
	Servos "github.com/matrixreality/autogo/peripherals/output"
)

//TODO env vars on viper
const (
	VERSION         = "v0.0.5"
	SERVOKIT_BUS    = 0
	SERVOKIT_ADDR   = 0x40
	ARDUINO_BUS     = 1
	ARDUINO_ADDR    = 0x18
	LCD_BUS         = 2
	LCD_ADDR        = 0x27
	LCD_COLLUMNS    = 16
	PAN_TILT_FACTOR = 30
)

func main() {
	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motorA, motorB := Motors.NewMotors(r)

	///SERVOKIT
	servoKit := Servos.NewDriver(r, SERVOKIT_BUS, SERVOKIT_ADDR)
	servoPan := servoKit.Add("0", "pan")
	servoTilt := servoKit.Add("1", "tilt")

	///ARDUINO SONAR SET
	arduinoConn, err := ArduinoSonarSet.GetConnection(r, ARDUINO_BUS, ARDUINO_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	///LCD
	lcd, err := LCD.NewLcd(LCD_BUS, LCD_ADDR, LCD_COLLUMNS)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.DeferAction()

	ip := GetOutboundIP()

	err = lcd.ShowMessage(string(ip), LCD.LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(VERSION+" Arrow key", LCD.LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	//Servos func, ArduinoSonarSet func, keys *Driver,
	firstRun := 1
	work := func() {
		if firstRun == 1 {
			firstRun = 0
			servoKit.Init()
			Servos.SetCenter(servoPan)
			Servos.SetAngle(servoTilt, uint8(Servos.TiltPos["horizon"]))
		}

		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			sonarData := ""
			if key.Key == keyboard.B {
				sonarData, err = ArduinoSonarSet.GetData(arduinoConn)
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
				Motors.Forward(255)
				lcd.ShowMessage("Front", LCD.LINE_2)
			} else if key.Key == keyboard.ArrowDown {
				Motors.Backward(255)
				lcd.ShowMessage("Back", LCD.LINE_2)
			} else if key.Key == keyboard.ArrowRight {
				Motors.Left(255)
				lcd.ShowMessage("Left", LCD.LINE_2)
			} else if key.Key == keyboard.ArrowLeft {
				Motors.Right(255)
				lcd.ShowMessage("Right", LCD.LINE_2)
			} else if key.Key == keyboard.Q {
				Motors.Stop()
				lcd.ShowMessage(VERSION+" Arrow key", LCD.LINE_2)
			} else {
				fmt.Println("keyboard event!", key, key.Char)
			}
		})
	}

	robot := gobot.NewRobot(
		"my-robot",
		[]gobot.Connection{r},
		[]gobot.Device{
			motorA,
			motorB,
			keys,
			servoKit.Driver,
			servoPan,
			servoTilt,
		},
		work,
	)

	robot.Start()
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "ip offline"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
