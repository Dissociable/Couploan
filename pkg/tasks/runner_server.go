package tasks

import (
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"log"
	"time"
)

// limiter Rate is 10 events/sec and permits burst of at most 30 events.
var limiter *rate.Limiter

func StartTasksRunnerServer(
	sharedRedisPool *SharedRedisPool, c *services.Container, async bool, srv *asynq.Server,
) {
	// limit 1 event per minute
	// limiter = rate.NewLimiter(rate.Every(1*time.Minute), 2)

	// Determine the database based on the environment
	db := c.Config.Cache.Database
	if c.Config.App.Environment == config.EnvTest {
		db = c.Config.Cache.TestDatabase
	}
	// Build the worker server
	logLevel := asynq.InfoLevel
	// err := logLevel.Set(c.Config.App.LogLevel)
	// if err != nil {
	//	c.Logger.Error("failed to set asynq runner server's log level", zap.Error(err))
	// }
	asynqConfig := asynq.Config{
		// See asynq.Config for all available options and explanation
		Concurrency:    1,
		IsFailure:      func(err error) bool { return !IsRateLimitError(err) },
		RetryDelayFunc: retryDelay,
		// RetryDelayFunc: func(n int, e error, tgUser *asynq.Task) time.Duration {
		//	return 10 * time.Second
		// },
		Queues: map[string]int{
			"post": 3,
		},
		Logger:                   NewAsynqLogger(c.Logger.Named(c.Config.App.Name + "TasksRunner")),
		LogLevel:                 logLevel,
		ShutdownTimeout:          2 * time.Minute,
		DelayedTaskCheckInterval: 0,
		GroupAggregator:          nil,
	}
	if sharedRedisPool != nil {
		*srv = *asynq.NewServer(sharedRedisPool, asynqConfig)
	} else {
		*srv = *asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
				DB:       db,
				Username: c.Config.Cache.Username,
				Password: c.Config.Cache.Password,
			}, asynqConfig,
		)
	}

	// Map task types to the handlers
	mux := asynq.NewServeMux()

	// Start the worker server
	if async {
		if err := srv.Start(mux); err != nil {
			log.Fatalf("could not run worker server: %v", err)
		}
	} else {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run worker server: %v", err)
		}
	}
}

type AsynqLogger struct {
	BaseLogger *zap.Logger
}

func NewAsynqLogger(base *zap.Logger) *AsynqLogger {
	asynqLogger := &AsynqLogger{
		BaseLogger: base,
	}
	return asynqLogger
}

func (l *AsynqLogger) GetMessageAndFields(args ...interface{}) (msg string, fields []zap.Field) {
	var firstMsg string
	if len(args) > 0 {
		// check if args is a string, then assign it to firstMsg
		s, ok := args[0].(string)
		if ok {
			firstMsg = s
		}
	}
	for i, arg := range args[1:] {
		fields = append(fields, zap.Any(fmt.Sprintf("field_%d", i), arg))
	}
	return firstMsg, fields
}

func (l *AsynqLogger) Debug(args ...interface{}) {
	m, f := l.GetMessageAndFields(args...)
	l.BaseLogger.Debug(m, f...)
}

func (l *AsynqLogger) Info(args ...interface{}) {
	m, f := l.GetMessageAndFields(args...)
	l.BaseLogger.Info(m, f...)
}

func (l *AsynqLogger) Warn(args ...interface{}) {
	m, f := l.GetMessageAndFields(args...)
	l.BaseLogger.Warn(m, f...)
}

func (l *AsynqLogger) Error(args ...interface{}) {
	m, f := l.GetMessageAndFields(args...)
	l.BaseLogger.Warn(m, f...)
}

func (l *AsynqLogger) Fatal(args ...interface{}) {
	m, f := l.GetMessageAndFields(args...)
	l.BaseLogger.Warn(m, f...)
}
