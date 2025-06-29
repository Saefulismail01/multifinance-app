package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"multifinance/delivery/dto"
	"multifinance/service"
	"multifinance/usecase/transaction"
)

type TransactionHandler struct {
	transactionUsecase transaction.TransactionUsecase
	validateService    service.ValidateService
}

func NewTransactionHandler(
	transactionUsecase transaction.TransactionUsecase,
	validateService service.ValidateService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
		validateService:    validateService,
	}
}

func (h *TransactionHandler) RegisterRoutes(router *gin.RouterGroup) {
	transactionGroup := router.Group("/transactions")
	{
		transactionGroup.POST("", h.CreateTransaction)
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "Invalid request"))
		return
	}

	if err := h.validateService.ValidateTransactionRequest(&req); err != nil {
		if vErr, ok := err.(interface{ GetErrors() []dto.ValidationError }); ok {
			c.JSON(http.StatusUnprocessableEntity, dto.Response{
				Code:    http.StatusUnprocessableEntity,
				Status:  http.StatusText(http.StatusUnprocessableEntity),
				Message: "Validation failed",
				Errors:  vErr.GetErrors(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "Internal server error"))
		return
	}

	tx, err := h.transactionUsecase.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case transaction.ErrCustomerNotFound:
			c.JSON(http.StatusNotFound, dto.NewErrorResponse(http.StatusNotFound, "Customer not found"))
		case transaction.ErrLimitExceeded:
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "Transaction amount exceeds available limit"))
		default:
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(http.StatusInternalServerError, "Failed to create transaction"))
		}
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(tx))
}
