package modules

import (
	"github.com/enrichoalkalas01/test-sharing-vision-golang/configs"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(NewViper),
)

func NewViper() (*viper.Viper, error) {
	return configs.NewViper(".env", "env", ".", "../../")
}
