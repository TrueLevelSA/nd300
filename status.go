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
)

type Status struct {
	msg   string
	value byte
}

func (s Status) String() string {
	if s.msg == "" {
		return "unknown"
	}

	return s.msg
}

// nolint: gomnd // Numbers are from the ND-300KM manual.
var (
	PayoutSuccess   = Status{value: 0xAA, msg: "payout success"}       // Manual description: 'Payout Successful'.
	PayoutFailure   = Status{value: 0xBB, msg: "payout failure"}       // Manual description: 'Payout Fails'.
	StatusOk        = Status{value: 0x00, msg: "ok"}                   // Manual description: 'Status fine'.
	NoteEmpty       = Status{value: 0x01, msg: "note dispenser empty"} // Manual description: 'empty note'.
	NoteLow         = Status{value: 0x02, msg: "notes amount low"}     // Manual description: 'Stock less'.
	NoteJam         = Status{value: 0x03, msg: "note jam"}
	OverLength      = Status{value: 0x04, msg: "note over length"}
	NoteNotExit     = Status{value: 0x05, msg: "note not exited"}
	SensorError     = Status{value: 0x06, msg: "sensor error"}      // Manual description: 'Sensor error (Reserve)'.
	DoubleNoteError = Status{value: 0x07, msg: "double note error"} // Manual description: 'Double note error (Reserve)'.
	MotorError      = Status{value: 0x08, msg: "motor error"}
	DispensingBusy  = Status{value: 0x09, msg: "note dispenser busy"}
	SensorAdjusting = Status{value: 0x0A, msg: "sensor adjusting"} // Manual description: 'Sensor adjusting (Reserve)'.
	ChecksumError   = Status{value: 0x0B, msg: "checksum error"}
	LowPowerError   = Status{value: 0x0C, msg: "low power error"}

	ErrBadLen        = errors.New("bad message length")
	ErrBadSTX        = errors.New("bad STX prefix")
	ErrChecksum      = errors.New("invalid checksum")
	ErrUnknownStatus = errors.New("unknown status")
	ErrNotAStatus    = errors.New("not a status")
)

func (m Msg) ValidateAsStatus() error {
	if len(m) != MsgLen {
		return ErrBadLen
	}

	if m[idxStx] != STX {
		return ErrBadSTX
	}

	if m[idxMS] != StatusFlag {
		return ErrNotAStatus
	}

	if cksum := computeChecksum(m); cksum != m[idxCksum] {
		return fmt.Errorf("%w: expected %x, go %x", ErrChecksum, cksum, m[idxCksum])
	}

	switch m[idxCmd] {
	case PayoutSuccess.value,
		PayoutFailure.value,
		StatusOk.value,
		NoteEmpty.value,
		NoteLow.value,
		NoteJam.value,
		OverLength.value,
		NoteNotExit.value,
		SensorError.value,
		DoubleNoteError.value,
		MotorError.value,
		DispensingBusy.value,
		SensorAdjusting.value,
		ChecksumError.value,
		LowPowerError.value:
		return nil
	default:
		return ErrUnknownStatus
	}
}

func (m Msg) Status() Status {
	var status Status

	if m[idxMS] != StatusFlag {
		return Status{
			msg:   "not a status",
			value: m[idxStatus],
		}
	}

	switch m[idxStatus] {
	case PayoutSuccess.value:
		status = PayoutSuccess
	case PayoutFailure.value:
		status = PayoutFailure
	case StatusOk.value:
		status = StatusOk
	case NoteEmpty.value:
		status = NoteEmpty
	case NoteLow.value:
		status = NoteLow
	case NoteJam.value:
		status = NoteJam
	case OverLength.value:
		status = OverLength
	case NoteNotExit.value:
		status = NoteNotExit
	case SensorError.value:
		status = SensorError
	case DoubleNoteError.value:
		status = DoubleNoteError
	case MotorError.value:
		status = MotorError
	case DispensingBusy.value:
		status = DispensingBusy
	case SensorAdjusting.value:
		status = SensorAdjusting
	case ChecksumError.value:
		status = ChecksumError
	case LowPowerError.value:
		status = LowPowerError
	default:
		status = Status{
			value: m[idxStatus],
			msg:   "unknown status",
		}
	}

	return status
}
