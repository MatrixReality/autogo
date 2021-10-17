package peripherals

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func GetConnection(a *raspi.Adaptor, bus int, addr uint8) (connection i2c.Connection, err error) {
	arduinoConn, err := a.GetConnection(0x18, 1)
	if err != nil {
		return nil, err
	}

	return arduinoConn, nil
}

func GetData(conn i2c.Connection) (string, error) {
	_, err := conn.Write([]byte("A"))
	if err != nil {
		return "", err
	}

	sonarByteLen := 28
	buf := make([]byte, sonarByteLen)
	bytesRead, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	sonarData := ""
	if bytesRead == sonarByteLen {
		sonarData = string(buf[:])
	}

	return sonarData, nil
}
