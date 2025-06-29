package service

import (
	"net/http"
	"multifinance/delivery/dto"
	"multifinance/errors"
)

// ValidateService defines the interface for validation services
type ValidateService interface {
	ValidateTransactionRequest(req *dto.CreateTransactionRequest) error
}

// ValidateServiceImpl implements ValidateService
type ValidateServiceImpl struct{}

// NewValidateService creates a new ValidateService
func NewValidateService() ValidateService {
	return &ValidateServiceImpl{}
}

// ValidateTransactionRequest validates the transaction request
func (s *ValidateServiceImpl) ValidateTransactionRequest(req *dto.CreateTransactionRequest) error {
	validationErrors := &errors.ValidationErrors{
		Message: "Validation failed",
	}

	if req.CustomerNIK == "" {
		validationErrors.Add("customer_nik", "is required")
	}

	if req.OTR <= 0 {
		validationErrors.Add("otr", "must be greater than 0")
	}

	if req.AdminFee < 0 {
		validationErrors.Add("admin_fee", "cannot be negative")
	}

	if req.Installment <= 0 {
		validationErrors.Add("installment", "must be greater than 0")
	}

	if req.Interest < 0 {
		validationErrors.Add("interest", "cannot be negative")
	}

	if req.AssetName == "" {
		validationErrors.Add("asset_name", "is required")
	}

	if req.Tenor <= 0 {
		validationErrors.Add("tenor", "must be greater than 0")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

// HandleError handles errors and returns appropriate HTTP status code and response
func HandleError(err error) (int, interface{}) {
	switch e := err.(type) {
	case *errors.ValidationErrors:
		return e.Code, e
	case *errors.ErrorResponse:
		return e.Code, e
	default:
		// For unhandled errors, return 500 Internal Server Error
		return http.StatusInternalServerError, errors.ErrInternalServerError
	}
}
