package nd300

import (
	"errors"
	"fmt"

	"go.bug.st/serial"
)

var (
	ErrPayoutFailed  = errors.New("payout failed")
	ErrPayoutWarning = errors.New("payout warning")
)

func (c *Conn) Payout(notes byte) (count byte, err error) {
	if notes == 0 {
		return
	}

	var port serial.Port

	port, err = c.Open()
	if err != nil {
		return
	}

	defer func() {
		err = AppendErr(err, port.Close(), closeErrMsg)
		c.port = nil
	}()

	if c.txBuff[idxCmd] != byte(SingleMachinePayout) || c.txBuff[idxData] != notes {
		c.txBuff[idxCmd] = byte(SingleMachinePayout)
		c.txBuff[idxData] = notes
		c.txBuff[idxCksum] = computeChecksum(c.txBuff)
	}

	if err = c.write(); err != nil {
		return
	}

	for count < notes {
		err = c.read()
		if err != nil {
			return
		}

		switch c.rxBuff.Status() {
		case PayoutSuccess:
			count = c.rxBuff[idxData]
		case PayoutFailure:
			count, err = c.rxBuff[idxData], ErrPayoutFailed
			return
		case StatusOk:
			if c.rxBuff[idxData] == notes {
				return
			}
		case DispensingBusy, SensorAdjusting:
			continue

		case ChecksumError, NoteLow:
			// Note the error but continue processing as error is recoverable.
			err = fmt.Errorf("%w: %s", ErrPayoutWarning, c.rxBuff.Status())

		case NoteEmpty:
			if count == notes {
				err = fmt.Errorf("%w: %s", ErrPayoutWarning, c.rxBuff.Status())
			} else {
				err = fmt.Errorf("%w: %s", ErrPayoutFailed, c.rxBuff.Status())
			}

			return
		case NoteJam,
			OverLength,
			NoteNotExit,
			SensorError,
			DoubleNoteError,
			MotorError,
			LowPowerError:
			err = fmt.Errorf("%w: %s", ErrPayoutFailed, c.rxBuff.Status())

			return
		}
	}

	return
}
