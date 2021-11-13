package domain

import (
	"encoding/json"
	"fmt"
	"time"

	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	locomotionDomain "github.com/jtonynet/autogo/domain/locomotion"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
)

type Sonar struct {
	SonarSet      *input.SonarSet
	Locomotion    *locomotionDomain.Locomotion
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	LCD           *LcdDomain.LCD
	Topic         string
}

//TODO: Change output.Motors to domain.Motors in future
func NewSonarSet(sonarSet *input.SonarSet, LCD *LcdDomain.LCD, locomotion *locomotionDomain.Locomotion, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *Sonar {
	this := &Sonar{SonarSet: sonarSet, LCD: LCD, Locomotion: locomotion, MessageBroker: messageBroker, Status: status, Topic: topic}
	return this
}

func (this *Sonar) sendDataToMessageBroker(sonarData map[string]float64) {
	j, err := json.Marshal(sonarData)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		this.MessageBroker.Pub(this.Topic, string(j))
	}
}

func (this *Sonar) SonarWorker() {
	status := this.Status
	delayInMS, _ := time.ParseDuration(
		fmt.Sprintf("%vms", this.SonarSet.Cfg.DelayInMS))

	for true {
		sonarData, err := this.SonarSet.GetData()
		if err != nil {
			return
		}

		if sonarData["center"] <= status.MinStopValue && status.Direction == "Front" && status.ColissionDetected == false {
			status.ColissionDetected = true

			if this.Locomotion != nil {
				this.Locomotion.Stop()
			}

			if this.LCD != nil {
				s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
				this.LCD.ShowMessage(s, 2)
			}

		} else if status.ColissionDetected && status.Direction != "Front" {
			status.ColissionDetected = false
		}

		if this.MessageBroker != nil {
			go this.sendDataToMessageBroker(sonarData)
		}

		status.SonarData = sonarData
		time.Sleep(delayInMS)
	}
}
