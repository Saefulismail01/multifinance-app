package transaction

import (
	"context"
	"errors"
	"sync"

	"multifinance/model"
	repo "multifinance/repository"
)

// DBTx is an alias for repository.DBTx
type DBTx = repo.DBTx

// DB is an alias for repository.DB
type DB = repo.DB

// CustomerRepository defines the interface for customer data operations.
type CustomerRepository interface {
	GetCustomer(ctx context.Context, nik string) (*model.Customer, error)
	CreateCustomer(ctx context.Context, customer *model.Customer) error
}

// LimitRepository defines the interface for limit data operations.
type LimitRepository interface {
	GetLimit(ctx context.Context, nik string, tenor int) (*model.CustomerLimit, error)
	UpdateLimit(ctx context.Context, tx repo.DBTx, nik string, tenor int, amount int64) error
}

// TransactionRepository defines the interface for transaction data operations.
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx repo.DBTx, transaction *model.Transaction) error
}

// TransactionUsecase defines the interface for transaction use cases
type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction, tenor int) error
}

type transactionUsecase struct {
	transactionRepo TransactionRepository
	customerRepo    CustomerRepository
	limitRepo       LimitRepository
	db              DB
	mu              sync.Mutex
}

// NewTransactionUsecase creates a new transaction usecase
func NewTransactionUsecase(
	transactionRepo TransactionRepository,
	customerRepo CustomerRepository,
	limitRepo LimitRepository,
	db DB,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
		limitRepo:       limitRepo,
		db:              db,
	}
}

// CreateTransaction handles the business logic for creating a transaction
func (u *transactionUsecase) CreateTransaction(ctx context.Context, transaction *model.Transaction, tenor int) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Check if customer exists
	customer, err := u.customerRepo.GetCustomer(ctx, transaction.CustomerNIK)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("customer not found")
	}

	// Get customer's limit for the requested tenor
	limit, err := u.limitRepo.GetLimit(ctx, transaction.CustomerNIK, tenor)
	if err != nil {
		return err
	}
	if limit == nil {
		return errors.New("limit not found for the requested tenor")
	}

	// Calculate total amount including interest and admin fee
	totalAmount := transaction.OTR + transaction.AdminFee
	if totalAmount > limit.LimitAmount {
		return errors.New("transaction amount exceeds available limit")
	}

	// Start a database transaction
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the customer's limit
	newLimitAmount := limit.LimitAmount - totalAmount
	if err := u.limitRepo.UpdateLimit(ctx, tx, transaction.CustomerNIK, tenor, newLimitAmount); err != nil {
		return err
	}

	// Create the transaction
	if err := u.transactionRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
