// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package nd300

// type Port interface {
// 	io.ReadWriteCloser
//
// 	SetReadTimeout(duration time.Duration) error
// }

type SerialType bool

const (
	RX          SerialType = false
	TX          SerialType = true
	closeErrMsg            = "failed to close serial port"
)

type SerialMsg struct {
	Err  error
	Data Msg
	Type SerialType
}
