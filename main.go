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

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/raspi"
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

	work := func() {
		/*
			motorA.Direction("forward")
			motorB.Direction("backward")

			gobot.Every(40*time.Millisecond, func() {
				motorControl(maIndex)
			})

			gobot.Every(20*time.Millisecond, func() {
				motorControl(mbIndex)
			})
		*/
		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == keyboard.W {
				motorA.Direction("forward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
			} else if key.Key == keyboard.S {
				motorA.Direction("backward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
			} else if key.Key == keyboard.A {
				motorA.Direction("forward")
				motorB.Direction("backward")
				motorA.Speed(255)
				motorB.Speed(255)
			} else if key.Key == keyboard.D {
				motorA.Direction("backward")
				motorB.Direction("forward")
				motorA.Speed(255)
				motorB.Speed(255)
			} else if key.Key == keyboard.Q {
				motorA.Speed(0)
				motorB.Speed(0)
			} else {
				fmt.Println("keyboard event!", key, key.Char)
			}
		})
	}

	robot := gobot.NewRobot(
		"my-robot",
		[]gobot.Connection{r},
		[]gobot.Device{motorA, motorB, keys},
		work,
	)

	robot.Start()
}
/*
func motorControl(idx int) {
	m := motors[idx]

	motorSpeed[idx] = byte(int(motorSpeed[idx]) + motorInc[idx])
	fmt.Println(motorSpeed[idx])
	m.Speed(motorSpeed[idx])

	counter[idx]++
	if counter[idx]%256 == 255 {
		if motorInc[idx] == 1 {
			motorInc[idx] = 0
		} else if motorInc[idx] == 0 {
			motorInc[idx] = -1
		} else {
			motorInc[idx] = 1
		}
	}

	if counter[idx]%766 == 765 {
		if m.CurrentDirection == "forward" {
			m.Direction("backward")
		} else {
			m.Direction("forward")
		}
	}
*/
}
