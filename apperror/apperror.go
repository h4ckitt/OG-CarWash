package apperror

import "net/http"

type AppError struct {
	status  int
	message string
	err     error
}

var (
	Conflict            = AppError{status: http.StatusConflict, message: "A conflict occurred while processing request"}
	BadRequest          = AppError{status: http.StatusBadRequest, message: "Bad request body received"}
	ServerError         = AppError{status: http.StatusInternalServerError, message: "An error occurred while processing that request"}
	Forbidden           = AppError{status: http.StatusForbidden, message: "Forbidden"}
	NotFound            = AppError{status: http.StatusNotFound, message: "The requested resource was not found"}
	Unauthorized        = AppError{status: http.StatusUnauthorized, message: "Unauthorized"}
	UnprocessableEntity = AppError{status: http.StatusUnprocessableEntity, message: "Can't Process This Request"}
)

func (e *AppError) StatusAndMessage() (int, string) {
	return e.status, e.message
}

func (e *AppError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)

	if !ok {
		return false
	}

	if e.status == t.status {
		return true
	}

	return false
}

func (e *AppError) Wrap(err error) {
	e.err = err
}
