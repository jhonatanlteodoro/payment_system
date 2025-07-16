package ports

import "context"

type DistributedLock interface {
	AcquireLock(ctx context.Context, key string) bool
	ReleaseLock(ctx context.Context, key string) error
}
