package modules

import (
	"context"

	"github.com/enrichoalkalas01/test-sharing-vision-golang/pkg/database"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var DatabaseModule = fx.Module("database",
	fx.Provide(NewMySQL),
	fx.Invoke(registerDatabaseLifecycle),
)

func NewMySQL(config *viper.Viper, log *zap.Logger) (*gorm.DB, error) {
	return database.NewMySQL(config, log)
}

func registerDatabaseLifecycle(lc fx.Lifecycle, db *gorm.DB, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection...")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
}
