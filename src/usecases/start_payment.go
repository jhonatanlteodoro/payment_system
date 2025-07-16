package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type StartPaymentUseCase struct {
	queue   ports.Queue
	cacheDB ports.DistributedLock

	processPaymentQueue ports.Queue
}

func NewStartPaymentUseCase(queue, processPaymentQueue ports.Queue) *StartPaymentUseCase {
	return &StartPaymentUseCase{
		queue:               queue,
		processPaymentQueue: processPaymentQueue,
	}
}

func (m *StartPaymentUseCase) StartPayment(ctx context.Context, data *types.Payment) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := m.queue.Publish(ctx, payload); err != nil {
		return err
	}

	log.Printf("Transaction published for processing. from %s to %s", data.FromAccount, data.ToAccount)
	return nil
}

func (m *StartPaymentUseCase) workerStartPayment(delivery amqp.Delivery) error {
	data := &types.Payment{}
	if err := json.Unmarshal(delivery.Body, data); err != nil {
		return err
	}

	if !m.cacheDB.AcquireLock(context.TODO(), data.FromAccount) {
		log.Println("Account is locked. Enqueuing msg")
		m.StartPayment(context.Background(), data)
		return nil
	}

	fmt.Println("Lock Acquired")
	fmt.Println("Processing something...")
	fmt.Println("Updating the database...")
	fmt.Println("generating next msg...")
	time.Sleep(1 * time.Second)

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := m.processPaymentQueue.Publish(context.TODO(), payload); err != nil {
		return err
	}

	log.Println("Message published for processing. into Payment Processing Queue")
	return nil
}

func (m *StartPaymentUseCase) ProcessStartPayment(ctx context.Context) error {
	errors := make(chan error)
	m.queue.WatchQueue(ctx, 5, errors, m.workerStartPayment)

	for err := range errors {
		log.Println(err)
	}
	return nil
}
