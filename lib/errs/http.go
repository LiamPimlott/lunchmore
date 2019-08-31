package errs

import "net/http"

var (
	// ErrNotFound represents a not found http response.
	ErrNotFound = &notFoundError{Code: http.StatusNotFound, Msg: "Not Found"}

	// ErrInvalid represents an invalid header http.
	ErrInvalid = &invalidError{Code: http.StatusBadRequest, Msg: "Bad Request"}

	// ErrInvalidAuth represents an invalid authorization header http response.
	ErrInvalidAuth = &invalidAuthorizationError{Code: http.StatusBadRequest, Msg: "Invalid Authorization Header"}

	// ErrForbidden represents an http forbidden response.
	ErrForbidden = &invalidAuthorizationError{Code: http.StatusForbidden, Msg: "Unauthorized"}

	// ErrInternal represents an http internal server error response.
	ErrInternal = &invalidAuthorizationError{Code: http.StatusInternalServerError, Msg: "Internal Server Error"}
)

type notFoundError struct {
	Code int
	Msg  string
}

func (e *notFoundError) Error() string { return e.Msg }

type invalidError struct {
	Code int
	Msg  string
}

func (e *invalidError) Error() string { return e.Msg }

type invalidAuthorizationError struct {
	Code int
	Msg  string
}

func (e *invalidAuthorizationError) Error() string { return e.Msg }

type forbiddenError struct {
	Code int
	Msg  string
}

func (e *forbiddenError) Error() string { return e.Msg }

type internalError struct {
	Code int
	Msg  string
}

func (e *internalError) Error() string { return e.Msg }
