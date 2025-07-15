package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func registerRoutes(router *gin.Engine) {
	newPayment().RegisterRoute(router)
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
