package mocks

import (
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
	"github.com/stretchr/testify/mock"
)

// MockLeveledLogger is a mock implementation of LeveledLogger for testing purposes.
type MockLeveledLogger struct {
	mock.Mock
}

// Debugf mocks the Debugf method of LeveledLogger.
func (m *MockLeveledLogger) Debugf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Infof mocks the Infof method of LeveledLogger.
func (m *MockLeveledLogger) Infof(format string, v ...interface{}) {
	m.Called(format, v)
}

// Warnf mocks the Warnf method of LeveledLogger.
func (m *MockLeveledLogger) Warnf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Errorf mocks the Errorf method of LeveledLogger.
func (m *MockLeveledLogger) Errorf(format string, v ...interface{}) {
	m.Called(format, v)
}

// Level mocks the Level method of LeveledLogger.
func (m *MockLeveledLogger) Level() logging.Level {
	args := m.Called()
	return args.Get(0).(logging.Level)
}

// SetLevel mocks the SetLevel method of LeveledLogger.
func (m *MockLeveledLogger) SetLevel(level logging.Level) {
	m.Called(level)
}
