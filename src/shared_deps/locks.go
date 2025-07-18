package shared_deps

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

// PaymentDistributedLock will handle the acquire and release step for our distributed lock
// initially it will simplify as much as possible, so I can continue progressing with the development
// when the core is ready I'll be back to refactor it and change the basic redis server to a redis cluster so we can
// play it properly. :)
type PaymentDistributedLock struct {
	dbConn *redis.Conn
}

func NewPaymentDistributedLock(dbClient *redis.Conn) *PaymentDistributedLock {
	return &PaymentDistributedLock{
		dbConn: dbClient,
	}
}

func (m *PaymentDistributedLock) AcquireLock(ctx context.Context, key string) bool {
	// Yes, that's not how a distributed lock should
	// be safe, do not copy it :)
	cmd, err := m.dbConn.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if lock exists: %s", err.Error())
		return false
	}

	if cmd > 0 {
		// lock already acquired
		return false
	}

	cmdResult := m.dbConn.Set(ctx, key, "acquired", time.Minute*5)
	if cmdResult.Err() != nil {
		log.Printf("Error acquiring lock: %s", cmdResult.Err().Error())
		return false
	}

	log.Printf("Acquired lock for: %s\n", key)
	return true
}

func (m *PaymentDistributedLock) ReleaseLock(ctx context.Context, key string) error {
	cmdResult := m.dbConn.Del(ctx, key)
	if cmdResult.Err() != nil {
		return fmt.Errorf("error releasing lock: %v", cmdResult.Err().Error())
	}

	log.Printf("Released lock for: %s\n", key)
	return nil
}
