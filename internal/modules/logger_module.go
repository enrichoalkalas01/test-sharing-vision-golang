package modules

import (
	"context"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var LoggerModule = fx.Module("logger",
	fx.Provide(NewLogger),
	fx.Provide(ExtractZapLogger),
	fx.Invoke(registerLoggerLifecycle),
)

func NewLogger() (*logger.Logger, error) {
	return logger.NewDefaultLogger()
}

func ExtractZapLogger(l *logger.Logger) *zap.Logger {
	return l.Logger
}

func registerLoggerLifecycle(lc fx.Lifecycle, l *logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = l.Sync()
			return nil
		},
	})
}
