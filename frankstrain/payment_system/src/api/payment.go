package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
	"github.com/jhonatanlteodoro/payment_system/src/usecases"
	"log"
	"net/http"
)

type payment struct {
	startPaymentUseCase *usecases.StartPaymentUseCase
}

func newPayment(startPaymentUseCase *usecases.StartPaymentUseCase) ports.Handler {
	return &payment{
		startPaymentUseCase: startPaymentUseCase,
	}
}

func (m *payment) RegisterRoute(router *gin.Engine) {
	router.POST(
		"/payment",
		m.StartTransaction,
	)
}

func (m *payment) StartTransaction(ctx *gin.Context) {
	payment := &types.Payment{}
	if err := ctx.ShouldBindJSON(payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := m.startPaymentUseCase.StartPayment(ctx, payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "payment started"})
	log.Println("payment started. sent to queue to be processed")
}
