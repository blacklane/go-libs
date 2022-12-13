package errors

import (
	"errors"
	"fmt"
)

const (
	DefaultErrorCode    = 500
	DefaultErrorMessage = "unknown_error"
)

type Error struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	Detail   string            `json:"detail,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	cause    error             `json:"-"`
}

func New(code int, message, detail string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func From(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return &Error{
		Code:    DefaultErrorCode,
		Message: DefaultErrorMessage,
		Detail:  err.Error(),
		cause:   err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error [%v]: %s, detail: %s, metadata: %v, cause: %v", e.Code, e.Message, e.Detail, e.Metadata, e.cause)
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Message == e.Message
	}
	return false
}

func (e *Error) WithCause(err error) *Error {
	ne := e.clone()
	ne.cause = err
	return ne
}

func (e *Error) WithDetail(detail string) *Error {
	err := e.clone()
	err.Detail = detail
	return err
}

func (e *Error) WithDetailf(format string, args ...any) *Error {
	err := e.clone()
	err.Detail = fmt.Sprintf(format, args...)
	return err
}

func (e *Error) WithMetadata(metadata map[string]string) *Error {
	err := e.clone()
	err.Metadata = metadata
	return err
}

func (e *Error) clone() *Error {
	metadata := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		metadata[k] = v
	}
	return &Error{
		Code:     e.Code,
		Message:  e.Message,
		Detail:   e.Detail,
		Metadata: metadata,
		cause:    e.cause,
	}
}
