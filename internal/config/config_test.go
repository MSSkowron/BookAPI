package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigValid(t *testing.T) {
	configFile := createTempConfigFile(t)
	defer os.Remove(configFile)

	cfg, err := LoadConfig(configFile)
	require.NoError(t, err)

	require.Equal(t, "test_db", cfg.DatabaseURL)
	require.Equal(t, "127.0.0.1:8080", cfg.HTTPServerListenAddress)
	require.Equal(t, "test_secret", cfg.TokenSecret)
	require.Equal(t, time.Hour, cfg.TokenDuration)
}

func TestLoadConfigInvalidPath(t *testing.T) {
	configFile := createTempConfigFile(t)
	defer os.Remove(configFile)

	_, err := LoadConfig("invalid_path_config.env")
	require.ErrorIs(t, err, os.ErrNotExist)
}

func createTempConfigFile(t *testing.T) string {
	configFile := "temp_config.env"
	file, err := os.Create(configFile)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString("DATABASE_URL=test_db\n")
	require.NoError(t, err)

	_, err = file.WriteString("HTTP_SERVER_LISTEN_ADDRESS=127.0.0.1:8080\n")
	require.NoError(t, err)

	_, err = file.WriteString("TOKEN_SECRET=test_secret\n")
	require.NoError(t, err)

	_, err = file.WriteString("TOKEN_DURATION=1h\n")
	require.NoError(t, err)

	return configFile
}
