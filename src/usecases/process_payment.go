package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/shared_deps"
	"github.com/jhonatanlteodoro/payment_system/src/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type ProcessPaymentUseCase struct {
	queue   ports.Queue
	cacheDB ports.DistributedLock

	notifyQueue ports.Queue

	paymentQuery ports.PaymentsQuery
	balanceQuery ports.BalanceQuery
	rules        ports.ProcessingRules
	summaryQuery ports.QuarterlyAccountSummary
}

func NewProcessPaymentUseCase(
	queue, notifyQueue ports.Queue, cacheDB ports.DistributedLock,
	paymentQuery ports.PaymentsQuery, balanceQuery ports.BalanceQuery,
	summaryQuery ports.QuarterlyAccountSummary,
	paymentRules ports.ProcessingRules,
) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		queue:       queue,
		cacheDB:     cacheDB,
		notifyQueue: notifyQueue,

		paymentQuery: paymentQuery,
		balanceQuery: balanceQuery,
		summaryQuery: summaryQuery,
		rules:        paymentRules,
	}
}

func (m *ProcessPaymentUseCase) workerProcessPayment(delivery amqp.Delivery) error {
	data := &types.Payment{}
	if err := json.Unmarshal(delivery.Body, data); err != nil {
		return err
	}

	ctx := context.TODO()
	tx, err := shared_deps.GetSharedDependencies().PaymentsDBConn.Begin(ctx)
	if err != nil {
		return err
	}

	exc := func() error {
		balanceFrom, err := m.balanceQuery.GetBalance(ctx, data.FromAccount, tx)
		if err != nil {
			return err
		}

		if balanceFrom.Amount <= 0 {
			return errors.New("no available balanceQuery")
		}

		balanceTo, err := m.balanceQuery.GetBalance(ctx, data.ToAccount, tx)
		if err != nil {
			return err
		}

		if balanceFrom.Amount < data.Amount {
			log.Println("insufficient funds")
			data.Status = "insufficient funds"
			return m.paymentQuery.UpdatePaymentStatus(ctx, data, tx)
		}

		summary, err := m.summaryQuery.GetSummary(ctx, data.ToAccount, tx)
		if err != nil {
			return err
		}

		reason, err := m.rules.ProcessRules(ctx, data, summary)
		if err != nil {
			return err
		}

		if reason != "" {
			log.Println("payment not allowed: ", reason)
			data.Status = reason
			return m.paymentQuery.UpdatePaymentStatus(ctx, data, tx)
		}
		log.Println("payment allowed")

		balanceFrom.Amount -= data.Amount
		if err := m.balanceQuery.UpdateBalance(ctx, balanceFrom, tx); err != nil {
			return err
		}

		balanceTo.Amount += data.Amount
		if err := m.balanceQuery.UpdateBalance(ctx, balanceTo, tx); err != nil {
			return err
		}
		return nil
	}()

	if exc != nil {
		return err
	}

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
