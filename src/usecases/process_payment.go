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

type ProcessPaymentUseCase struct {
	queue   ports.Queue
	cacheDB ports.DistributedLock

	notifyQueue ports.Queue
}

func NewProcessPaymentUseCase(queue, notifyQueue ports.Queue, cacheDB ports.DistributedLock) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		queue:       queue,
		cacheDB:     cacheDB,
		notifyQueue: notifyQueue,
	}
}

func (m *ProcessPaymentUseCase) workerProcessPayment(delivery amqp.Delivery) error {
	data := &types.Payment{}
	if err := json.Unmarshal(delivery.Body, data); err != nil {
		return err
	}

	fmt.Println("Processing something...")
	fmt.Println("Updating the database...")
	fmt.Println("generating next msg...")
	time.Sleep(2 * time.Second)

	log.Println("Releasing lock...")
	if err := m.cacheDB.ReleaseLock(context.TODO(), data.FromAccount); err != nil {
		log.Println("Failed to release lock ", err.Error())
		// the lock will expire so, no need to retry to release it here
	}
	log.Println("Lock released")

	log.Println("Sending to notifyQueue")
	if err := m.notifyQueue.Publish(context.TODO(), []byte("payment done")); err != nil {
		log.Println("Failed to send message ", err.Error())
		// failed to send to notification should not return an error since it wont change anything for us at this point
	}

	log.Println("Process Payment done...")
	return nil
}

func (m *ProcessPaymentUseCase) Process(ctx context.Context) error {
	errors := make(chan error)
	m.queue.WatchQueue(ctx, 5, errors, m.workerProcessPayment)

	for err := range errors {
		log.Println(err)
	}
	return nil
}
