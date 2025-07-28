package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	"github.com/jhonatanlteodoro/payment_system/src/usecases"
	"log"
	"net/http"
	"os"
	"time"
)

func registerRoutes(router *gin.Engine) {
	deps := shared_deps.GetSharedDependencies()
	startPayment := usecases.NewStartPaymentUseCase(deps.StartPaymentQueue, deps.ProcessPaymentQueue, deps.PaymentDistributedLock)

	newPayment(startPayment).RegisterRoute(router)
}

func StartApiServer(serverDown chan os.Signal) {
	router := gin.Default()
	registerRoutes(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-serverDown
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server error:", err)
	}

	log.Println("Server exiting")
}
