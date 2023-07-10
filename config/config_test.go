package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigValid(t *testing.T) {
	configFile := createTempConfigFile()
	defer os.Remove(configFile)

	cfg, err := LoadConfig(configFile)
	assert.NoError(t, err)

	assert.Equal(t, "test_db", cfg.DBSource)
	assert.Equal(t, "127.0.0.1:8080", cfg.HTTPServerListenAddress)
	assert.Equal(t, "test_secret", cfg.TokenSecret)
	assert.Equal(t, time.Hour, cfg.TokenDuration)
}

func TestLoadConfigInvalidPath(t *testing.T) {
	configFile := createTempConfigFile()
	defer os.Remove(configFile)

	_, err := LoadConfig("invalid_path_config.env")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func createTempConfigFile() string {
	configFile := "temp_config.env"
	file, _ := os.Create(configFile)
	defer file.Close()

	file.WriteString("DB_SOURCE=test_db\n")
	file.WriteString("HTTP_SERVER_LISTEN_ADDRESS=127.0.0.1:8080\n")
	file.WriteString("TOKEN_SECRET=test_secret\n")
	file.WriteString("TOKEN_DURATION=1h\n")

	return configFile
}
