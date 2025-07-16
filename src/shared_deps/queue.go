package shared_deps

import (
	"context"
	"fmt"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type Queue struct {
	SeverConn *amqp.Connection
	QueueName string

	publisherChLock  sync.RWMutex
	PublisherChannel *amqp.Channel

	watcherChLock  sync.RWMutex
	watcherChannel *amqp.Channel
}

func NewQueue(serverConnection *amqp.Connection, queueName string) ports.Queue {
	return &Queue{
		SeverConn: serverConnection,
		QueueName: queueName,
	}
}

// channelIsValid return false if ch is nil or is closed otherwise return true
func (m *Queue) channelIsValid(ch *amqp.Channel) bool {
	return ch != nil && ch.IsClosed() == false
}

func (m *Queue) createNewPublisherChannel() error {
	m.publisherChLock.Lock()
	defer m.publisherChLock.Unlock()

	// double-check grants that if another thread was trying to reconnect
	// we do not override channels
	if !m.channelIsValid(m.PublisherChannel) {
		ch, err := m.SeverConn.Channel()
		if err != nil {
			return err
		}

		m.PublisherChannel = ch
	}
	return nil
}

func (m *Queue) getPublisherChannel() (*amqp.Channel, error) {
	m.publisherChLock.RLock()
	isValid := m.channelIsValid(m.PublisherChannel)
	m.publisherChLock.RUnlock()

	if !isValid {
		if err := m.createNewPublisherChannel(); err != nil {
			return nil, fmt.Errorf("error creating new publisher channel: %v", err)
		}
	}

	m.publisherChLock.RLock()
	defer m.publisherChLock.RUnlock()
	return m.PublisherChannel, nil
}

func (m *Queue) Publish(ctx context.Context, data []byte) error {

	ch, err := m.getPublisherChannel()
	if err != nil {
		return err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// even dough this will not guarantee the message delivery to the server
	// I'll keep it this way for now, until i complete the other things
	return ch.PublishWithContext(
		ctxWithTimeout, "",
		m.QueueName,
		false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		},
	)
}

func (m *Queue) createNewWatcherChannel() error {
	m.watcherChLock.Lock()
	defer m.watcherChLock.Unlock()

	// double-check grants that if another thread was trying to reconnect
	// we do not override channels
	if !m.channelIsValid(m.watcherChannel) {
		ch, err := m.SeverConn.Channel()
		if err != nil {
			return err
		}
		m.watcherChannel = ch
	}
	return nil
}

func (m *Queue) getWatcherChannel() (*amqp.Channel, error) {
	m.watcherChLock.RLock()
	isValid := m.channelIsValid(m.watcherChannel)
	m.watcherChLock.RUnlock()

	if !isValid {
		if err := m.createNewWatcherChannel(); err != nil {
			return nil, fmt.Errorf("error creating new publisher channel: %v", err)
		}
	}

	m.watcherChLock.RLock()
	defer m.watcherChLock.RUnlock()
	return m.watcherChannel, nil
}

func (m *Queue) WatchQueue(_ context.Context, maxWorkers int, errors chan error, worker func(delivery amqp.Delivery) error) {

	ch, err := m.getWatcherChannel()
	if err != nil {
		errors <- fmt.Errorf("error creating new publisher channel: %v", err)
		return
	}

	queuedChanMessages, err := ch.Consume(
		m.QueueName,
		"",    // consumer
		false, // auto-ack - as false to ensure we process it properly before the server mark it as done
		false, // exclusive - as false to guarantee that rabbitmq will deliver for more than one consumer
		false, // no-local - following what the doc says, this flag is useless
		false, // no-wait - we should keep it false to guarantee the server will always confirm that we have received the msg
		nil,   // args
	)
	if err != nil {
		errors <- fmt.Errorf("error creating new publisher channel: %v", err)
		return
	}

	workersChan := make(chan struct{}, maxWorkers)

	for msg := range queuedChanMessages {
		workersChan <- struct{}{}
		go func(c chan struct{}) {
			defer func() { <-c }()
			errWorker := worker(msg)
			if errWorker != nil {
				errors <- fmt.Errorf("error processing message: %v", errWorker)
				if errNack := msg.Nack(false, true); errNack != nil {
					errors <- fmt.Errorf("error nacking message: %v", errNack)
				}
				return
			}

			if errAck := msg.Ack(false); errAck != nil {
				errors <- fmt.Errorf("error acking message: %v", errAck)
			}
		}(workersChan)
	}
}
