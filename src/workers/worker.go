package workers

import (
	"context"
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"time"
)

type DummyWorker struct {
	data []string
}

func (m *DummyWorker) HasMessageToProcess() bool {
	if len(m.data) > 0 {
		return true
	}

	// mock some random items
	num := rand.Intn(10)
	if num%2 == 0 {
		m.data = append(m.data, strconv.Itoa(num))
		fmt.Println("Added value to process: ", num)
	}
	return false
}

func (m *DummyWorker) Process(serverDown chan os.Signal) {
	for {
		select {
		case <-serverDown:
			log.Println("Quit Signal Received - Worker is shutting down...")
			return
		default:
			if m.HasMessageToProcess() {
				fmt.Println("dummy msg: ", m.data[len(m.data)-1])
				m.data = slices.Delete(m.data, len(m.data)-1, len(m.data))
				continue
			}

			fmt.Println("No message to process")
			time.Sleep(2 * time.Second)
		}
	}
}

func StartWorker(serverDown chan os.Signal) {
	deps := shared_deps.GetSharedDependencies()

	go func() {
		ctx := context.TODO()
		for i := 0; i < 10; i++ {
			if err := deps.StartPaymentQueue.Publish(ctx, []byte(fmt.Sprintf("num: %d", i))); err != nil {
				fmt.Println("Failed to publish message to queue ", err)
				continue
			}
			time.Sleep(1 * time.Second)
			fmt.Println("Published message to queue ", i)
		}
	}()

	ctx := context.TODO()
	errors := make(chan error)
	deps.StartPaymentQueue.WatchQueue(ctx, 3, errors, func(delivery amqp.Delivery) error {
		log.Println("Received a message: ", string(delivery.Body))
		return nil
	})

	for e := range errors {
		log.Println("Error processing message: ", e)
	}
	fmt.Println("Worker stopped")
}
