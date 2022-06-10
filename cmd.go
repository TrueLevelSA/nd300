package nd300

// nd300 is a driver for the ND-300 note dispenser.
// It implements the necessary functions to interact over a serial connection.

//go:generate stringer -linecomment -type Cmd ./cmd.go

//nolint:godot // Comments are used by stringer for string representation.
const (
	SingleMachinePayout    Cmd = 0x10 // machine payout
	RequestMachineStatus   Cmd = 0x11 // machine status
	ResetDispenser         Cmd = 0x12 // machine reset
	MultipleMachinesPayout Cmd = 0x13 // multiple machines payout
)

type Cmd byte

func (m Msg) CmdString() string {
	if m[idxMS] != CmdFlag {
		return "not a command"
	}

	cmd := Cmd(m[idxCmd])
	switch cmd {
	case SingleMachinePayout, RequestMachineStatus, ResetDispenser, MultipleMachinesPayout:
		return cmd.String()
	default:
		return "unknown command"
	}
}
