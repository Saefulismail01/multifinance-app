package repository

import (
	"context"
	"database/sql"

	"multifinance/model"

	"github.com/jmoiron/sqlx"
)

// CustomerRepository implements the customer repository interface
type CustomerRepository struct {
	db *sqlx.DB
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetCustomer retrieves a customer by NIK.
func (r *CustomerRepository) GetCustomer(ctx context.Context, nik string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.GetContext(ctx, &customer, "SELECT * FROM customers WHERE nik = ?", nik)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &customer, err
}

// CreateCustomer creates a new customer.
func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO customers (nik, full_name, legal_name, birth_place, birth_date, salary, photo_ktp, photo_selfie)
		VALUES (:nik, :full_name, :legal_name, :birth_place, :birth_date, :salary, :photo_ktp, :photo_selfie)
	`, customer)
	return err
}
