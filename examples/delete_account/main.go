package main

import (
	"time"

	"github.com/aabri-assignments/form3-accounts/v1"
	"github.com/aabri-assignments/form3-accounts/v1/accounts"
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
)

func main() {
	opts := accounts.Options{
		BaseURL:      "http://accountapi:8080",
		Duration:     3 * time.Minute,
		Retries:      3,
		InitialDelay: 300 * time.Millisecond,
		Multiplier:   2,
		Factor:       0.1,
		LogLevel:     logging.LevelInfo,
	}
	// Create a new Form3 client.
	form3Client, err := form3.NewForm3(opts)
	if err != nil {
		panic(err)
	}

	// Delete an account.
	err = form3Client.Accounts().Delete("{account-id}", 0)
	if err != nil {
		panic(err)
	}
}
