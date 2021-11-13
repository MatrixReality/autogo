package domain

import (
	"fmt"
	"net"
	"strings"

	infrastructure "github.com/jtonynet/autogo/infrastructure"

	SonarDomain "github.com/jtonynet/autogo/domain/arduinoSonarSet"
	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	LocomotionDomain "github.com/jtonynet/autogo/domain/locomotion"
	domain "github.com/jtonynet/autogo/domain/locomotion"
	ServosDomain "github.com/jtonynet/autogo/domain/servos"
	StatusDomain "github.com/jtonynet/autogo/domain/status"

	config "github.com/jtonynet/autogo/config"

	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

type Robot struct {
	MessageBroker *infrastructure.MessageBroker

	LCD        *LcdDomain.LCD
	Locomotion *LocomotionDomain.Locomotion
	Servos     *ServosDomain.Servos
	SonarSet   *SonarDomain.Sonar
	Status     *StatusDomain.Status

	Cfg *config.Config
}

func NewRobot(messageBroker *infrastructure.MessageBroker, motors *output.Motors, servos *output.Servos, display *output.Display, sonarSet *input.SonarSet, cfg *config.Config) *Robot {
	Status := &StatusDomain.Status{
		ColissionDetected: false,
		Direction:         "",
		Version:           cfg.Version,
		ProjectName:       cfg.ProjectName,
		RobotName:         cfg.RobotName,
		MinStopValue:      cfg.ArduinoSonar.MinStopValue,
	}

	this := &Robot{MessageBroker: messageBroker, Status: Status, Cfg: cfg}

	if servos != nil {
		servosDomain := ServosDomain.NewServos(servos)
		this.Servos = servosDomain
	}

	if display != nil {
		msgLine1 := getOutboundIP()
		if cfg.Camera.Enabled {
			s := []string{msgLine1, cfg.Camera.Port}
			msgLine1 = strings.Join(s, ":")
		}

		LCDTopic := fmt.Sprintf("%s/%s/lcd", cfg.ProjectName, cfg.RobotName)
		this.LCD = LcdDomain.NewLCD(display, messageBroker, LCDTopic)

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			messageBroker.Sub(LCDTopic)
		}

		this.LCD.ShowMessage(msgLine1, 1)
		this.LCD.ShowMessage(cfg.Version+" Arrow key", 2)
	}

	if motors != nil {
		locomotionDomain := domain.NewLocomotion(motors, this.LCD, this.Status)
		this.Locomotion = locomotionDomain
	}

	if sonarSet != nil {
		sonarTopic := fmt.Sprintf("%s/%s/sonar", cfg.ProjectName, cfg.RobotName)
		sonarDomain := SonarDomain.NewSonarSet(sonarSet, this.LCD, this.Locomotion, messageBroker, Status, sonarTopic)
		this.SonarSet = sonarDomain

		//TODO: Test only, remove after create robot client subscription
		if messageBroker != nil {
			messageBroker.Sub(sonarTopic)
		}

		go sonarDomain.SonarWorker()
	}

	return this
}

func (this *Robot) ControllByKeyboard(data interface{}) {
	key := input.GetKeyEvent(data)

	if this.Servos != nil {
		go this.Servos.ControllPanAndTilt(key.Key)
	}

	if this.Locomotion != nil {
		go this.Locomotion.ControllMoviment(key.Key)
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
