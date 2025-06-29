package model

import "time"

// Customer represents the customer entity.
type Customer struct {
	NIK         string `db:"nik"`
	FullName    string `db:"full_name"`
	LegalName   string `db:"legal_name"`
	BirthPlace  string `db:"birth_place"`
	BirthDate   string `db:"birth_date"`
	Salary      int64  `db:"salary"`
	PhotoKTP    string `db:"photo_ktp"`
	PhotoSelfie string `db:"photo_selfie"`
}

// CustomerLimit represents the customer's credit limit for a specific tenor.
type CustomerLimit struct {
	CustomerNIK string `db:"customer_nik"`
	Tenor       int    `db:"tenor"`
	LimitAmount int64  `db:"limit_amount"`
}

// Transaction represents a financial transaction.
type Transaction struct {
	ContractNumber string    `db:"contract_number"`
	CustomerNIK    string    `db:"customer_nik"`
	OTR            int64     `db:"otr"`
	AdminFee       int64     `db:"admin_fee"`
	Installment    int64     `db:"installment"`
	Interest       int64     `db:"interest"`
	AssetName      string    `db:"asset_name"`
	CreatedAt      time.Time `db:"created_at"`
}

// CreateTransactionDTO represents the data transfer object for creating a transaction.
type CreateTransactionDTO struct {
	CustomerNIK string `json:"customer_nik"`
	OTR         int64  `json:"otr"`
	AdminFee    int64  `json:"admin_fee"`
	Installment int64  `json:"installment"`
	Interest    int64  `json:"interest"`
	AssetName   string `json:"asset_name"`
	Tenor       int    `json:"tenor"`
}
