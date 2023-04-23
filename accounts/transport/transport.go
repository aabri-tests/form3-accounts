package transport

import (
	"context"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
)

// Transport defines the interface for the transport layer of the client library.
type Transport interface {
	Create(context context.Context, req *utils.CreateAccountRequest) (*utils.CreateAccountResponse, error)
	Fetch(context context.Context, accountID string) (*utils.FetchAccountResponse, error)
	Delete(context context.Context, req *utils.DeleteAccountRequest) error
}
