package utils

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
)

// CreateAccountRequest represents the request structure for creating an account.
type CreateAccountRequest struct {
	Data models.AccountData `json:"data"`
}

// DeleteAccountRequest represents the request structure for deleting an account.
type DeleteAccountRequest struct {
	ID      string
	Version int64
}

// EncodeJSONRequest encodes a JSON request body.
func EncodeJSONRequest(requestBody interface{}) (io.Reader, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}
