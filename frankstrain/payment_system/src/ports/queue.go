package ports

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	// Publish will be used to handle publishing message for the configured queue in use
	Publish(ctx context.Context, data []byte) error

	// WatchQueue will be used to handle processing messages for the configured queue in use
	WatchQueue(_ context.Context, maxWorkers int, errors chan error, worker func(delivery amqp.Delivery) error)
}
