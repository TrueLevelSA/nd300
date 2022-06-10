package nd300

import (
	"errors"
	"fmt"
	"time"

	"go.bug.st/serial"
)

const (
	resetDelay  = 500 * time.Millisecond // Based on manual testing.
	resetChecks = 8
)

var ErrResetFailed = errors.New("machine reset failed")

func (c *Conn) ResetMachine() (err error) {
	var port serial.Port

	port, err = c.Open()
	if err != nil {
		return
	}

	defer func() {
		err = AppendErr(err, port.Close(), closeErrMsg)
		c.port = nil
	}()

	if c.txBuff[idxCmd] != byte(ResetDispenser) || c.txBuff[idxData] != 0x0 {
		c.txBuff[idxCmd] = byte(ResetDispenser)
		c.txBuff[idxData] = 0x0
		c.txBuff[idxCksum] = computeChecksum(c.txBuff)
	}

	if err = c.write(); err != nil {
		return
	}

	checks := 0

	// Whilst reseting the machine must be adjusting the sensor and returns SensorAdjusting.
	for status := SensorAdjusting; status == SensorAdjusting; checks++ {
		time.Sleep(resetDelay)

		status, _, err = requestStatus(c)
		if err != nil {
			return
		}

		if checks >= resetChecks {
			return fmt.Errorf("%w: %s", ErrResetFailed, status)
		}
	}

	return nil
}
