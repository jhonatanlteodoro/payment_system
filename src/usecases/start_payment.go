package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	"github.com/jhonatanlteodoro/payment_system/src/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type StartPaymentUseCase struct {
	queue   ports.Queue
	cacheDB ports.DistributedLock

	processPaymentQueue ports.Queue

	paymentQuery ports.PaymentsQuery
}

func NewStartPaymentUseCase(queue, processPaymentQueue ports.Queue, cacheDB ports.DistributedLock, paymentQuerySvc ports.PaymentsQuery) *StartPaymentUseCase {
	return &StartPaymentUseCase{
		queue:               queue,
		processPaymentQueue: processPaymentQueue,
		cacheDB:             cacheDB,
		paymentQuery:        paymentQuerySvc,
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

	ctx := context.TODO()

	tx, err := shared_deps.GetSharedDependencies().PaymentsDBConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if !m.cacheDB.AcquireLock(ctx, data.FromAccount) {
		log.Println("Account is locked. Enqueuing msg")
		return m.StartPayment(context.Background(), data)
	}
	fmt.Println("Lock Acquired")

	exc := func() error {
		log.Println("Starting payment transaction")
		err := m.paymentQuery.CreatePayment(ctx, data, tx)
		if err != nil {
			return err
		}
		log.Println("payment created")

		payload, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err = m.processPaymentQueue.Publish(context.TODO(), payload); err != nil {
			return err
		}

		return tx.Commit(ctx)
	}()

	if exc != nil {
		m.cacheDB.ReleaseLock(ctx, data.FromAccount)
		return m.StartPayment(context.Background(), data)
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
