package logger

import (
	"context"
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type ctxKey struct{}

var once sync.Once

var logger *zap.Logger

var LogLevel string

var LokiConfig config.LokiConfig

// Get initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
func Get() *zap.Logger {
	once.Do(
		func() {
			stdout := zapcore.AddSync(os.Stdout)

			file := zapcore.AddSync(
				&lumberjack.Logger{
					Filename:   "logs/couploan.log",
					MaxSize:    5,
					MaxBackups: 3,
					MaxAge:     7,
					Compress:   true,
				},
			)

			level := zap.InfoLevel
			levelEnv := LogLevel
			if levelEnv != "" {
				levelFromEnv, err := zapcore.ParseLevel(levelEnv)
				if err != nil {
					log.Println(
						fmt.Errorf("invalid level, defaulting to INFO: %w", err),
					)
				}

				level = levelFromEnv
			}

			logLevel := zap.NewAtomicLevelAt(level)

			productionCfg := zap.NewProductionEncoderConfig()
			productionCfg.TimeKey = "timestamp"
			productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

			developmentCfg := zap.NewDevelopmentEncoderConfig()
			developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

			consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
			fileEncoder := zapcore.NewJSONEncoder(productionCfg)

			var gitRevision string

			buildInfo, ok := debug.ReadBuildInfo()
			if ok {
				for _, v := range buildInfo.Settings {
					if v.Key == "vcs.revision" {
						gitRevision = v.Value
						break
					}
				}
			}

			// log to multiple destinations (console and file)
			// extra fields are added to the JSON output alone
			core := zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, stdout, logLevel),
				zapcore.NewCore(fileEncoder, file, logLevel).
					With(
						[]zapcore.Field{
							zap.String("git_revision", gitRevision),
							zap.String("go_version", buildInfo.GoVersion),
						},
					),
			)

			if LokiConfig.Url != "" && LokiConfig.Username != "" && LokiConfig.Password != "" {
				loki := NewZapLoki(
					context.Background(), ZapLokiConfig{
						Url:          LokiConfig.Url,
						BatchMaxSize: 1000,
						BatchMaxWait: 10 * time.Second,
						Labels:       map[string]string{"app": LokiConfig.App},
						Username:     LokiConfig.Username,
						Password:     LokiConfig.Password,
					},
				)

				logger = zap.New(core, zap.Hooks(loki.Hook))
			} else {
				logger = zap.New(core)
			}
		},
	)

	return logger
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store the same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
