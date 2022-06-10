// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
