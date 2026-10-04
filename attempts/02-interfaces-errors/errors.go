package objects

import (
	"errors"
	"fmt"
)

// StatusReason mirrors metav1.StatusReason.
type StatusReason string

const (
	ReasonNotFound      StatusReason = "NotFound"
	ReasonAlreadyExists StatusReason = "AlreadyExists"
)

// StatusError mirrors k8s.io/apimachinery/pkg/api/errors.StatusError —
// the error you get back from client-go when a request fails.
type StatusError struct {
	Reason  StatusReason
	Code    int // HTTP status code
	Message string
}

// Error makes *StatusError satisfy the built-in `error` interface:
//
//	type error interface { Error() string }
//
// It should return the Message.
func (e *StatusError) Error() string {
	// TODO
	return e.Message
}

// NewNotFound returns a 404 NotFound error with message: <kind> "<name>" not found
// e.g. Pod "web" not found
func NewNotFound(kind, name string) *StatusError {
	// TODO
	return &StatusError{
		Reason:  ReasonNotFound,
		Code:    404,
		Message: fmt.Sprintf("%s %q not found", kind, name),
	}
}

// NewAlreadyExists returns a 409 AlreadyExists error with message: <kind> "<name>" already exists
func NewAlreadyExists(kind, name string) *StatusError {
	// TODO
	return &StatusError{
		Reason:  ReasonAlreadyExists,
		Code:    409,
		Message: fmt.Sprintf("%s %q already exists", kind, name),
	}
}

// ReasonForError returns the Reason if err is (or wraps) a *StatusError,
// otherwise "". Hint: errors.As — it looks through wrapped errors.
func ReasonForError(err error) StatusReason {
	// TODO
	var se *StatusError
	if errors.As(err, &se) {
		return se.Reason
	}
	return ""
}

// IsNotFound reports whether err is (or wraps) a NotFound StatusError.
// You'll write `if apierrors.IsNotFound(err)` constantly in real controllers.
func IsNotFound(err error) bool {
	// TODO
	// 	var se *StatusError
	// 	if errors.As(err, &se) {
	// 		return se.Code == 404
	// 	}
	// 	return false
	return ReasonForError(err) == ReasonNotFound
}

// IsAlreadyExists reports whether err is (or wraps) an AlreadyExists StatusError.
func IsAlreadyExists(err error) bool {
	// TODO
	return ReasonForError(err) == ReasonAlreadyExists
}
