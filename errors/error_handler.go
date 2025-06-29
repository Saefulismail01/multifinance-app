package errors

import "net/http"

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
	}
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return e.Message
}

// Common error responses
var (
	ErrBadRequest          = NewErrorResponse(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = NewErrorResponse(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = NewErrorResponse(http.StatusForbidden, "forbidden")
	ErrNotFound            = NewErrorResponse(http.StatusNotFound, "not found")
	ErrInternalServerError = NewErrorResponse(http.StatusInternalServerError, "internal server error")
)

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Code    int               `json:"code"`
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// NewValidationError creates a new validation error response
func NewValidationError(message string, errors []ValidationError) *ValidationErrors {
	return &ValidationErrors{
		Code:    http.StatusBadRequest,
		Status:  http.StatusText(http.StatusBadRequest),
		Message: message,
		Errors:  errors,
	}
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	return ve.Message
}

// Add adds a new validation error
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors checks if there are any validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}
