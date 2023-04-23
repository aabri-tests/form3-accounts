package accounts_test

import (
	"context"
	"testing"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1/accounts"
	mocks_retry "github.com/aabri-assignments/form3-accounts/v1/accounts/mocks"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/retry"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport/http"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type MockFailingTransport struct{}

func (m *MockFailingTransport) Create(ctx context.Context, req *utils.CreateAccountRequest) (*utils.CreateAccountResponse, error) {
	return nil, nil
}

func (m *MockFailingTransport) Delete(ctx context.Context, req *utils.DeleteAccountRequest) error {
	return nil
}

func (m *MockFailingTransport) Fetch(ctx context.Context, accountID string) (*utils.FetchAccountResponse, error) {
	return nil, errors.New("Failing operation")
}

func TestNewClient(t *testing.T) {
	opts := accounts.Options{
		BaseURL:      "https://api.example.com",
		Duration:     5 * time.Second,
		Retries:      3,
		InitialDelay: 2 * time.Second,
		Multiplier:   2.0,
		Factor:       1.5,
		LogLevel:     logging.LevelInfo,
	}
	client := accounts.New(opts)

	// Check if httpTransport has the correct BaseURL
	httpTransport, ok := client.GetTransport().(*http.Transport)
	assert.True(t, ok, "Expected HTTPTransport")
	assert.Equal(t, opts.BaseURL, httpTransport.BaseURL, "BaseURL does not match")

	// Check if Retry is an ExponentialBackOff instance
	backOff, ok := client.GetRetry().(*retry.ExponentialBackOff)
	assert.True(t, ok, "Expected ExponentialBackOff")
	assert.Equal(t, opts.Retries, backOff.MaxRetries, "MaxRetries does not match")
	assert.Equal(t, opts.InitialDelay, backOff.InitialDelay, "InitialDelay does not match")
	assert.Equal(t, float64(opts.Multiplier), backOff.Multiplier, "Multiplier does not match")
	assert.Equal(t, opts.Factor, backOff.RandFactor, "Factor does not match")

	// Check if Logger has the correct LogLevel
	logger := client.GetLogger().(*logging.Leveled)
	assert.Equal(t, opts.LogLevel, logger.Level, "LogLevel does not match")
}
func TestClient(t *testing.T) {
	mockTransport := &mocks_retry.MockTransport{}
	mockRetries := &mocks_retry.MockRetrier{}
	client := &accounts.AccountClient{
		Transport: mockTransport,
		Retry:     mockRetries,
		Logger:    &logging.Leveled{Level: logging.LevelError},
	}

	// Test Create method
	account := &models.AccountData{
		ID:             "some-id",
		OrganisationID: "some-org-id",
	}

	t.Run("Create", func(t *testing.T) {
		mockTransport.On("Create", context.Background(), &utils.CreateAccountRequest{Data: *account}).Return(&utils.CreateAccountResponse{Data: *account}, nil)
		createdAccount, err := client.Create(account)
		assert.NoError(t, err)
		assert.Equal(t, account, createdAccount)
		mockTransport.AssertExpectations(t)
	})
	t.Run("Fetch", func(t *testing.T) {
		// Test Fetch method
		accountID := "some-id"
		mockTransport.On("Fetch", context.Background(), accountID).Return(&utils.FetchAccountResponse{Data: *account}, nil)
		fetchedAccount, err := client.Fetch(accountID)
		assert.NoError(t, err)
		assert.Equal(t, account, fetchedAccount)
		mockTransport.AssertExpectations(t)
	})
	t.Run("Delete", func(t *testing.T) {
		// Test Delete method
		deleteAccountID := "some-id"
		version := int64(1)
		mockTransport.On("Delete", context.Background(), &utils.DeleteAccountRequest{ID: deleteAccountID, Version: version}).Return(nil)
		err := client.Delete(deleteAccountID, version)
		assert.NoError(t, err)
		mockTransport.AssertExpectations(t)
	})

}
func TestCreateWithError(t *testing.T) {
	mockTransport := &mocks_retry.MockTransport{}
	mockRetries := &mocks_retry.MockRetrier{}
	client := &accounts.AccountClient{
		Transport: mockTransport,
		Retry:     mockRetries,
		Logger:    &logging.Leveled{Level: logging.LevelError},
	}

	// Test Create method
	account := &models.AccountData{
		ID:             "some-id",
		OrganisationID: "some-org-id",
	}
	t.Run("Create With Error", func(t *testing.T) {
		// Test Create method with an error
		mockTransport.On("Create", context.Background(), &utils.CreateAccountRequest{Data: *account}).Return((*utils.CreateAccountResponse)(nil), errors.New("create error"))
		_, err := client.Create(account)
		assert.Error(t, err)
		mockTransport.AssertExpectations(t)
	})
}
func TestClientFetch(t *testing.T) {
	opts := accounts.Options{
		BaseURL:      "https://api.example.com",
		Duration:     5 * time.Second,
		Retries:      3,
		InitialDelay: 2 * time.Second,
		Multiplier:   2.0,
		Factor:       1.5,
		LogLevel:     logging.LevelInfo,
	}
	client := accounts.New(opts)

	// Use the custom transport that simulates a failing operation
	client.SetTransport(&MockFailingTransport{})

	accountID := "test-account-id"
	// Call the Fetch method that uses the retry mechanism
	_, err := client.Fetch(accountID)

	// Assert that the error is not nil and the error message is as expected
	assert.NotNil(t, err)
	assert.Equal(t, "Failing operation", err.Error())
}

func TestClientSetRetry(t *testing.T) {
	mockTransport := &mocks_retry.MockTransport{}
	client := &accounts.AccountClient{
		Transport: mockTransport,
		Retry:     nil,
		Logger:    &logging.Leveled{Level: logging.LevelError},
	}
	// Test SetRetry method
	mockRetries := &mocks_retry.MockRetrier{}
	client.SetRetry(mockRetries)

	// Check if the Retry has been set correctly
	assert.Equal(t, mockRetries, client.Retry, "SetRetry does not set the expected Retrier")
}
func TestClientSetLogger(t *testing.T) {
	mockTransport := &mocks_retry.MockTransport{}
	client := &accounts.AccountClient{
		Transport: mockTransport,
		Retry:     nil,
		Logger:    &logging.Leveled{Level: logging.LevelError},
	}
	// Test SetLogger method
	mockLogger := &logging.Leveled{Level: logging.LevelInfo}
	client.SetLogger(mockLogger)

	// Check if the Logger has been set correctly
	assert.Equal(t, mockLogger, client.Logger, "SetLogger does not set the expected LeveledLogger")
}
