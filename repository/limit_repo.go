package repository

import (
	"context"
	"database/sql"

	"multifinance/model"
	"github.com/jmoiron/sqlx"
)

// LimitRepository implements the limit repository interface
type LimitRepository struct {
	db *sqlx.DB
}

// NewLimitRepository creates a new LimitRepository.
func NewLimitRepository(db *sqlx.DB) *LimitRepository {
	return &LimitRepository{db: db}
}

// GetLimit retrieves a customer's limit by NIK and tenor.
func (r *LimitRepository) GetLimit(ctx context.Context, nik string, tenor int) (*model.CustomerLimit, error) {
	var limit model.CustomerLimit
	err := r.db.GetContext(ctx, &limit, "SELECT * FROM customer_limits WHERE customer_nik = ? AND tenor = ?", nik, tenor)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &limit, err
}

// UpdateLimit updates a customer's limit within a transaction.
func (r *LimitRepository) UpdateLimit(ctx context.Context, tx DBTx, nik string, tenor int, amount int64) error {
	query := `
		UPDATE customer_limits 
		SET limit_amount = ? 
		WHERE customer_nik = ? AND tenor = ?`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, amount, nik, tenor)
	} else {
		_, err = r.db.ExecContext(ctx, query, amount, nik, tenor)
	}

	return err
}
