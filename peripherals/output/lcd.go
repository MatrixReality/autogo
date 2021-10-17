package peripherals

import (
	"log"

	deviceD2r2 "github.com/d2r2/go-hd44780"
	i2cD2r2 "github.com/d2r2/go-i2c"
)

const (
	DEFAULT_COLLUMNS = 16

	LINE_1 = deviceD2r2.SHOW_LINE_1
	LINE_2 = deviceD2r2.SHOW_LINE_2
	LINE_3 = deviceD2r2.SHOW_LINE_3
	LINE_4 = deviceD2r2.SHOW_LINE_4
)

func NewLcd(bus int, addr uint8, collumns int) (*deviceD2r2.Lcd, func(), error) {
	//TODO: use lcd i2c gobot solution to 16x2 screen

	// Create new connection to i2c-bus on 2 line with address 0x27.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2cD2r2.NewI2C(addr, bus)
	if err != nil {
		log.Fatal(err)
	}

	// Construct lcd-device connected via I2C connection
	lcd, err := deviceD2r2.NewLcd(i2c, deviceD2r2.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}

	if collumns == 0 {
		collumns = DEFAULT_COLLUMNS
	}

	// Turn on the backlight
	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}

	return lcd,
		func() {
			// Free I2C connection on exit
			defer i2c.Close()
		},
		nil
}
