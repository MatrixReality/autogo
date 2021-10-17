package main

// Circuit: esp8266-and-l298n-motor-controller
// Objective: dual speed and direction control using MotorDriver
//
// | Enable | Dir 1 | Dir 2 | Motor         |
// +--------+-------+-------+---------------+
// | 0      | X     | X     | Off           |
// | 1      | 0     | 0     | 0ff           |
// | 1      | 0     | 1     | On (forward)  |
// | 1      | 1     | 0     | On (backward) |
// | 1      | 1     | 1     | Off           |

import (
	"fmt"
	"log"
	"net"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	ArduinoSonarSet "github.com/matrixreality/autogo/peripherals/input"
	LCD "github.com/matrixreality/autogo/peripherals/output"
	Servos "github.com/matrixreality/autogo/peripherals/output"
)

/*
Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 12	   | Speed
A-Dir1        | DIR1 (Motor A) | 15	   | Direction
A-Dir2        | DIR2 (Motor A) | 11	   | Direction
B-Enable      | PWMA (Motor B) | 35	   | Speed
B-Dir1        | DIR1 (Motor B) | 16	   | Direction
B-Dir2        | DIR2 (Motor B) | 18	   | Direction
*/

const (
	maPWMPin  = "12"
	maDir1Pin = "15"
	maDir2Pin = "11"
	mbPWMPin  = "35"
	mbDir1Pin = "16"
	mbDir2Pin = "18"
)

//TODO env vars on viper
const (
	VERSION         = "v0.0.4"
	SERVOKIT_BUS    = 0
	SERVOKIT_ADDR   = 0x40
	ARDUINO_BUS     = 1
	ARDUINO_ADDR    = 0x18
	LCD_BUS         = 2
	LCD_ADDR        = 0x27
	LCD_COLLUMNS    = 16
	PAN_TILT_FACTOR = 30
)

const (
	maIndex = iota
	mbIndex
)

var (
	motorSpeed [2]byte
	motorInc   = [2]int{1, 1}
	counter    = [2]int{}
	motors     [2]*gpio.MotorDriver
)

func main() {
	r := raspi.NewAdaptor()
	keys := keyboard.NewDriver()

	///MOTORS
	motorA := gpio.NewMotorDriver(r, maPWMPin)
	motorA.ForwardPin = maDir1Pin
	motorA.BackwardPin = maDir2Pin
	motorA.SetName("Motor-A")

	motorB := gpio.NewMotorDriver(r, mbPWMPin)
	motorB.ForwardPin = mbDir1Pin
	motorB.BackwardPin = mbDir2Pin
	motorB.SetName("Motor-B")

	motors[maIndex] = motorA
	motors[mbIndex] = motorB
	///----

	///SERVOKIT
	servoKit := Servos.NewKitDriver(r, SERVOKIT_BUS, SERVOKIT_ADDR)
	servoPan := Servos.NewServo(servoKit, "0", "pan")
	servoTilt := Servos.NewServo(servoKit, "1", "tilt")

	///ARDUINO SONAR SET
	arduinoConn, err := ArduinoSonarSet.GetConnection(r, ARDUINO_BUS, ARDUINO_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	///LCD
	lcd, lcdClose, err := LCD.NewLcd(LCD_BUS, LCD_ADDR, LCD_COLLUMNS)
	if err != nil {
		log.Fatal(err)
	}
	defer lcdClose()

	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	ip := GetOutboundIP()

	err = lcd.ShowMessage(string(ip), LCD.LINE_1)
	if err != nil {
		log.Fatal(err)
	}

	err = lcd.ShowMessage(VERSION+" Arrow key", LCD.LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	firstRun := 1
	work := func() {
		servoKit.SetPWMFreq(60)
		if firstRun == 1 {
			firstRun = 0
			servoPan.Center()
			servoTilt.Move(uint8(Servos.TiltPos["horizon"]))
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
				servoTilt.Move(uint8(newTilt))

			} else if key.Key == keyboard.S {
				newTilt := tiltAngle + PAN_TILT_FACTOR
				if newTilt > Servos.TiltPos["down"] {
					newTilt = Servos.TiltPos["down"]
				}
				servoTilt.Move(uint8(newTilt))

			} else if key.Key == keyboard.A {
				newPan := panAngle + PAN_TILT_FACTOR
				if newPan > Servos.PanPos["left"] {
					newPan = Servos.PanPos["left"]
				}
				servoPan.Move(uint8(newPan))

			} else if key.Key == keyboard.D {
				newPan := panAngle - PAN_TILT_FACTOR
				if newPan < Servos.PanPos["right"] {
					newPan = Servos.PanPos["right"]
				}
				servoPan.Move(uint8(newPan))

			} else if key.Key == keyboard.X {
				servoPan.Center()
				servoTilt.Move(uint8(Servos.TiltPos["horizon"]))
			}

			if key.Key == keyboard.ArrowUp {
				motorA.Direction("forward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Front", " ", LCD_COLLUMNS), LCD.LINE_2)
			} else if key.Key == keyboard.ArrowDown {
				motorA.Direction("backward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Back", " ", LCD_COLLUMNS), LCD.LINE_2)
			} else if key.Key == keyboard.ArrowRight {
				motorA.Direction("forward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Left", " ", LCD_COLLUMNS), LCD.LINE_2)
			} else if key.Key == keyboard.ArrowLeft {
				motorA.Direction("backward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
				lcd.ShowMessage(rightPad("Right", " ", LCD_COLLUMNS), LCD.LINE_2)
			} else if key.Key == keyboard.Q {
				motorA.Speed(0)
				motorB.Speed(0)
				motorA.Direction("none")
				motorB.Direction("none")
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
			servoKit,
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

func rightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, (pLen-len(s)))
}
