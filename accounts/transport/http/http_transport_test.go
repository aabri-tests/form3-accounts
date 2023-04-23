package http_test

import (
	"fmt"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/transport/http"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type HttpTestSuite struct {
	t   *testing.T
	ctx context.Context
}

func TestHttpTransport(t *testing.T) {
	suite := &HttpTestSuite{t: t, ctx: context.Background()}
	t.Run("TestCreate", suite.TestCreate)
	t.Run("TestCreateWithError", suite.TestCreateWithError)
	t.Run("TestFetch", suite.TestFetch)
	t.Run("TestDelete", suite.TestDelete)
}
func (suite *HttpTestSuite) TestCreate(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	baseURL := server.URL
	basePath := "/create"
	transport := http.New(baseURL, basePath)
	req := &utils.CreateAccountRequest{
		Data: models.AccountData{
			Type: "accounts",
			Attributes: &models.AccountAttributes{
				AccountNumber: "123456",
			},
		},
	}

	resp, err := transport.Create(suite.ctx, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.Data.ID)
	assert.Equal(t, "accounts", resp.Data.Type)
	assert.Equal(t, "123456", resp.Data.Attributes.AccountNumber)
}
func (suite *HttpTestSuite) TestCreateWithError(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	baseURL := server.URL
	basePath := "/create"
	transport := http.New(baseURL, basePath)
	req := &utils.CreateAccountRequest{
		Data: models.AccountData{
			Type: "accounts",
			Attributes: &models.AccountAttributes{
				AccountNumber: "123456",
			},
		},
	}

	// Test with valid request
	resp, err := transport.Create(suite.ctx, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.Equal(t, "test-id", resp.Data.ID)

	// Test with invalid request
	server.Close()
	server = httptest.NewServer(http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		w.WriteHeader(http2.StatusOK)
		fmt.Fprint(w, `{"invalid": "json"`) // invalid JSON
	}))
	defer server.Close()

	transport = http.New(server.URL, basePath)
	_, err = transport.Create(suite.ctx, req)
	assert.Error(t, err)
	assert.IsType(t, &errors.ErrPermanentFailure{}, err)
}
func (suite *HttpTestSuite) TestFetch(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	baseURL := server.URL
	basePath := "/fetch/"
	transport := http.New(baseURL, basePath)
	resp, err := transport.Fetch(suite.ctx, "test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.Data.ID)
	assert.Equal(t, "accounts", resp.Data.Type)
	assert.Equal(t, "123456", resp.Data.Attributes.AccountNumber)
}
func (suite *HttpTestSuite) TestDelete(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	baseURL := server.URL
	basePath := "/delete/"
	transport := http.New(baseURL, basePath)
	req := &utils.DeleteAccountRequest{
		ID:      "test-id",
		Version: 1,
	}

	err := transport.Delete(suite.ctx, req)
	assert.NoError(t, err)
}

func createMockServer() *httptest.Server {
	server := httptest.NewServer(http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		switch r.URL.Path {
		case "/create":
			resp := `{"data": {"id": "test-id", "type": "accounts", "attributes": {"account_number": "123456"}}}`
			fmt.Fprint(w, resp)
		case "/fetch/test-id":
			resp := `{"data": {"id": "test-id", "type": "accounts", "attributes": {"account_number": "123456"}}}`
			fmt.Fprint(w, resp)
		case "/delete/test-id":
			w.WriteHeader(http2.StatusNoContent)
		default:
			w.WriteHeader(http2.StatusNotFound)
			resp := `{"message": "Resource not found"}`
			fmt.Fprint(w, resp)
		}
	}))
	return server
}
