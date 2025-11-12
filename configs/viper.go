package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper(filename, filetype string, path ...string) (*viper.Viper, error) {
	config := viper.New()

	config.SetConfigFile(filename)
	config.SetConfigType(filetype)

	for _, p := range path {
		config.AddConfigPath(p)
	}

	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config viper file : %w", err)
	}

	return config, nil
}
