package repository

import (
	"context"
	"database/sql"

	"multifinance/model"

	"github.com/jmoiron/sqlx"
)

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) GetCustomer(ctx context.Context, nik string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.GetContext(ctx, &customer, "SELECT * FROM customers WHERE nik = ?", nik)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO customers (nik, full_name, legal_name, birth_place, birth_date, salary, photo_ktp, photo_selfie)
		VALUES (:nik, :full_name, :legal_name, :birth_place, :birth_date, :salary, :photo_ktp, :photo_selfie)
	`, customer)
	return err
}
