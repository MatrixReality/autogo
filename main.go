package main

import (
	"log"
	"net"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"

	handlers "github.com/matrixreality/autogo/handlers"
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
	motors := Motors.NewMotors(r)

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
	work := func() {
		handlers.InitKeyboard(servoKit, arduinoConn, lcd, motors, keys)
	}

	robot := gobot.NewRobot(
		"my-robot",
		[]gobot.Connection{r},
		[]gobot.Device{
			motors.MotorA,
			motors.MotorB,
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
