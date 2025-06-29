package repository

import (
	"context"

	"multifinance/model"

	"github.com/jmoiron/sqlx"
)

// TransactionRepository implements the transaction repository interface
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository creates a new TransactionRepository.
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction within a database transaction.
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx DBTx, t *model.Transaction) error {
	query := `
		INSERT INTO transactions (contract_number, customer_nik, otr, admin_fee, installment, interest, asset_name, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	args := []interface{}{
		t.ContractNumber,
		t.CustomerNIK,
		t.OTR,
		t.AdminFee,
		t.Installment,
		t.Interest,
		t.AssetName,
		t.CreatedAt,
	}

	var err error
	if tx != nil {
		// Use the transaction if provided
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		// Otherwise use the database connection directly
		_, err = r.db.ExecContext(ctx, query, args...)
	}

	return err
}
