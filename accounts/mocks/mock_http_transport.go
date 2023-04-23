package mocks

import (
	"context"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/stretchr/testify/mock"
)

type MockTransport struct {
	transport.Transport
	mock.Mock
}

func (m *MockTransport) Create(ctx context.Context, req *utils.CreateAccountRequest) (*utils.CreateAccountResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*utils.CreateAccountResponse), args.Error(1)
}

func (m *MockTransport) Fetch(ctx context.Context, accountID string) (*utils.FetchAccountResponse, error) {
	args := m.Called(ctx, accountID)

	return args.Get(0).(*utils.FetchAccountResponse), args.Error(1)
}

func (m *MockTransport) Delete(ctx context.Context, req *utils.DeleteAccountRequest) error {
	args := m.Called(ctx, req)

	return args.Error(0)
}
