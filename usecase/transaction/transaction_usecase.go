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

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrLimitExceeded    = errors.New("transaction amount exceeds available limit")
)

type DBTx = repo.DBTx
type DB = repo.DB

type CustomerRepository interface {
	GetCustomer(ctx context.Context, nik string) (*model.Customer, error)
	CreateCustomer(ctx context.Context, customer *model.Customer) error
}

type LimitRepository interface {
	GetLimit(ctx context.Context, nik string, tenor int) (*model.CustomerLimit, error)
	UpdateLimit(ctx context.Context, tx repo.DBTx, nik string, tenor int, amount int64) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx repo.DBTx, transaction *model.Transaction) error
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error)
}

type transactionUsecase struct {
	db           *sqlx.DB
	customerRepo CustomerRepository
	limitRepo    LimitRepository
	txRepo       TransactionRepository
}

func NewTransactionUsecase(db *sqlx.DB, customerRepo CustomerRepository, limitRepo LimitRepository, txRepo TransactionRepository) TransactionUsecase {
	return &transactionUsecase{
		db:           db,
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		txRepo:       txRepo,
	}
}

func (u *transactionUsecase) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error) {
	dbTx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	customer, err := u.customerRepo.GetCustomer(ctx, req.CustomerNIK)
	if err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("gagal mendapatkan data customer: %w", err)
	}
	if customer == nil {
		dbTx.Rollback()
		return nil, ErrCustomerNotFound
	}

	limit, err := u.limitRepo.GetLimit(ctx, req.CustomerNIK, req.Tenor)
	if err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("gagal mendapatkan limit customer: %w", err)
	}

	totalAmount := req.OTR + req.AdminFee

	if limit.LimitAmount < totalAmount {
		dbTx.Rollback()
		return nil, ErrLimitExceeded
	}

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

	if err := u.txRepo.CreateTransaction(ctx, dbTx, transaction); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("gagal membuat transaksi: %w", err)
	}

	newLimitAmount := limit.LimitAmount - totalAmount
	if err := u.limitRepo.UpdateLimit(ctx, dbTx, req.CustomerNIK, req.Tenor, newLimitAmount); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("gagal memperbarui limit: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		dbTx.Rollback()
		return nil, fmt.Errorf("gagal melakukan commit transaksi: %w", err)
	}

	return transaction, nil
}
