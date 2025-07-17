package usecases

import (
	"context"
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type NotifyUseCase struct {
	queue ports.Queue
}

func NewNotify(queue ports.Queue) *NotifyUseCase {
	return &NotifyUseCase{
		queue: queue,
	}
}

func (m *NotifyUseCase) workerNotify(delivery amqp.Delivery) error {
	msg := string(delivery.Body)

	fmt.Println("Msg received: ", msg)
	fmt.Println("Sending email...")
	time.Sleep(1 * time.Second)
	fmt.Println("done")

	return nil
}

func (m *NotifyUseCase) Notify(ctx context.Context) error {
	errors := make(chan error)
	m.queue.WatchQueue(ctx, 5, errors, m.workerNotify)

	for err := range errors {
		log.Println(err)
	}

	return nil
}
