package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
	assert "github.com/stretchr/testify/require"
)

func TestLeveledLogger(t *testing.T) {
	var buf bytes.Buffer

	logger := &logging.Leveled{
		Level:          logging.LevelDebug,
		StdoutOverride: &buf,
		StderrOverride: &buf,
	}

	t.Run("Debugf logs a message", func(t *testing.T) {
		buf.Reset()
		logger.Debugf("debug %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[DEBUG] debug message\n"))
	})
	t.Run("Infof logs a message", func(t *testing.T) {
		buf.Reset()
		logger.Infof("info %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[INFO] info message\n"))
	})
	t.Run("Warnf logs a message", func(t *testing.T) {
		buf.Reset()
		logger.Warnf("warn %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[WARN] warn message\n"))
	})
	t.Run("Errorf logs a message", func(t *testing.T) {
		buf.Reset()
		logger.Errorf("error %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[ERROR] error message\n"))
	})
	t.Run("Logger respects log level", func(t *testing.T) {
		logger.Level = logging.LevelWarn
		buf.Reset()

		logger.Debugf("debug %s", "message")
		assert.False(t, strings.Contains(buf.String(), "[DEBUG] debug message\n"))

		logger.Infof("info %s", "message")
		assert.False(t, strings.Contains(buf.String(), "[INFO] info message\n"))

		logger.Warnf("warn %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[WARN] warn message\n"))

		logger.Errorf("error %s", "message")
		assert.True(t, strings.Contains(buf.String(), "[ERROR] error message\n"))
	})
}

func TestLeveledLogger_NoOverrides(t *testing.T) {
	logger := &logging.Leveled{
		Level:          logging.LevelDebug,
		StderrOverride: nil,
		StdoutOverride: nil,
	}

	t.Run("Debugf logs a message without StdoutOverride", func(t *testing.T) {
		logger.Debugf("debug %s", "message")
	})

	t.Run("Errorf logs a message without StderrOverride", func(t *testing.T) {
		logger.Errorf("error %s", "message")
	})
}
