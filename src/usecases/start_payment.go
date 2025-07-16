package usecases

import (
	"context"
	"encoding/json"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type StartPaymentUseCase struct {
	queue   ports.Queue
	cacheDB any
}

func NewStartPaymentUseCase(queue ports.Queue) *StartPaymentUseCase {
	return &StartPaymentUseCase{
		queue: queue,
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
