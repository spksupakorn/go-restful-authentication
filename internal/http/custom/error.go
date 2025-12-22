package custom

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnauthorizedError() error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
}

func NewValidationError(message string) error {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func NewUnexpctedError(message string) error {
	if message != "" {
		return &AppError{
			Code:    http.StatusInternalServerError,
			Message: message,
		}
	}

	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpeced error",
	}
}

func NewForrbidden(message string) error {
	if message != "" {
		return &AppError{
			Code:    http.StatusForbidden,
			Message: message,
		}
	}

	return &AppError{
		Code:    http.StatusForbidden,
		Message: "unauthorized",
	}
}
