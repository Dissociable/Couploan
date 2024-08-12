package tasks

import (
	"context"
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/rate"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

func PoolWatchdog(c *services.Container) {
	go func(c *services.Container) {
		checkInterval := 15 * time.Second
		for {
			time.Sleep(checkInterval)
			if err := sharedRedisPool.Client.Ping(context.Background()).Err(); err != nil {
				c.Logger.Error("redis client disconnected, trying to reconnect...", zap.Error(err))
				ReconnectRedisPool(c, true)
			}
		}
	}(c)
}

func ReconnectRedisPool(c *services.Container, async bool) {
	var err error
	db := c.Config.Cache.Database
	if c.Config.App.Environment == config.EnvTest {
		db = c.Config.Cache.TestDatabase
	}
	if asynqServer != nil {
		asynqServer.Shutdown()
	}
	// Close
	if sharedRedisPool != nil {
		_ = sharedRedisPool.Client.Close()
	}
	// Reconnect
	sharedRedisPool, err = NewSharedRedisPool(
		SharedRedisPoolOpts{
			DSN:      fmt.Sprintf("redis://%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
			DB:       db,
			Username: c.Config.Cache.Username,
			Password: c.Config.Cache.Password,
		},
	)
	if err != nil {
		c.Logger.Fatal("failed to create shared redis pool", zap.Error(err))
	}
	time.Sleep(5 * time.Second)
	InitializeSemaphore("global", 1)
	asynqServer = new(asynq.Server)
	StartTasksRunnerServer(sharedRedisPool, c, async, asynqServer)
	asynqInspector = new(asynq.Inspector)
	StartInspector(sharedRedisPool, c, asynqInspector)
	time.Sleep(5 * time.Second)
}

// InitializeSemaphore initializes the semaphore
func InitializeSemaphore(scope string, maxTokens int) {
	Semaphore = rate.NewSemaphore(sharedRedisPool, scope, maxTokens)
}

type RateLimitError struct {
	RetryIn time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited (retry in  %v)", e.RetryIn)
}

func IsRateLimitError(err error) bool {
	var rateLimitError *RateLimitError
	ok := errors.As(err, &rateLimitError)
	return ok
}

func retryDelay(n int, err error, task *asynq.Task) time.Duration {
	var rateLimitErr *RateLimitError
	if errors.As(err, &rateLimitErr) {
		return rateLimitErr.RetryIn
	}
	return asynq.DefaultRetryDelayFunc(n, err, task)
}
