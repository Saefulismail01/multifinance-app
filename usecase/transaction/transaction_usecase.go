package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"multifinance/delivery/dto"
	"multifinance/model"
	repo "multifinance/repository"

	"github.com/jmoiron/sqlx"
)

// Error variables
var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrLimitExceeded    = errors.New("transaction amount exceeds available limit")
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
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error)
}

type transactionUsecase struct {
	db           *sqlx.DB
	customerRepo CustomerRepository
	limitRepo    LimitRepository
	txRepo       TransactionRepository
}

// NewTransactionUsecase creates a new transaction usecase
func NewTransactionUsecase(db *sqlx.DB, customerRepo CustomerRepository, limitRepo LimitRepository, txRepo TransactionRepository) TransactionUsecase {
	return &transactionUsecase{
		db:           db,
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		txRepo:       txRepo,
	}
}

// CreateTransaction creates a new transaction
func (u *transactionUsecase) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error) {
	// Start a database transaction
	dbTx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	// Check if customer exists
	customer, err := u.customerRepo.GetCustomer(ctx, req.CustomerNIK)
	if err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		dbTx.Rollback()
		return nil, ErrCustomerNotFound
	}

	// Get customer limit for the requested tenor
	limit, err := u.limitRepo.GetLimit(ctx, req.CustomerNIK, req.Tenor)
	if err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("failed to get customer limit: %w", err)
	}

	// Calculate total amount
	totalAmount := req.OTR + req.AdminFee

	// Check if limit is sufficient
	if limit.LimitAmount < totalAmount {
		dbTx.Rollback()
		return nil, ErrLimitExceeded
	}

	// Create transaction record
	transaction := &model.Transaction{
		ContractNumber: fmt.Sprintf("CON-%d", time.Now().UnixNano()),
		CustomerNIK:    req.CustomerNIK,
		OTR:            req.OTR,
		AdminFee:       req.AdminFee,
		Installment:    req.Installment,
		Interest:       req.Interest,
		AssetName:      req.AssetName,
		CreatedAt:      time.Now(),
	}

	// Save transaction
	if err := u.txRepo.CreateTransaction(ctx, dbTx, transaction); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update customer limit
	newLimitAmount := limit.LimitAmount - totalAmount
	if err := u.limitRepo.UpdateLimit(ctx, dbTx, req.CustomerNIK, req.Tenor, newLimitAmount); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("failed to update customer limit: %w", err)
	}

	// Commit transaction
	if err := dbTx.Commit(); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}
