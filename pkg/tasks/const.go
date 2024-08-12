package tasks

import (
	"github.com/hibiken/asynq"

	"github.com/hibiken/asynq/x/rate"
)

var (
	sharedRedisPool *SharedRedisPool
	Semaphore       *rate.Semaphore
	asynqServer     *asynq.Server
	asynqInspector  *asynq.Inspector
)
