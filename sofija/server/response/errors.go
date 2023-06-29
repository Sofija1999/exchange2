package response

import "net/http"

type UserError struct {
	Code      int    `json:"-"`
	Message   string `json:"error"`
	ErrorCode *int   `json:"error_code"`
	Internal    error  `json:"-"`
}

func NewValidationError(msg string) UserError {
	return UserError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NewForbiddenError(msg string) UserError {
	return UserError{
		Code:    http.StatusForbidden,
		Message: msg,
	}
}

func NewConflictError(msg string) UserError {
	return UserError{
		Code:    http.StatusConflict,
		Message: msg,
	}
}

func NewNotAcceptableError(msg string) UserError {
	return UserError{
		Code:    http.StatusNotAcceptable,
		Message: msg,
	}
}

func NewNotAcceptableForLegalReasonsError(msg string) UserError {
	return UserError{
		Code:    http.StatusUnavailableForLegalReasons,
		Message: msg,
	}
}

func NewPaymentRequiredError(msg string) UserError {
	return UserError{
		Code:    http.StatusPaymentRequired,
		Message: msg,
	}
}

func NewUnauthorizedError(msg string) UserError {
	return UserError{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

func NewTooEarlyError(msg string) UserError {
	return UserError{
		Code:    http.StatusTooEarly,
		Message: msg,
	}
}

func NewNotFoundError(msg string) UserError {
	return UserError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func NewInternalServerError(msg string) UserError {
	return UserError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

func NewServiceUnavailableError(msg string) UserError {
	return UserError{
		Code:    http.StatusServiceUnavailable,
		Message: msg,
	}
}

func NewUnsupportMediaTypeError(msg string) UserError {
	return UserError{
		Code:    http.StatusUnsupportedMediaType,
		Message: msg,
	}
}

func NewBadRequestError(msg string) UserError {
	return UserError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func (e UserError) Error() string {
	return e.Message
}

func (e UserError) WithInternal(err error) UserError {
	e.Internal = err

	return e
}

func (e UserError) WithErrorCode(errorCode int) UserError {
	e.ErrorCode = &errorCode

	return e
}
