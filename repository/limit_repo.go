package repository

import (
	"context"
	"database/sql"
	"log"

	"multifinance/model"

	"github.com/jmoiron/sqlx"
)

type LimitRepository struct {
	db *sqlx.DB
}

func NewLimitRepository(db *sqlx.DB) *LimitRepository {
	return &LimitRepository{db: db}
}

func (r *LimitRepository) GetLimit(ctx context.Context, nik string, tenor int) (*model.CustomerLimit, error) {
	var limit model.CustomerLimit
	query := "SELECT * FROM customer_limits WHERE customer_nik = ? AND tenor = ?"
	
	log.Printf("Debug - Executing query: %s with nik=%s, tenor=%d", query, nik, tenor)
	
	err := r.db.GetContext(ctx, &limit, query, nik, tenor)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Debug - No limit found for nik=%s, tenor=%d", nik, tenor)
			return nil, nil
		}
			log.Printf("Error - Failed to get limit: %v", err)
		return nil, err
	}
	
	log.Printf("Debug - Found limit: %+v", limit)
	return &limit, nil
}

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
