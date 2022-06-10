// Copyright 2022 TrueLevel SA
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package nd300

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

const (
	DefaultBaudRate  = 9600
	DefaultParityBit = 8
	interReadDelay   = 5 * time.Millisecond
)

var SerialMode = &serial.Mode{
	BaudRate: DefaultBaudRate,
	DataBits: DefaultParityBit,
	Parity:   serial.EvenParity,
	StopBits: serial.OneStopBit,
}

type BadWriteSizeError int

func (e BadWriteSizeError) Error() string {
	return fmt.Sprintf("bad write size %d", e)
}

type Conn struct {
	port     serial.Port
	mode     *serial.Mode
	log      chan<- SerialMsg
	Log      <-chan SerialMsg
	portName string
	_rxBuff  []byte
	rxBuff   Msg
	txBuff   Msg
	timeout  time.Duration
}

func NewDefaultConn(port string, machineno byte, timeout time.Duration, logBuff int) *Conn {
	return NewConn(
		port,
		SerialMode,
		STX,
		machineno,
		MsgLen,
		timeout,
		logBuff,
	)
}

func NewConn(
	port string,
	mode *serial.Mode,
	stx byte,
	machineno byte,
	msgLen int,
	timeout time.Duration,
	logBuff int,
) *Conn {
	if mode == nil {
		mode = SerialMode
	}

	c := &Conn{
		port:     nil,
		mode:     mode,
		portName: port,
		_rxBuff:  make([]byte, msgLen),
		rxBuff:   make([]byte, msgLen),
		txBuff:   make([]byte, msgLen),
		timeout:  timeout,
	}

	c.txBuff[idxStx] = stx
	c.txBuff[idxMS] = CmdFlag
	c.txBuff[idxMachine] = machineno

	if logBuff >= 0 {
		log := make(chan SerialMsg, logBuff)
		c.log, c.Log = log, log
	}

	return c
}

func (c *Conn) Open() (serial.Port, error) { //nolint:ireturn // Return what serial.Open returns.
	if c.port != nil {
		return c.port, nil
	}

	var err error
	if c.port, err = serial.Open(c.portName, c.mode); err != nil {
		return c.port, fmt.Errorf("failed to open serial port: %w", err)
	}

	if err = c.port.SetReadTimeout(c.timeout); err != nil {
		return c.port, fmt.Errorf("failed to set timeout: %w", err)
	}

	return c.port, nil
}

func (c *Conn) read() error {
	var (
		size int
		err  error
		i    int
	)

	pos := 0

	for pos < len(c._rxBuff) {
		size, err = c.port.Read(c._rxBuff[pos:])
		if err != nil {
			return fmt.Errorf("failed to read msg: %w", err)
		}

		for i = 0; i < size; i++ {
			c.rxBuff[pos+i] = c._rxBuff[pos+i]
		}

		pos += size

		time.Sleep(interReadDelay) // Sadly no better option right now...
	}

	if c.log != nil {
		c.log <- SerialMsg{
			Type: RX,
			Data: c.rxBuff[:],
			Err:  err,
		}
	}

	return c.rxBuff.ValidateAsStatus()
}

func (c *Conn) write() error {
	size, err := c.port.Write(c.txBuff)

	switch {
	case err != nil:
		err = fmt.Errorf("failed to write msg: %w", err)
	case size != len(c.txBuff):
		err = fmt.Errorf("%w: should have written %d", BadWriteSizeError(size), len(c.txBuff))
	}

	if c.log != nil {
		c.log <- SerialMsg{
			Type: TX,
			Data: c.txBuff[:],
			Err:  err,
		}
	}

	return err
}
