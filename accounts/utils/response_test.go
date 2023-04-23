package utils_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/stretchr/testify/assert"
)

func TestDecodeJSONResponse(t *testing.T) {
	t.Run("Successful JSON decoding", func(t *testing.T) {
		requestBody := TestStruct{
			FieldA: "test value",
			FieldB: 42,
		}

		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader(jsonData)),
		}

		var decoded TestStruct
		err = utils.DecodeJSONResponse(resp, &decoded)

		assert.NoError(t, err)
		assert.Equal(t, requestBody, decoded)
	})

	t.Run("Error on invalid JSON decoding", func(t *testing.T) {
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte(`{invalid json}`))),
		}

		var decoded TestStruct
		err := utils.DecodeJSONResponse(resp, &decoded)

		assert.Error(t, err)
	})
}

func TestCloseResponse(t *testing.T) {
	t.Run("Close response without error", func(t *testing.T) {
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte(`{"field_a": "test value", "field_b": 42}`))),
		}
		assert.NotNil(t, resp.Body)

		var readErr error
		var writeBuf bytes.Buffer
		_, copyErr := io.Copy(&writeBuf, resp.Body)
		assert.NoError(t, copyErr)

		utils.CloseResponse(resp)

		_, readErr = resp.Body.Read(make([]byte, 1))
		assert.Error(t, readErr)
		assert.Equal(t, "EOF", readErr.Error())
	})
}
