package accounts

import (
	"context"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/retry"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport/http"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
)

const basePath = "/v1/organisation/accounts/"

type Client interface {
	Create(account *models.AccountData) (*models.AccountData, error)
	Fetch(accountID string) (*models.AccountData, error)
	Delete(accountID string, version int64) error
	GetRetry() retry.Retrier
	SetRetry(r retry.Retrier)
	GetTransport() transport.Transport
	SetTransport(t transport.Transport)
	GetLogger() logging.LeveledLogger
	SetLogger(l logging.LeveledLogger)
}

type Options struct {
	BaseURL      string
	Duration     time.Duration
	Retries      int
	InitialDelay time.Duration
	Multiplier   int
	Factor       float64
	LogLevel     logging.Level
}
type AccountClient struct {
	Transport transport.Transport
	Retry     retry.Retrier
	Logger    logging.LeveledLogger
}

func New(opt Options) Client {
	httpTransport := http.New(opt.BaseURL, basePath)
	backOff := retry.NewExponentialBackOff(opt.Duration, opt.Retries, opt.InitialDelay, float64(opt.Multiplier), opt.Factor)

	return &AccountClient{
		Transport: httpTransport,
		Retry:     backOff,
		Logger: &logging.Leveled{
			Level: opt.LogLevel,
		},
	}
}
func (c *AccountClient) Create(account *models.AccountData) (*models.AccountData, error) {
	req := &utils.CreateAccountRequest{Data: *account}

	var resp *utils.CreateAccountResponse

	operation := func() error {
		var err error
		resp, err = c.GetTransport().Create(context.Background(), req)

		return err
	}

	if err := retry.Retry(operation, c.Retry, c.Logger); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
func (c *AccountClient) Fetch(accountID string) (*models.AccountData, error) {
	var resp *utils.FetchAccountResponse

	operation := func() error {
		var err error

		resp, err = c.GetTransport().Fetch(context.Background(), accountID)

		return err
	}

	if err := retry.Retry(operation, c.Retry, c.Logger); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
func (c *AccountClient) Delete(accountID string, version int64) error {
	req := &utils.DeleteAccountRequest{ID: accountID, Version: version}

	operation := func() error {
		c.Logger.Infof("Attempting operation...")

		return c.Transport.Delete(context.Background(), req)
	}

	return retry.Retry(operation, c.Retry, c.Logger)
}
func (c *AccountClient) GetTransport() transport.Transport {
	return c.Transport
}
func (c *AccountClient) SetTransport(t transport.Transport) {
	c.Transport = t
}
func (c *AccountClient) GetRetry() retry.Retrier {
	return c.Retry
}
func (c *AccountClient) SetRetry(r retry.Retrier) {
	c.Retry = r
}
func (c *AccountClient) GetLogger() logging.LeveledLogger {
	return c.Logger
}
func (c *AccountClient) SetLogger(l logging.LeveledLogger) {
	c.Logger = l
}
