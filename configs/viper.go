package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper(filename, filetype string, paths ...string) (*viper.Viper, error) {
	config := viper.New()

	// Bind env variables
	config.AutomaticEnv()

	// Try to load config from file
	config.SetConfigName(filename)
	config.SetConfigType(filetype)

	for _, p := range paths {
		config.AddConfigPath(p)
	}

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config viper file: %w", err)
		}
	}

	return config, nil
}

// docker run -d --name sharing-vision-golang -e APP_NAME=TestApp -e APP_ENV=production -e LOG_LEVEL=info -p 8004:3000 golang-sharing-vision:latest
