package errors_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	errs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errs.New("read error")
}

func TestHandleHTTPError(t *testing.T) {
	testCases := []struct {
		name           string
		response       *http.Response
		expectedResult error
	}{
		{
			name: "success",
			response: &http.Response{
				StatusCode: http.StatusOK,
			},
			expectedResult: nil,
		},
		{
			name: "bad_request",
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "missing required field"}`))),
			},
			expectedResult: &errors.ErrBadRequest{Detail: "missing required field"},
		},
		{
			name: "not_found",
			response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "123456"}`))),
			},
			expectedResult: &errors.ErrNotFound{ResourceID: "123456"},
		},
		{
			name: "other_error",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "internal server error"}`))),
			},
			expectedResult: &errors.APIError{StatusCode: http.StatusInternalServerError, Message: "internal server error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.HandleHTTPError(tc.response)
			assert.Equal(t, tc.expectedResult, err)
		})
	}
}
func TestAPIError_Error(t *testing.T) {
	apiError := &errors.APIError{
		StatusCode: 500,
		Message:    "internal server error",
	}

	expectedMessage := "internal server error"
	assert.Equal(t, expectedMessage, apiError.Error())
}

func TestHandleHTTPError_FailedToReadResponseBody(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(&errorReader{}),
	}

	err := errors.HandleHTTPError(response)
	assert.Error(t, err)
	assert.IsType(t, &errors.ErrPermanentFailure{}, err)
	assert.Equal(t, "permanent failure: failed to read response body", err.Error())
}

func TestHandleHTTPError_FailedToUnmarshalAPIError(t *testing.T) {
	invalidJSON := `{"message": "internal server error",`

	response := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewReader([]byte(invalidJSON))),
	}

	err := errors.HandleHTTPError(response)
	assert.Error(t, err)
	assert.IsType(t, &errors.ErrPermanentFailure{}, err)
	assert.Equal(t, "permanent failure: failed to unmarshal API error", err.Error())
}
