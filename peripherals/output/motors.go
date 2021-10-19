package peripherals

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

/*
// Objective: dual speed and direction control using MotorDriver
//
// | Enable | Dir 1 | Dir 2 | Motor         |
// +--------+-------+-------+---------------+
// | 0      | X     | X     | Off           |
// | 1      | 0     | 0     | 0ff           |
// | 1      | 0     | 1     | On (forward)  |
// | 1      | 1     | 0     | On (backward) |
// | 1      | 1     | 1     | Off           |

Motor Shield  | NodeMCU        | GPIO  | Purpose
--------------+----------------+-------+----------
A-Enable      | PWMA (Motor A) | 12	   | Speed
A-Dir1        | DIR1 (Motor A) | 15	   | Direction
A-Dir2        | DIR2 (Motor A) | 11	   | Direction
B-Enable      | PWMA (Motor B) | 35	   | Speed
B-Dir1        | DIR1 (Motor B) | 16	   | Direction
B-Dir2        | DIR2 (Motor B) | 18	   | Direction
*/

//TODO env vars on viper
const (
	maPWMPin  = "12"
	maDir1Pin = "15"
	maDir2Pin = "11"
	mbPWMPin  = "35"
	mbDir1Pin = "16"
	mbDir2Pin = "18"
)

var (
	motorSpeed [2]byte
	motorInc   = [2]int{1, 1}
	counter    = [2]int{}
	motors     [2]*gpio.MotorDriver
)

const (
	maIndex = iota
	mbIndex
)

var motorA *gpio.MotorDriver
var motorB *gpio.MotorDriver

func NewMotors(a *raspi.Adaptor) (*gpio.MotorDriver, *gpio.MotorDriver) {
	motorA = gpio.NewMotorDriver(a, maPWMPin)
	motorA.ForwardPin = maDir1Pin
	motorA.BackwardPin = maDir2Pin
	motorA.SetName("Motor-A")

	motorB = gpio.NewMotorDriver(a, mbPWMPin)
	motorB.ForwardPin = mbDir1Pin
	motorB.BackwardPin = mbDir2Pin
	motorB.SetName("Motor-B")

	motors[maIndex] = motorA
	motors[mbIndex] = motorB

	return motorA, motorB
}

func Forward(speed byte) {
	motorA.Direction("forward")
	motorB.Direction("forward")
	motorA.Speed(speed)
	motorB.Speed(speed)
}

func Backward(speed byte) {
	motorA.Direction("backward")
	motorB.Direction("backward")
	motorA.Speed(speed)
	motorB.Speed(speed)
}

func Right(speed byte) {
	motorA.Direction("forward")
	motorB.Direction("backward")
	motorA.Speed(speed)
	motorB.Speed(speed)
}

func Left(speed byte) {
	motorA.Direction("backward")
	motorB.Direction("forward")
	motorA.Speed(speed)
	motorB.Speed(speed)
}

func Stop() {
	motorA.Speed(0)
	motorB.Speed(0)
	motorA.Direction("none")
	motorB.Direction("none")
}
