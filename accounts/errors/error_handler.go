package errors

import (
	"encoding/json"
	"io"
	"net/http"
)

// APIError represents a general API error with a status code and message.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// HandleHTTPError checks the HTTP response for errors and returns the appropriate custom error type.
func HandleHTTPError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return &ErrPermanentFailure{Detail: "failed to read response body"}
	}

	var apiErr APIError

	err = json.Unmarshal(bodyBytes, &apiErr)

	if err != nil {
		return &ErrPermanentFailure{Detail: "failed to unmarshal API error"}
	}

	apiErr.StatusCode = resp.StatusCode

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return &ErrBadRequest{Detail: apiErr.Message}
	case http.StatusNotFound:
		return &ErrNotFound{ResourceID: apiErr.Message}
	default:
		return &apiErr
	}
}
