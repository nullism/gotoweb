package logging

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Configure(t *testing.T) {
	// Test case 1: Set log level to Debug
	Configure(slog.LevelError)
	assert.Equal(t, slog.LevelError, logLevel.Level())

	// Test case 1: Set log level to Debug
	Configure(slog.LevelDebug)
	assert.Equal(t, slog.LevelDebug, logLevel.Level())
}

func Test_GetLogger(t *testing.T) {
	// Test case 1: Get logger
	logger := GetLogger()
	assert.NotNil(t, logger)
	assert.Equal(t, logger, &CLIHandler{slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &logLevel})})
}
