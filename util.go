package nd300

import (
	"fmt"
)

// AppendErr appends a new error to an existing one with an optional message.
// If both are nil, nil is returned.
// If the existing error only is nil, then the new error is returned.
// If the new error only is nil, then the existing error is returned.
// If both the existing error and the new are not nil, return an error of the form:
// "<existingErr>, <message>: <newErr>"
// or "<existingErr>, <newErr>" if `message` is an empty string.
// In this case, only the existing error can be unwrapped.
func AppendErr(existingErr error, newErr error, message string) error {
	switch {
	case existingErr != nil && newErr != nil && message != "":
		existingErr = fmt.Errorf("%w, %s: %v", existingErr, message, newErr)
	case existingErr != nil && newErr != nil:
		existingErr = fmt.Errorf("%w, %v", existingErr, newErr)
	case newErr != nil:
		existingErr = newErr
	}

	return existingErr
}
