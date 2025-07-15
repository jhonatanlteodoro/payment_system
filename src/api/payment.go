package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"net/http"
)

type payment struct{}

func newPayment() ports.Handler {
	return &payment{}
}

func (m *payment) RegisterRoute(router *gin.Engine) {
	router.POST(
		"/payment",
		m.StartTransaction,
	)
}

func (m *payment) StartTransaction(ctx *gin.Context) {
	fmt.Println("Hello there!")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "starting transaction",
	})
}
