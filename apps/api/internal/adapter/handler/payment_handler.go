package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/dto/request"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/response"
)

// PaymentHandler handles payment endpoints.
type PaymentHandler struct {
	paymentUC usecase.PaymentUsecase
}

// NewPaymentHandler creates a new payment handler.
func NewPaymentHandler(paymentUC usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUC: paymentUC}
}

// InitiatePayment godoc
// @Summary Initiate payment for an order
// @Tags Payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.InitiatePaymentRequest true "Payment payload"
// @Success 201 {object} response.Response
// @Router /payments [post]
func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req request.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID := middleware.GetUserID(c)

	input := &usecase.InitiatePaymentInput{
		OrderID: req.OrderID,
		UserID:  userID,
	}

	result, err := h.paymentUC.InitiatePayment(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "payment initiated", result)
}

// HandleCallback godoc
// @Summary Handle Midtrans payment callback
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /payments/callback [post]
func (h *PaymentHandler) HandleCallback(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := &usecase.PaymentCallbackInput{
		OrderID:           getStringFromMap(payload, "order_id"),
		TransactionID:     getStringFromMap(payload, "transaction_id"),
		TransactionStatus: getStringFromMap(payload, "transaction_status"),
		PaymentType:       getStringFromMap(payload, "payment_type"),
		FraudStatus:       getStringFromMap(payload, "fraud_status"),
	}

	if err := h.paymentUC.HandleCallback(c.Request.Context(), input); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "callback processed", nil)
}

// CheckStatus godoc
// @Summary Check payment status
// @Tags Payments
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} response.Response
// @Router /payments/{id}/status [get]
func (h *PaymentHandler) CheckStatus(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentUC.CheckStatus(c.Request.Context(), paymentID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "payment status retrieved", payment)
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
