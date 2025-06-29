package dto

import (
	"errors"
	"fmt"
)

// CreateTransactionRequest represents the request payload for creating a transaction.
type CreateTransactionRequest struct {
	CustomerNIK string `json:"customer_nik"`
	OTR         int64  `json:"otr"`
	AdminFee    int64  `json:"admin_fee"`
	Installment int64  `json:"installment"`
	Interest    int64  `json:"interest"`
	AssetName   string `json:"asset_name"`
	Tenor       int    `json:"tenor"`
}

// Validate validates the CreateTransactionRequest fields
func (r *CreateTransactionRequest) Validate() error {
	if r.CustomerNIK == "" {
		return errors.New("customer_nik is required")
	}
	if r.OTR <= 0 {
		return errors.New("otr must be greater than 0")
	}
	if r.AdminFee < 0 {
		return errors.New("admin_fee cannot be negative")
	}
	if r.Installment <= 0 {
		return errors.New("installment must be greater than 0")
	}
	if r.Interest < 0 {
		return errors.New("interest cannot be negative")
	}
	if r.AssetName == "" {
		return errors.New("asset_name is required")
	}
	if r.Tenor <= 0 {
		return errors.New("tenor must be greater than 0")
	}
	return nil
}

// CreateTransactionResponse represents the transaction response to the client.
type CreateTransactionResponse struct {
	ContractNumber string `json:"contract_number"`
	CustomerNIK    string `json:"customer_nik"`
	OTR            int64  `json:"otr"`
	AdminFee       int64  `json:"admin_fee"`
	Installment    int64  `json:"installment"`
	Interest       int64  `json:"interest"`
	AssetName      string `json:"asset_name"`
	Tenor          int    `json:"tenor"`
}

// ToString returns a string representation of the response
func (r *CreateTransactionResponse) ToString() string {
	return fmt.Sprintf(
		"ContractNumber: %s, CustomerNIK: %s, OTR: %d, Installment: %d, Tenor: %d",
		r.ContractNumber,
		r.CustomerNIK,
		r.OTR,
		r.Installment,
		r.Tenor,
	)
}
