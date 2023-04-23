package main

import (
	"encoding/json"
	"fmt"
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

	// Fetch an account.
	fetchedAccount, err := form3Client.Accounts().Fetch("{account-id}")
	if err != nil {
		panic(err)
	}

	fetchedAccountJSON, err := json.Marshal(fetchedAccount)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(fetchedAccountJSON))
}
