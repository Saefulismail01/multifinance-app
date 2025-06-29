package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"multifinance/delivery/dto"
	"multifinance/model"
	"multifinance/usecase/transaction"
)

// Response is a standard API response structure
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	transactionUsecase transaction.TransactionUsecase
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(transactionUsecase transaction.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase: transactionUsecase}
}

// RegisterRoutes registers all transaction routes
func (h *TransactionHandler) RegisterRoutes(router *gin.RouterGroup) {
	transactionGroup := router.Group("/transactions")
	{
		transactionGroup.POST("", h.CreateTransaction)
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest

	// Bind and validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Additional validation
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Create the transaction domain object
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

	// Call the usecase to create the transaction
	if err := h.transactionUsecase.CreateTransaction(c.Request.Context(), transaction, req.Tenor); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create transaction",
			Error:   err.Error(),
		})
		return
	}

	// Create the response
	response := dto.CreateTransactionResponse{
		ContractNumber: transaction.ContractNumber,
		CustomerNIK:    transaction.CustomerNIK,
		OTR:            transaction.OTR,
		AdminFee:       transaction.AdminFee,
		Installment:    transaction.Installment,
		Interest:       transaction.Interest,
		AssetName:      transaction.AssetName,
		Tenor:          req.Tenor,
	}

	// Send the response
	c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Transaction created successfully",
		Data:    response,
	})
}
