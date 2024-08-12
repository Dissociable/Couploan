package tasks

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type (
	SharedRedisPoolOpts struct {
		DSN      string
		DB       int
		Username string
		Password string
	}

	SharedRedisPool struct {
		Client *redis.Client
	}
)

func NewSharedRedisPool(o SharedRedisPoolOpts) (*SharedRedisPool, error) {
	redisOpts, err := redis.ParseURL(o.DSN)
	if err != nil {
		return nil, err
	}

	redisOpts.DB = o.DB
	redisOpts.ClientName = "asynq_shared_redis"
	redisOpts.MaxActiveConns = 0
	redisOpts.MaxRetries = 7
	redisOpts.MinRetryBackoff = 1 * time.Second
	redisOpts.Username = o.Username
	redisOpts.Password = o.Password
	redisClient := redis.NewClient(redisOpts)

	return &SharedRedisPool{
		Client: redisClient,
	}, nil
}

// MakeRedisClient Fulfills asynq's Interface
func (r *SharedRedisPool) MakeRedisClient() interface{} {
	return r.Client
}

// Somewhere else when creating an asynq server asynq.NewServer()
// use *pool.SharedRedisPool
