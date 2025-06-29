package dto

import (
	"net/http"
	"multifinance/errors"
)

type ErrorResponseDto struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type SuccessResponseDto struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse creates a standard error response
func ErrorResponse(statusCode int, message string, err error) ErrorResponseDto {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return ErrorResponseDto{
		Code:    statusCode,
		Status:  http.StatusText(statusCode),
		Message: message,
		Error:   errMsg,
	}
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(validationErr *errors.ValidationErrors) ErrorResponseDto {
	return ErrorResponseDto{
		Code:    validationErr.Code,
		Status:  http.StatusText(validationErr.Code),
		Message: validationErr.Message,
		Errors:  validationErr.Errors,
	}
}

// SuccessResponse creates a standard success response
func SuccessResponse(statusCode int, message string, data interface{}) SuccessResponseDto {
	return SuccessResponseDto{
		Code:    statusCode,
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}
}
