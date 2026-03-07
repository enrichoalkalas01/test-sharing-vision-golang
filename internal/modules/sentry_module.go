package modules

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var SentryModule = fx.Module("sentry",
	fx.Invoke(initSentry),
)

func initSentry(config *viper.Viper) error {
	dsn := config.GetString("SENTRY_DSN")
	if dsn == "" {
		return nil
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      config.GetString("APP_ENV"),
		Release:          config.GetString("APP_VERSION"),
		TracesSampleRate: 1.0,
	}); err != nil {
		return fmt.Errorf("sentry initialization failed: %w", err)
	}

	return nil
}
