// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package nd300

const (
	MsgLen     = 6 // Length of a serial message (command or status).
	idxStx     = 0
	idxMS      = 1 // Indicate if master (command) or slave (status).
	idxMachine = 2
	idxCmd     = 3
	idxStatus  = 3 // Same as idxCmd, helps make the code clearer.
	idxData    = 4
	idxCksum   = 5

	STX        byte = 0x01
	CmdFlag    byte = 0x10
	StatusFlag byte = 0x01
)

type Msg []byte

func (m Msg) SetData(data byte) {
	m[idxData] = data
	m[idxCksum] = computeChecksum(m)
}

func (m Msg) Bytes() []byte {
	return m
}

func (m Msg) CmdOrStatus() string {
	switch m[idxMS] {
	case CmdFlag:
		return m.CmdString()
	case StatusFlag:
		return m.Status().String()
	default:
		return "invalid msg"
	}
}

func computeChecksum(msg []byte) byte {
	return msg[idxStx] + msg[idxMS] + msg[idxMachine] + msg[idxCmd] + msg[idxData]
}
