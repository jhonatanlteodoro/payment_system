package workers

import (
	"context"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	"github.com/jhonatanlteodoro/payment_system/src/usecases"
	"log"
	"os"
	"sync"
	"syscall"
)

func RunAllWorkers(serverDown chan os.Signal) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		RunStartPaymentWorker(serverDown)
	}()
	go func() {
		defer wg.Done()
		RunProcessPaymentWorker(serverDown)
	}()

	go func() {
		defer wg.Done()
		RunNotifyWorker(serverDown)
	}()
	wg.Wait()
}

func RunStartPaymentWorker(serverDown chan os.Signal) {
	deps := shared_deps.GetSharedDependencies()

	ctx, cancel := context.WithCancel(context.TODO())
	log.Println("Running Start payment worker...")
	u := usecases.NewProcessPaymentUseCase(deps.ProcessPaymentQueue, deps.NotifyUserQueue, deps.PaymentDistributedLock)
	go func() {
		if err := u.Process(ctx); err != nil {
			log.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-serverDown
		log.Println("Quit Signal Received - Worker StartPayment is shutting down...")
		cancel()
		serverDown <- syscall.SIGINT // propagate
	}()

	wg.Wait()
	return
}

func RunProcessPaymentWorker(serverDown chan os.Signal) {
	deps := shared_deps.GetSharedDependencies()

	ctx, cancel := context.WithCancel(context.TODO())

	log.Println("Running Process Payment worker...")
	u := usecases.NewProcessPaymentUseCase(deps.ProcessPaymentQueue, deps.NotifyUserQueue, deps.PaymentDistributedLock)
	go func() {
		if err := u.Process(ctx); err != nil {
			log.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-serverDown
		log.Println("Quit Signal Received - Worker ProcessPayment is shutting down...")
		cancel()
		serverDown <- syscall.SIGINT // propagate
	}()

	wg.Wait()
	return
}

func RunNotifyWorker(serverDown chan os.Signal) {
	deps := shared_deps.GetSharedDependencies()

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Running Notify worker...")
	u := usecases.NewNotify(deps.NotifyUserQueue)
	go func() {
		if err := u.Notify(ctx); err != nil {
			log.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-serverDown
		log.Println("Quit Signal Received - Worker Notify is shutting down...")
		cancel()
		serverDown <- syscall.SIGINT // propagate
	}()

	wg.Wait()
}
