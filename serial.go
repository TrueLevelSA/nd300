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
