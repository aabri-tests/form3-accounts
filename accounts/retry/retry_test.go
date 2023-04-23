package retry_test

import (
	"testing"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/retry"
	errs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Debugf(format string, v ...interface{})    {}
func (m *mockLogger) Warnf(format string, v ...interface{})     {}
func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}

func TestExponentialBackOff(t *testing.T) {
	t.Run("NewExponentialBackOff uses default values", func(t *testing.T) {
		backOff := retry.NewExponentialBackOff(0, 0, 0, 0, 0)
		assert.Equal(t, 5*time.Minute, backOff.MaxElapsedTime)
		assert.Equal(t, 5, backOff.MaxRetries)
		assert.Equal(t, 0*time.Millisecond, backOff.InitialDelay)
		assert.Equal(t, 2.0, backOff.Multiplier)
		assert.Equal(t, 0.1, backOff.RandFactor)
	})
	t.Run("NextBackOff respects MaxRetries", func(t *testing.T) {
		backOff := retry.NewExponentialBackOff(5*time.Minute, 3, 0*time.Millisecond, 2, 0.1)
		for i := 0; i < 3; i++ {
			assert.NotEqual(t, -1, backOff.NextBackOff())
		}
		assert.Equal(t, time.Duration(-1), backOff.NextBackOff())
	})
	t.Run("NextBackOff respects MaxElapsedTime", func(t *testing.T) {
		backOff := retry.NewExponentialBackOff(10*time.Millisecond, 5, 1*time.Millisecond, 2, 0.1)
		foundMaxElapsedTime := false
		for i := 0; i < 5; i++ {
			delay := backOff.NextBackOff()
			if delay == time.Duration(-1) {
				foundMaxElapsedTime = true
			} else {
				assert.LessOrEqual(t, delay, backOff.MaxElapsedTime)
			}
		}
		assert.True(t, foundMaxElapsedTime, "Expected to find delay with -1 when MaxElapsedTime is reached")
	})
	t.Run("Attempt and Reset work correctly", func(t *testing.T) {
		backOff := retry.NewExponentialBackOff(5*time.Minute, 3, 10*time.Millisecond, 2, 0.1)

		assert.Equal(t, 0, backOff.Attempt())

		// Simulate 2 attempts
		backOff.NextBackOff()
		backOff.NextBackOff()

		assert.Equal(t, 2, backOff.Attempt())
		assert.Equal(t, 1, backOff.RemainingRetries())

		// Reset the backOff
		backOff.Reset()

		// Check if attempt and remaining retries have been reset
		assert.Equal(t, 0, backOff.Attempt())
		assert.Equal(t, 3, backOff.RemainingRetries())
	})
}

func TestRetry(t *testing.T) {
	logger := &mockLogger{}
	temporaryError := errs.New("temporary error")
	permanentError := &errors.ErrPermanentFailure{Detail: errs.New("permanent error").Error()}

	t.Run("Retries on temporary errors", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			return temporaryError
		}
		backOff := retry.NewExponentialBackOff(5*time.Minute, 3, 10*time.Millisecond, 2, 0.1)
		err := retry.Retry(operation, backOff, logger)
		assert.Error(t, err)
		assert.Equal(t, 4, attempts)
	})
	t.Run("Does not retry on permanent errors", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			return permanentError
		}
		backOff := retry.NewExponentialBackOff(5*time.Minute, 3, 10*time.Millisecond, 2, 0.1)
		err := retry.Retry(operation, backOff, logger)
		assert.Error(t, err)
		assert.Equal(t, 1, attempts)
	})
	t.Run("Stops retrying when successful", func(t *testing.T) {
		attempts := 0
		operation := func() error {
			attempts++
			if attempts == 3 {
				return nil
			}
			return temporaryError
		}
		backOff := retry.NewExponentialBackOff(5*time.Minute, 5, 10*time.Millisecond, 2, 0.1)
		err := retry.Retry(operation, backOff, logger)
		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})
	t.Run("Respects remaining retries", func(t *testing.T) {
		operation := func() error {
			return temporaryError
		}
		backOff := retry.NewExponentialBackOff(5*time.Minute, 2, 10*time.Millisecond, 2, 0.1)
		err := retry.Retry(operation, backOff, logger)
		assert.Error(t, err)
		assert.Equal(t, 0, backOff.RemainingRetries())
	})
}
