package service

import (
	"context"
	"fmt"
	"sync"

	"multifinance/model"
	"multifinance/repository"

	"github.com/jmoiron/sqlx"
)

// TransactionService defines the interface for transaction-related business logic.
type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction, tenor int) error
}

// TransactionServiceImpl implements the TransactionService interface.
type TransactionServiceImpl struct {
	customerRepo    repository.CustomerRepository
	limitRepo       repository.LimitRepository
	transactionRepo repository.TransactionRepository
	db              *sqlx.DB
	mutex           sync.Mutex
}

// NewTransactionService creates a new TransactionServiceImpl.
func NewTransactionService(
	customerRepo repository.CustomerRepository,
	limitRepo repository.LimitRepository,
	transactionRepo repository.TransactionRepository,
	db *sqlx.DB,
) TransactionService {
	return &TransactionServiceImpl{
		customerRepo:    customerRepo,
		limitRepo:       limitRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

// CreateTransaction handles the business logic for creating a new financial transaction.
// It ensures the operation is atomic by using a database transaction.
func (s *TransactionServiceImpl) CreateTransaction(ctx context.Context, transaction *model.Transaction, tenor int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Begin a new database transaction.
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case of an error. The rollback will be ignored if the transaction is committed.
	defer tx.Rollback() // nolint:errcheck

	// Check if the customer has sufficient limit for the transaction.
	limit, err := s.limitRepo.GetLimit(ctx, transaction.CustomerNIK, tenor)
	if err != nil {
		return fmt.Errorf("failed to get customer limit: %w", err)
	}
	if limit == nil || limit.LimitAmount < transaction.OTR+transaction.AdminFee {
		return fmt.Errorf("insufficient limit")
	}

	// Calculate the new limit and update it.
	newLimit := limit.LimitAmount - (transaction.OTR + transaction.AdminFee)
	if err := s.limitRepo.UpdateLimit(ctx, tx, transaction.CustomerNIK, tenor, newLimit); err != nil {
		return fmt.Errorf("failed to update customer limit: %w", err)
	}

	// Create the transaction record.
	if err := s.transactionRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit the transaction if all operations were successful.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
