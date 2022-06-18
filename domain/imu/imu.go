package domain

import (
	"encoding/json"
	"time"

	StatusDomain "github.com/jtonynet/autogo/domain/status"
	infrastructure "github.com/jtonynet/autogo/infrastructure"
	sensors "github.com/jtonynet/autogo/peripherals/sensors"
	"gobot.io/x/gobot/drivers/i2c"
)

type IMU struct {
	IMU           *sensors.IMU
	MessageBroker *infrastructure.MessageBroker
	Status        *StatusDomain.Status
	Topic         string
	Delay         time.Duration
}

//TODO cast ThreeDData to remove 12c dependence
type IMUMessage struct {
	Accel i2c.ThreeDData
	Gyro  i2c.ThreeDData
	Temp  int16
}

func NewIMU(imu *sensors.IMU, messageBroker *infrastructure.MessageBroker, status *StatusDomain.Status, topic string) *IMU {
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
		this.IMU.GetData()

		//fmt.Println("Model", this.IMU.GetModel())
		//fmt.Println("Accelerometer", this.IMU.GetAccelerometer())
		//fmt.Println("Gyroscope", this.IMU.GetGyroscope())
		//fmt.Println("Temperature", this.IMU.GetTemperature())

		if this.MessageBroker != nil {
			m := IMUMessage{
				this.IMU.GetAccelerometer(),
				this.IMU.GetGyroscope(),
				this.IMU.GetTemperature(),
			}

			m_marshalled, _ := json.Marshal(m)

			this.MessageBroker.Pub(this.Topic, string(m_marshalled))
		}

		time.Sleep(this.Delay)
	}
}
