package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"multifinance/delivery/dto"
	"multifinance/model"
	transactionUsecase "multifinance/usecase/transaction"
)

// MockTransactionUsecase is a mock implementation of TransactionUsecase
type MockTransactionUsecase struct {
	mock.Mock
}

func (m *MockTransactionUsecase) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*model.Transaction, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

// MockValidateService is a mock implementation of ValidateService
type MockValidateService struct {
	mock.Mock
}

func (m *MockValidateService) ValidateTransactionRequest(req *dto.CreateTransactionRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func setupRouter(handler *TransactionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api")
	handler.RegisterRoutes(api)
	return r
}

func TestTransactionHandler_CreateTransaction_Success(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Mock data
	req := dto.CreateTransactionRequest{
		CustomerNIK: "1234567890123456",
		OTR:         1000000,
		AdminFee:    50000,
		Installment: 1000000,
		Interest:    100000,
		AssetName:   "Laptop",
		Tenor:       6,
	}

	expectedTx := &model.Transaction{
		ContractNumber: "CON-1234567890",
		CustomerNIK:    req.CustomerNIK,
		OTR:           req.OTR,
		AdminFee:      req.AdminFee,
		Installment:    req.Installment,
		Interest:      req.Interest,
		AssetName:     req.AssetName,
		CreatedAt:      time.Now(),
	}

	// Expectations
	mockValidate.On("ValidateTransactionRequest", &req).Return(nil)
	mockUsecase.On("CreateTransaction", mock.Anything, &req).Return(expectedTx, nil)

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	jsonReq, _ := json.Marshal(req)
	reqBody := bytes.NewBuffer(jsonReq)
	httpReq, _ := http.NewRequest("POST", "/api/transactions", reqBody)

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response dto.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusOK), response.Status)

	// Verify mocks
	mockValidate.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransaction_InvalidJSON(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBufferString("{invalid json}"))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request", response["message"])
}

func TestTransactionHandler_CreateTransaction_ValidationError(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Mock data
	req := dto.CreateTransactionRequest{
		// Missing required fields to trigger validation error
	}

	// Expectations
	validationErrs := []dto.ValidationError{
		{Field: "customer_nik", Message: "customer_nik is required"},
		{Field: "otr", Message: "otr must be greater than 0"},
		{Field: "asset_name", Message: "asset_name is required"},
		{Field: "tenor", Message: "tenor must be greater than 0"},
	}
	mockValidate.On("ValidateTransactionRequest", &req).Return(dto.NewValidationError(validationErrs))

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	jsonReq, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(jsonReq))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response dto.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusUnprocessableEntity), response.Status)
	assert.Equal(t, "Validation failed", response.Message)
	assert.NotEmpty(t, response.Errors)

	// Verify mocks
	mockValidate.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransaction_CustomerNotFound(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Mock data
	req := dto.CreateTransactionRequest{
		CustomerNIK: "9999999999999999",
		OTR:         1000000,
		AdminFee:    50000,
		Installment: 1000000,
		Interest:    100000,
		AssetName:   "Laptop",
		Tenor:       6,
	}

	// Expectations
	mockValidate.On("ValidateTransactionRequest", &req).Return(nil)
	mockUsecase.On("CreateTransaction", mock.Anything, &req).
		Return(nil, transactionUsecase.ErrCustomerNotFound)

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	jsonReq, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(jsonReq))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Customer not found", response["message"])

	// Verify mocks
	mockValidate.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransaction_LimitExceeded(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Mock data
	req := dto.CreateTransactionRequest{
		CustomerNIK: "1234567890123456",
		OTR:         1000000000, // Very large amount to trigger limit exceeded
		AdminFee:    50000,
		Installment: 1000000,
		Interest:    100000,
		AssetName:   "Laptop",
		Tenor:       6,
	}

	// Expectations
	mockValidate.On("ValidateTransactionRequest", &req).Return(nil)
	mockUsecase.On("CreateTransaction", mock.Anything, &req).
		Return(nil, transactionUsecase.ErrLimitExceeded)

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	jsonReq, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(jsonReq))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transaction amount exceeds available limit", response["message"])

	// Verify mocks
	mockValidate.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransaction_InternalServerError(t *testing.T) {
	// Setup
	mockUsecase := new(MockTransactionUsecase)
	mockValidate := new(MockValidateService)
	handler := NewTransactionHandler(mockUsecase, mockValidate)

	// Mock data
	req := dto.CreateTransactionRequest{
		CustomerNIK: "1234567890123456",
		OTR:         1000000,
		AdminFee:    50000,
		Installment: 1000000,
		Interest:    100000,
		AssetName:   "Laptop",
		Tenor:       6,
	}

	// Expectations
	errMsg := "database connection failed"
	mockValidate.On("ValidateTransactionRequest", &req).Return(nil)
	mockUsecase.On("CreateTransaction", mock.Anything, &req).
		Return(nil, errors.New(errMsg))

	// Execute
	r := setupRouter(handler)
	w := httptest.NewRecorder()
	jsonReq, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(jsonReq))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to create transaction", response["message"])

	// Verify mocks
	mockValidate.AssertExpectations(t)
	mockUsecase.AssertExpectations(t)
}
