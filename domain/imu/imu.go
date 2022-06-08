package domain

import (
	"fmt"
	"time"

	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	input "github.com/jtonynet/autogo/peripherals/input"
)

type IMU struct {
	IMU           *input.IMU
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	Topic         string
	Delay         time.Duration
}

func NewIMU(imu *input.IMU, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *IMU {
	delay, _ := time.ParseDuration(imu.Cfg.Delay)
	this := &IMU{
		IMU:           imu,
		MessageBroker: messageBroker,
		Status:        status,
		Topic:         topic,
		Delay:         delay,
	}

	imu.Init()
	time.Sleep(time.Second * 5)

	return this
}

func (this *IMU) Worker() {
	for true {
		this.IMU.Driver.GetData()

		fmt.Println("Accelerometer D:", this.IMU.Driver.Accelerometer)
		fmt.Println("Gyroscope D:", this.IMU.Driver.Gyroscope)
		fmt.Println("Temperature D:", this.IMU.Driver.Temperature)
		fmt.Println()
		fmt.Println()

		time.Sleep(this.Delay)
	}
}
