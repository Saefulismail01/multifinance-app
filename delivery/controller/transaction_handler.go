package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"multifinance/delivery/dto"
	"multifinance/errors"
	"multifinance/model"
	"multifinance/service"
	"multifinance/usecase/transaction"
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	transactionUsecase transaction.TransactionUsecase
	validateService    service.ValidateService
}

// NewTransactionHandler creates a new TransactionHandler.
func NewTransactionHandler(
	transactionUsecase transaction.TransactionUsecase,
	validateService service.ValidateService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
		validateService:    validateService,
	}
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

	// Bind the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Validate the request using the validation service
	if err := h.validateService.ValidateTransactionRequest(&req); err != nil {
		switch v := err.(type) {
		case *errors.ValidationErrors:
			c.JSON(v.Code, dto.ValidationErrorResponse(v))
		case *errors.ErrorResponse:
			c.JSON(v.Code, dto.ErrorResponse(v.Code, v.Message, v))
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, "Validation error", err))
		}
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
		switch e := err.(type) {
		case *errors.ErrorResponse:
			c.JSON(e.Code, e)
		case *errors.ValidationErrors:
			c.JSON(e.Code, e)
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse(http.StatusInternalServerError, "Failed to create transaction", err))
		}
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

	// Send the success response
	c.JSON(http.StatusCreated, dto.SuccessResponse(http.StatusCreated, "Transaction created successfully", response))
}
