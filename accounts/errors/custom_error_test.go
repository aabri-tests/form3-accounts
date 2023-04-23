package errors_test

import (
	"fmt"
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1/accounts/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("ErrBadRequest", func(t *testing.T) {
		detail := "missing required field"
		err := &errors.ErrBadRequest{Detail: detail}
		expectedMessage := fmt.Sprintf("bad request: %s", detail)

		assert.Equal(t, expectedMessage, err.Error())
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		resourceID := "123456"
		err := &errors.ErrNotFound{ResourceID: resourceID}
		expectedMessage := fmt.Sprintf("resource not found with ID: %s", resourceID)

		assert.Equal(t, expectedMessage, err.Error())
	})

	t.Run("ErrPermanentFailure", func(t *testing.T) {
		detail := "database connection error"
		err := &errors.ErrPermanentFailure{Detail: detail}
		expectedMessage := fmt.Sprintf("permanent failure: %s", detail)

		assert.Equal(t, expectedMessage, err.Error())
	})
}
