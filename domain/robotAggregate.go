package domain

import (
	"fmt"
	"net"
	"strings"

	config "github.com/jtonynet/autogo/config"
	SonarDomain "github.com/jtonynet/autogo/domain/arduinoSonarSet"
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"

	"gobot.io/x/gobot/platforms/keyboard"
)

type Robot struct {
	MessageBroker *infrastructure.MessageBroker

	Motors   *output.Motors
	ServoKit *output.Servos

	Cfg *config.Config

	LCD      *LcdDomain.LCD
	SonarSet *SonarDomain.Sonar
	Status   *StatusDomain.Status
}

func NewRobot(messageBroker *infrastructure.MessageBroker, motors *output.Motors, servoKit *output.Servos, display *output.Display, sonarSet *input.SonarSet, cfg *config.Config) *Robot {
	Status := &StatusDomain.Status{ColissionDetected: false, Direction: "", MinStopValue: cfg.ArduinoSonar.MinStopValue}
	this := &Robot{MessageBroker: messageBroker, Motors: motors, ServoKit: servoKit, Status: Status, Cfg: cfg}

	if servoKit != nil {
		servoPan := servoKit.GetByName("pan")
		servoTilt := servoKit.GetByName("tilt")

		servoKit.Init()
		servoKit.SetCenter(servoPan)
		servoKit.SetAngle(servoTilt, uint8(servoKit.TiltPos["horizon"]))
	}

	if display != nil {
		msgLine1 := getOutboundIP()
		if cfg.Camera.Enabled {
			s := []string{msgLine1, cfg.Camera.Port}
			msgLine1 = strings.Join(s, ":")
		}

		LCDTopic := fmt.Sprintf("autogo/%s/lcd", cfg.RobotName)
		this.LCD = LcdDomain.NewLCD(display, messageBroker, LCDTopic)

		//TODO test only, remove
		if messageBroker != nil {
			messageBroker.Sub(LCDTopic)
		}

		this.LCD.ShowMessage(msgLine1, 1)
		this.LCD.ShowMessage(cfg.Version+" Arrow key", 2)
	}

	if sonarSet != nil && motors != nil {
		//fmt.Print("cfg.ArduinoSonar.DelayInMS: ")
		//fmt.Print(cfg.ArduinoSonar.DelayInMS)
		//fmt.Print("------------------\n\n")

		sonarTopic := fmt.Sprintf("autogo/%s/sonar", cfg.RobotName)
		sonarDomain := SonarDomain.NewSonarSet(sonarSet, this.LCD, motors, messageBroker, Status, sonarTopic, cfg.ArduinoSonar.DelayInMS)
		this.SonarSet = sonarDomain

		//TODO test only, remove
		if messageBroker != nil {
			messageBroker.Sub(sonarTopic)
		}

		go sonarDomain.SonarWorker()
	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	oldDirection := this.Status.Direction
	key := input.GetKeyEvent(data)
	cfg := this.Cfg

	if this.ServoKit != nil {
		servoPan := this.ServoKit.GetByName("pan")
		servoTilt := this.ServoKit.GetByName("tilt")

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)

		if key.Key == keyboard.W {
			newTilt := tiltAngle - cfg.ServoKit.PanTiltFactor
			if newTilt < this.ServoKit.TiltPos["top"] {
				newTilt = this.ServoKit.TiltPos["top"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + cfg.ServoKit.PanTiltFactor
			if newTilt > this.ServoKit.TiltPos["down"] {
				newTilt = this.ServoKit.TiltPos["down"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + cfg.ServoKit.PanTiltFactor
			if newPan > this.ServoKit.PanPos["left"] {
				newPan = this.ServoKit.PanPos["left"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - cfg.ServoKit.PanTiltFactor
			if newPan < this.ServoKit.PanPos["right"] {
				newPan = this.ServoKit.PanPos["right"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))
		} else if key.Key == keyboard.X {
			this.ServoKit.SetCenter(servoPan)
			this.ServoKit.SetAngle(servoTilt, uint8(this.ServoKit.TiltPos["horizon"]))
		}
	}

	if this.Motors != nil {
		if key.Key == keyboard.ArrowUp && this.Status.ColissionDetected == false {
			this.Motors.Forward(cfg.Motors.MaxSpeed)
			this.Status.Direction = "Front"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowDown {
			this.Motors.Backward(cfg.Motors.MaxSpeed)
			this.Status.Direction = "Back"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowRight {
			this.Motors.Left(cfg.Motors.MaxSpeed)
			this.Status.Direction = "Right"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowLeft {
			this.Motors.Right(cfg.Motors.MaxSpeed)
			this.Status.Direction = "Left"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.Q {
			this.Motors.Stop()
			this.Status.Direction = ""
			this.Status.LCDMsg = cfg.Version + " Arrow key"
		} else {
			fmt.Println(this.Status.LCDMsg, key, key.Char)
		}
	}

	if this.LCD != nil && oldDirection != this.Status.Direction {
		this.LCD.ShowMessage(this.Status.LCDMsg, 2)
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "offline"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
