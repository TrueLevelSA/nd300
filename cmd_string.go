// Code generated by "stringer -linecomment -type Cmd ./cmd.go"; DO NOT EDIT.

package nd300

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SingleMachinePayout-16]
	_ = x[RequestMachineStatus-17]
	_ = x[ResetDispenser-18]
	_ = x[MultipleMachinesPayout-19]
}

const _Cmd_name = "machine payoutmachine statusmachine resetmultiple machines payout"

var _Cmd_index = [...]uint8{0, 14, 28, 41, 65}

func (i Cmd) String() string {
	i -= 16
	if i >= Cmd(len(_Cmd_index)-1) {
		return "Cmd(" + strconv.FormatInt(int64(i+16), 10) + ")"
	}
	return _Cmd_name[_Cmd_index[i]:_Cmd_index[i+1]]
}
