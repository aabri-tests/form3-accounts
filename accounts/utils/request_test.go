package utils_test

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/utils"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	FieldA string `json:"field_a"`
	FieldB int    `json:"field_b"`
}

func TestEncodeJSONRequest(t *testing.T) {
	t.Run("Successful JSON encoding", func(t *testing.T) {
		requestBody := TestStruct{
			FieldA: "test value",
			FieldB: 42,
		}

		reader, err := utils.EncodeJSONRequest(requestBody)

		assert.NoError(t, err)

		buf := new(strings.Builder)
		_, err = io.Copy(buf, reader)
		assert.NoError(t, err)

		var decoded TestStruct
		err = json.Unmarshal([]byte(buf.String()), &decoded)

		assert.NoError(t, err)
		assert.Equal(t, requestBody, decoded)
	})

	t.Run("Error on invalid JSON encoding", func(t *testing.T) {
		requestBody := make(chan int) // An unsupported type for JSON encoding
		reader, err := utils.EncodeJSONRequest(requestBody)

		assert.Error(t, err)
		assert.Nil(t, reader)
	})
}
