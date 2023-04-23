package retry

import (
	errs "errors"
	"math"
	"math/rand"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
)

const (
	defaultDuration          = 5 * time.Minute
	defaultDelay             = 00 * time.Millisecond
	defaultMultiplier        = 2
	defaultFactor            = 0.1
	defaultMaxNetworkRetries = 5
)

// Retrier is an interface that represents a backoff retry strategy.
type Retrier interface {
	Attempt() int
	NextBackOff() time.Duration
	Reset()
	RemainingRetries() int
}

// ExponentialBackOff implements a backoff policy for retrying an operation using exponential backoff.
type ExponentialBackOff struct {
	MaxElapsedTime   time.Duration
	MaxRetries       int
	InitialDelay     time.Duration
	Multiplier       float64
	RandFactor       float64
	attempt          int
	remainingRetries int
}

// NewExponentialBackOff creates a new ExponentialBackOff with the specified parameters.
func NewExponentialBackOff(maxElapsedTime time.Duration, maxNetworkRetries int, initialDelay time.Duration, multiplier, randFactor float64) *ExponentialBackOff {
	if maxElapsedTime == 0 {
		maxElapsedTime = defaultDuration
	}

	if maxNetworkRetries == 0 {
		maxNetworkRetries = defaultMaxNetworkRetries // set default maxNetworkRetries to 3
	}

	if initialDelay == 0 {
		initialDelay = defaultDelay
	}

	if multiplier == 0 {
		multiplier = defaultMultiplier
	}

	if randFactor == 0 {
		randFactor = defaultFactor
	}

	return &ExponentialBackOff{
		MaxElapsedTime:   maxElapsedTime,
		MaxRetries:       maxNetworkRetries,
		InitialDelay:     initialDelay,
		Multiplier:       multiplier,
		RandFactor:       randFactor,
		attempt:          0,
		remainingRetries: maxNetworkRetries,
	}
}

// Attempt returns the current attempt number.
func (b *ExponentialBackOff) Attempt() int {
	return b.attempt
}

func (b *ExponentialBackOff) RemainingRetries() int {
	return b.remainingRetries
}

// NextBackOff calculates the next backoff interval based on the exponential backoff strategy.
func (b *ExponentialBackOff) NextBackOff() time.Duration {
	if b.attempt >= b.MaxRetries {
		return -1
	}

	b.attempt++
	b.remainingRetries--

	backoff := float64(b.InitialDelay) * math.Pow(b.Multiplier, float64(b.attempt-1))
	jitter := (rand.Float64()*2 - 1) * b.RandFactor * backoff //nolint:gosec

	delay := time.Duration(backoff + jitter)
	if delay > b.MaxElapsedTime {
		return -1
	}

	return delay
}

// Reset resets the backoff attempts counter.
func (b *ExponentialBackOff) Reset() {
	b.attempt = 0
	b.remainingRetries = b.MaxRetries
}

// Retry retries the provided function using the provided Retries strategy.
func Retry(operation func() error, retries Retrier, logger logging.LeveledLogger) error {
	var err error

	for {
		err = operation()
		if err == nil {
			break
		}

		var permErr *errors.ErrPermanentFailure

		if errs.As(err, &permErr) {
			break
		}

		delay := retries.NextBackOff()
		if delay == -1 {
			break
		}

		logger.Infof("Retrying..")
		logger.Debugf("Remaining retries: %d", retries.RemainingRetries())

		time.Sleep(delay)
	}

	return err
}
