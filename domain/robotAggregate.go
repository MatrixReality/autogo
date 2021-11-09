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

func NewRobot(MessageBroker *infrastructure.MessageBroker, Motors *output.Motors, ServoKit *output.Servos, display *output.Display, SonarSet *input.SonarSet, Cfg *config.Config) *Robot {
	Status := &StatusDomain.Status{ColissionDetected: false, Direction: "", MinStopValue: Cfg.ArduinoSonar.MinStopValue}
	this := &Robot{MessageBroker: MessageBroker, Motors: Motors, ServoKit: ServoKit, Status: Status, Cfg: Cfg}

	if Cfg.ServoKit.Enabled {
		servoPan := ServoKit.GetByName("pan")
		servoTilt := ServoKit.GetByName("tilt")

		ServoKit.Init()
		ServoKit.SetCenter(servoPan)
		ServoKit.SetAngle(servoTilt, uint8(ServoKit.TiltPos["horizon"]))
	}

	var LCDDomain *LcdDomain.LCD
	if Cfg.LCD.Enabled {
		msgLine1 := getOutboundIP()
		if Cfg.Camera.Enabled {
			s := []string{msgLine1, Cfg.Camera.Port}
			msgLine1 = strings.Join(s, ":")
		}

		LCDTopic := fmt.Sprintf("autogo/%s/lcd", Cfg.RobotName)
		LCDDomain := LcdDomain.NewLCD(display, MessageBroker, LCDTopic)
		this.LCD = LCDDomain

		//TODO test only, remove
		if MessageBroker != nil {
			MessageBroker.Sub(LCDTopic)
		}

		LCDDomain.ShowMessage(msgLine1, 1)
		LCDDomain.ShowMessage(Cfg.Version+" Arrow key", 2)
	}

	if Cfg.ArduinoSonar.Enabled && Cfg.Motors.Enabled {
		sonarTopic := fmt.Sprintf("autogo/%s/sonar", Cfg.RobotName)
		sonarDomain := SonarDomain.NewSonarSet(SonarSet, LCDDomain, Motors, MessageBroker, Status, sonarTopic)
		this.SonarSet = sonarDomain

		//TODO test only, remove
		if MessageBroker != nil {
			MessageBroker.Sub(this.SonarSet.Topic)
		}

		go sonarDomain.SonarWorker()
	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	oldDirection := this.Status.Direction
	key := input.GetKeyEvent(data)

	if this.Cfg.ServoKit.Enabled {
		servoPan := this.ServoKit.GetByName("pan")
		servoTilt := this.ServoKit.GetByName("tilt")

		panAngle := int(servoPan.CurrentAngle)
		tiltAngle := int(servoTilt.CurrentAngle)

		if key.Key == keyboard.W {
			newTilt := tiltAngle - this.Cfg.ServoKit.PanTiltFactor
			if newTilt < this.ServoKit.TiltPos["top"] {
				newTilt = this.ServoKit.TiltPos["top"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.S {
			newTilt := tiltAngle + this.Cfg.ServoKit.PanTiltFactor
			if newTilt > this.ServoKit.TiltPos["down"] {
				newTilt = this.ServoKit.TiltPos["down"]
			}
			this.ServoKit.SetAngle(servoTilt, uint8(newTilt))

		} else if key.Key == keyboard.A {
			newPan := panAngle + this.Cfg.ServoKit.PanTiltFactor
			if newPan > this.ServoKit.PanPos["left"] {
				newPan = this.ServoKit.PanPos["left"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))

		} else if key.Key == keyboard.D {
			newPan := panAngle - this.Cfg.ServoKit.PanTiltFactor
			if newPan < this.ServoKit.PanPos["right"] {
				newPan = this.ServoKit.PanPos["right"]
			}
			this.ServoKit.SetAngle(servoPan, uint8(newPan))
		} else if key.Key == keyboard.X {
			this.ServoKit.SetCenter(servoPan)
			this.ServoKit.SetAngle(servoTilt, uint8(this.ServoKit.TiltPos["horizon"]))
		}
	}

	if this.Cfg.Motors.Enabled {
		if key.Key == keyboard.ArrowUp && this.Status.ColissionDetected == false {
			this.Motors.Forward(this.Cfg.Motors.MaxSpeed)
			this.Status.Direction = "Front"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowDown {
			this.Motors.Backward(this.Cfg.Motors.MaxSpeed)
			this.Status.Direction = "Back"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowRight {
			this.Motors.Left(this.Cfg.Motors.MaxSpeed)
			this.Status.Direction = "Right"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.ArrowLeft {
			this.Motors.Right(this.Cfg.Motors.MaxSpeed)
			this.Status.Direction = "Left"
			this.Status.LCDMsg = this.Status.Direction
		} else if key.Key == keyboard.Q {
			this.Motors.Stop()
			this.Status.Direction = ""
			this.Status.LCDMsg = this.Cfg.Version + " Arrow key"
		} else {
			fmt.Println(this.Status.LCDMsg, key, key.Char)
		}
	}

	if this.Cfg.LCD.Enabled && oldDirection != this.Status.Direction {
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
