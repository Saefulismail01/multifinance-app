package dto

import (
	"net/http"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type validationError struct {
	Errors []ValidationError
}

// Response represents a standard response
type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Success",
		Data:    data,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(errs []ValidationError) error {
	return &validationError{
		Errors: errs,
	}
}

// Error implements the error interface
func (v *validationError) Error() string {
	return "validation failed"
}

// GetErrors returns the list of validation errors
func (v *validationError) GetErrors() []ValidationError {
	return v.Errors
}

// SuccessResponse creates a standard success response
func SuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Success",
		Data:    data,
	}
}
