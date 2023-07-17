package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or envirionment variables.
type Config struct {
	DatabaseURL             string        `mapstructure:"DATABASE_URL"`
	HTTPServerListenAddress string        `mapstructure:"HTTP_SERVER_LISTEN_ADDRESS"`
	TokenSecret             string        `mapstructure:"TOKEN_SECRET"`
	TokenDuration           time.Duration `mapstructure:"TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
// Config values specified in the config file are overwritten with environment variables if set.
func LoadConfig(path string) (Config, error) {
	config := Config{}

	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
