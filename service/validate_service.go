package service

import (
	"net/http"
	"multifinance/delivery/dto"
)

type ValidateService interface {
	ValidateTransactionRequest(req *dto.CreateTransactionRequest) error
}

type ValidateServiceImpl struct{}

func NewValidateService() ValidateService {
	return &ValidateServiceImpl{}
}

func (s *ValidateServiceImpl) ValidateTransactionRequest(req *dto.CreateTransactionRequest) error {
	var validationErrs []dto.ValidationError

	if req.CustomerNIK == "" {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "customer_nik",
			Message: "customer_nik is required",
		})
	}

	if req.OTR <= 0 {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "otr",
			Message: "otr must be greater than 0",
		})
	}

	if req.AdminFee < 0 {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "admin_fee",
			Message: "admin_fee cannot be negative",
		})
	}

	if req.Installment <= 0 {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "installment",
			Message: "installment must be greater than 0",
		})
	}

	if req.Interest < 0 {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "interest",
			Message: "interest cannot be negative",
		})
	}

	if req.AssetName == "" {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "asset_name",
			Message: "asset_name is required",
		})
	}

	if req.Tenor <= 0 {
		validationErrs = append(validationErrs, dto.ValidationError{
			Field:   "tenor",
			Message: "tenor must be greater than 0",
		})
	}

	if len(validationErrs) > 0 {
		return dto.NewValidationError(validationErrs)
	}

	return nil
}

func HandleError(err error) (int, interface{}) {
	return http.StatusInternalServerError, map[string]interface{}{
		"error":   "Internal Server Error",
		"message": err.Error(),
	}
}
