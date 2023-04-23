package form3_test

import (
	"testing"

	"github.com/aabri-assignments/form3-accounts/v1"
	"github.com/aabri-assignments/form3-accounts/v1/accounts"
	"github.com/stretchr/testify/assert"
)

func TestNewForm3(t *testing.T) {
	options := accounts.Options{
		BaseURL: "http://accountapi:8080",
	}

	form3Client, err := form3.NewForm3(options)

	assert.NoError(t, err, "NewForm3 should not return an error")
	assert.NotNil(t, form3Client, "NewForm3 should return a valid Form3 client")
}

func TestAccounts(t *testing.T) {
	options := accounts.Options{
		BaseURL: "http://accountapi:8080",
	}

	form3Client, _ := form3.NewForm3(options)
	accountsClient := form3Client.Accounts()

	assert.NotNil(t, accountsClient, "Accounts should return a valid accounts.Client")
}
