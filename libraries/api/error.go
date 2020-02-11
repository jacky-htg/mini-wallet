package api

import "net/http"

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err           error
	MessageStatus string
	HTTPStatus    int
}

// ErrNew wraps a provided error with an HTTP status code and custome status code. This
// function should be used when handlers encounter expected errors.
func ErrNew(err error, messageStatus string, httpStatus int) error {
	return &Error{err, messageStatus, httpStatus}
}

// ErrBadRequest wraps a provided error with an HTTP status code and custome status code for bad request. This
// function should be used when handlers encounter expected errors.
func ErrBadRequest(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageBadRequest
	}
	return &Error{err, message, http.StatusBadRequest}
}

// ErrNotFound wraps a provided error with an HTTP status code and custome status code for not found. This
// function should be used when handlers encounter expected errors.
func ErrNotFound(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageNotFound
	}
	return &Error{err, message, http.StatusNotFound}
}

// ErrForbidden wraps a provided error with an HTTP status code and custome status code for forbidden. This
// function should be used when handlers encounter expected errors.
func ErrForbidden(err error, message string) error {
	if len(message) <= 0 || message == "" {
		message = StatusMessageForbidden
	}
	return &Error{err, message, http.StatusForbidden}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *Error) Error() string {
	return err.Err.Error()
}
