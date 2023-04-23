package utils

import (
	"encoding/json"
	"net/http"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
)

// CreateAccountResponse represents the response structure for creating an account.
type CreateAccountResponse struct {
	Data models.AccountData `json:"data"`
}

// FetchAccountResponse represents the response structure for fetching an account.
type FetchAccountResponse struct {
	Data models.AccountData `json:"data"`
}

// DecodeJSONResponse decodes a JSON response body.
func DecodeJSONResponse(resp *http.Response, responseBody interface{}) error {
	return json.NewDecoder(resp.Body).Decode(responseBody)
}

// CloseResponse safely closes the response body after reading.
func CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
