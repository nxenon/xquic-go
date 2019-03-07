package qerr

import (
	"fmt"
	"net"
)

// A QuicError consists of an error code plus a error reason
type QuicError struct {
	ErrorCode    ErrorCode
	ErrorMessage string
	isTimeout    bool
}

var _ net.Error = &QuicError{}

// Error creates a new QuicError instance
func Error(errorCode ErrorCode, errorMessage string) *QuicError {
	return &QuicError{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}

// TimeoutError creates a new QuicError instance for a timeout error
func TimeoutError(errorMessage string) *QuicError {
	return &QuicError{
		ErrorMessage: errorMessage,
		isTimeout:    true,
	}
}

// CryptoError create a new QuicError instance for a crypto error
func CryptoError(tlsAlert uint8) *QuicError {
	return &QuicError{
		ErrorCode: 0x100 + ErrorCode(tlsAlert),
	}
}

func (e *QuicError) Error() string {
	if len(e.ErrorMessage) == 0 {
		return e.ErrorCode.String()
	}
	return fmt.Sprintf("%s: %s", e.ErrorCode.String(), e.ErrorMessage)
}

// Temporary says if the error is temporary.
func (e *QuicError) Temporary() bool {
	return false
}

// Timeout says if this error is a timeout.
func (e *QuicError) Timeout() bool {
	return e.isTimeout
}

// ToQuicError converts an arbitrary error to a QuicError. It leaves QuicErrors
// unchanged, and properly handles `ErrorCode`s.
func ToQuicError(err error) *QuicError {
	switch e := err.(type) {
	case *QuicError:
		return e
	case ErrorCode:
		return Error(e, "")
	}
	return Error(InternalError, err.Error())
}
