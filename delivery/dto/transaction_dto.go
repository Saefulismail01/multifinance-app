package dto

import (
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

// Validate is a convenience method that delegates to the validation service
// This is kept for backward compatibility but will be deprecated in the future
func (r *CreateTransactionRequest) Validate() error {
	// This is now a placeholder that will be removed in future versions
	// Validation logic has been moved to the service layer
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
