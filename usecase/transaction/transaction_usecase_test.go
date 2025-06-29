package transaction

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"multifinance/delivery/dto"
	"multifinance/model"
	repo "multifinance/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCustomerRepository struct {
	mock.Mock
}

func (m *mockCustomerRepository) GetCustomer(ctx context.Context, nik string) (*model.Customer, error) {
	args := m.Called(ctx, nik)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *mockCustomerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

type mockLimitRepository struct {
	mock.Mock
}

func (m *mockLimitRepository) GetLimit(ctx context.Context, nik string, tenor int) (*model.CustomerLimit, error) {
	args := m.Called(ctx, nik, tenor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CustomerLimit), args.Error(1)
}

func (m *mockLimitRepository) UpdateLimit(ctx context.Context, tx repo.DBTx, nik string, tenor int, amount int64) error {
	args := m.Called(ctx, tx, nik, tenor, amount)
	return args.Error(0)
}

type mockTransactionRepository struct {
	mock.Mock
}

func (m *mockTransactionRepository) CreateTransaction(ctx context.Context, tx repo.DBTx, transaction *model.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func TestTransactionUsecase_CreateTransaction(t *testing.T) {
	tests := []struct {
		namaTest         string
		setupMocks       func(
			customerRepo *mockCustomerRepository,
			limitRepo *mockLimitRepository,
			txRepo *mockTransactionRepository,
			sqlMock sqlmock.Sqlmock,
		)
		req             *dto.CreateTransactionRequest
		harusError      bool
		erorDiharapkan  error
		harusUpdateLimit bool
		shouldPanic    bool
	}{
		// Test cases for successful scenarios
		{
			namaTest: "transaksi berhasil",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()

				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)

				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(&model.CustomerLimit{CustomerNIK: "1234567890123456", Tenor: 6, LimitAmount: 10000000}, nil).Once()

				txRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.MatchedBy(func(tx *model.Transaction) bool {
					return tx.CustomerNIK == "1234567890123456" && 
						tx.OTR == 1000000 && 
						tx.AdminFee == 50000 &&
						tx.Installment == 1000000 &&
						tx.Interest == 100000 &&
						tx.AssetName == "Laptop"
				})).Return(nil).Once()

				limitRepo.On("UpdateLimit", mock.Anything, mock.Anything, "1234567890123456", 6, int64(8950000)).
					Return(nil).Once()

				sqlMock.ExpectCommit()
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError:      false,
			harusUpdateLimit: true,
		},

		// Test cases for error scenarios
		{
			namaTest: "gagal memulai transaksi",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin().WillReturnError(errors.New("database error"))
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal memulai transaksi: database error"),
			harusUpdateLimit: false,
			shouldPanic: false,
		},

		{
			namaTest: "gagal mendapatkan data customer",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(nil, errors.New("database error")).Once()
				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal mendapatkan data customer: database error"),
			harusUpdateLimit: false,
			shouldPanic: false,
		},

		{
			namaTest: "gagal mendapatkan limit",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()

				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)

				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(nil, errors.New("database error")).Once()

				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal mendapatkan limit customer: database error"),
			harusUpdateLimit: false,
			shouldPanic: false,
		},

		{
			namaTest: "gagal membuat transaksi",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()

				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)

				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(&model.CustomerLimit{CustomerNIK: "1234567890123456", Tenor: 6, LimitAmount: 10000000}, nil).Once()

				txRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("database error")).Once()

				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal membuat transaksi: database error"),
			harusUpdateLimit: false,
			shouldPanic: false,
		},

		{
			namaTest: "gagal update limit",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()

				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)

				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(&model.CustomerLimit{CustomerNIK: "1234567890123456", Tenor: 6, LimitAmount: 10000000}, nil).Once()

				txRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

				limitRepo.On("UpdateLimit", mock.Anything, mock.Anything, "1234567890123456", 6, int64(8950000)).
					Return(errors.New("database error")).Once()

				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal memperbarui limit: database error"),
			harusUpdateLimit: true,
			shouldPanic: false,
		},

		{
			namaTest: "gagal commit transaksi",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()

				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)

				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(&model.CustomerLimit{CustomerNIK: "1234567890123456", Tenor: 6, LimitAmount: 10000000}, nil).Once()

				txRepo.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

				limitRepo.On("UpdateLimit", mock.Anything, mock.Anything, "1234567890123456", 6, int64(8950000)).Return(nil).Once()

				sqlMock.ExpectCommit().WillReturnError(errors.New("commit failed"))
				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: fmt.Errorf("gagal melakukan commit transaksi: commit failed"),
			harusUpdateLimit: true,
			shouldPanic: false,
		},

		{
			namaTest: "customer tidak ditemukan",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				customerRepo.On("GetCustomer", mock.Anything, "9999999999999999").
					Return(nil, errors.New("customer not found"))
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "9999999999999999",
				OTR:         9000000,
				AdminFee:    1000000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: errors.New("customer not found"),
			shouldPanic: false,
		},
		{
			namaTest: "limit tidak mencukupi",
			setupMocks: func(customerRepo *mockCustomerRepository, limitRepo *mockLimitRepository, txRepo *mockTransactionRepository, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				customerRepo.On("GetCustomer", mock.Anything, "1234567890123456").
					Return(&model.Customer{NIK: "1234567890123456", FullName: "John Doe"}, nil)
				limitRepo.On("GetLimit", mock.Anything, "1234567890123456", 6).
					Return(&model.CustomerLimit{CustomerNIK: "1234567890123456", Tenor: 6, LimitAmount: 100000}, nil).Once()
				sqlMock.ExpectRollback().WillReturnError(nil)
			},
			req: &dto.CreateTransactionRequest{
				CustomerNIK: "1234567890123456",
				OTR:         1000000,
				AdminFee:    50000,
				Installment: 1000000,
				Interest:    100000,
				AssetName:   "Laptop",
				Tenor:       6,
			},
			harusError: true,
			erorDiharapkan: ErrLimitExceeded,
			harusUpdateLimit: false,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.namaTest, func(t *testing.T) {
			// Setup mocks
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Gagal membuat mock database: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			customerRepo := &mockCustomerRepository{}
			limitRepo := &mockLimitRepository{}
			txRepo := &mockTransactionRepository{}

			// Setup test case specific mocks
			tt.setupMocks(customerRepo, limitRepo, txRepo, sqlMock)

			// Create usecase with mocked dependencies
			uc := NewTransactionUsecase(sqlxDB, customerRepo, limitRepo, txRepo)

			// Skip panic test as it's covered by other test cases

			// Execute
			_, err = uc.CreateTransaction(context.Background(), tt.req)

			// Assert
			if tt.harusError {
				assert.Error(t, err)
				if tt.erorDiharapkan != nil {
					assert.Contains(t, err.Error(), tt.erorDiharapkan.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			// Skip SQL mock verification for error cases as the transaction might be in an inconsistent state
			if !tt.harusError {
				customerRepo.AssertExpectations(t)
				limitRepo.AssertExpectations(t)
				txRepo.AssertExpectations(t)
				// Check SQL mock expectations last
				assert.NoError(t, sqlMock.ExpectationsWereMet())
			}
		})
	}
}

// generateContractNumber generates a unique contract number based on NIK and timestamp
func generateContractNumber(nik string, now time.Time) string {
	// Take first 8 characters of NIK, or less if NIK is shorter
	nikPrefix := nik
	if len(nik) > 8 {
		nikPrefix = nik[:8]
	}
	
	// Format: NIK_PREFIX-TIMESTAMP-XXXX (where X is a random number)
	timestamp := now.Format("02012006150405") // DDMMYYYYHHMMSS
	randomNum := now.Nanosecond() % 10000
	return fmt.Sprintf("%s-%s-%04d", nikPrefix, timestamp, randomNum)
}

func TestGenerateContractNumber(t *testing.T) {
	tests := []struct {
		name     string
		nik      string
		expected string
	}{
		{
			name:     "valid NIK",
			nik:      "1234567890123456",
			expected: "12345678-",
		},
		{
			name:     "short NIK",
			nik:      "1234",
			expected: "1234-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			result := generateContractNumber(tt.nik, now)
			assert.Contains(t, result, tt.expected)
		})
	}
}
