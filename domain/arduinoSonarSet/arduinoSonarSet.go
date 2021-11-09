package domain

import (
	"encoding/json"
	"fmt"

	LcdDomain "github.com/jtonynet/autogo/domain/lcd"
	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
	output "github.com/jtonynet/autogo/peripherals/output"
)

type Sonar struct {
	SonarSet      *input.SonarSet
	Motors        *output.Motors //TODO: convert to domain
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	LCD           *LcdDomain.LCD
	Topic         string
}

/*
TODO: change output.Motors to domain.Motors in future
*/
func NewSonarSet(SonarSet *input.SonarSet, LCD *LcdDomain.LCD, Motors *output.Motors, MessageBroker *infrastructure.MessageBroker, Status *StatusDomain.Status, Topic string) *Sonar {
	this := &Sonar{SonarSet: SonarSet, LCD: LCD, Motors: Motors, MessageBroker: MessageBroker, Status: Status, Topic: Topic}
	return this
}

func (this *Sonar) SonarWorker() {
	for true {
		sonarData, err := this.SonarSet.GetData()
		if err == nil {
			if sonarData["center"] <= this.Status.MinStopValue && this.Status.Direction == "Front" && this.Status.ColissionDetected == false {
				this.Status.ColissionDetected = true
				this.Motors.Stop()

				if this.LCD != nil {
					s := fmt.Sprintf("STOP CRASH %.2f", sonarData["center"])
					this.LCD.ShowMessage(s, 2)
				}

			} else if this.Status.ColissionDetected && this.Status.Direction != "Front" {
				this.Status.ColissionDetected = false
			}

			if this.MessageBroker != nil {
				j, err := json.Marshal(sonarData)
				if err != nil {
					fmt.Printf("Error: %s", err.Error())
				} else {
					this.MessageBroker.Pub(this.SonarSet.Topic, string(j))
				}
			}

		}
	}
}
