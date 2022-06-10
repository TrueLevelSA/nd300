package nd300

import (
	"go.bug.st/serial"
)

func (c *Conn) RequestStatus() (status Status, count byte, err error) {
	var port serial.Port

	port, err = c.Open()
	if err != nil {
		return
	}

	defer func() {
		err = AppendErr(err, port.Close(), closeErrMsg)
		c.port = nil
	}()

	status, count, err = requestStatus(c)

	return
}

func requestStatus(conn *Conn) (status Status, count byte, err error) {
	//nolint:goconst // False positive, see: https://github.com/jgautheron/goconst/issues/19
	if conn.txBuff[idxCmd] != byte(RequestMachineStatus) || conn.txBuff[idxData] != 0x0 {
		conn.txBuff[idxCmd] = byte(RequestMachineStatus)
		conn.txBuff[idxData] = 0x0
		conn.txBuff[idxCksum] = computeChecksum(conn.txBuff)
	}

	if err = conn.write(); err != nil {
		return
	}

	if err = conn.read(); err != nil {
		return
	}

	count = conn.rxBuff[idxData]
	status = conn.rxBuff.Status()

	return
}
