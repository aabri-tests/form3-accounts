package main

import (
	"fmt"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1"
	"github.com/aabri-assignments/form3-accounts/v1/accounts"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
	"github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
	"github.com/google/uuid"
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

	// Create a new account.
	switched := true
	status := "confirmed"
	jointAccount := false
	country := "GB"
	out := true
	classification := "Personal"
	id := uuid.New()
	orgID := uuid.New()
	account := &models.AccountData{
		Type:           "accounts",
		ID:             id.String(),
		OrganisationID: orgID.String(),
		Attributes: &models.AccountAttributes{
			Country:                 &country,
			BaseCurrency:            "GBP",
			BankID:                  "123456",
			BankIDCode:              "GBDSC",
			Bic:                     "EXMPLGB2XXX",
			AccountClassification:   &classification,
			AccountMatchingOptOut:   &out,
			AccountNumber:           "12345678",
			AlternativeNames:        []string{"AltName1", "AltName2"},
			Iban:                    "GB11NWBK56000123456789",
			JointAccount:            &jointAccount,
			Name:                    []string{"John Smith"},
			SecondaryIdentification: "1234",
			Status:                  &status,
			Switched:                &switched,
		},
	}

	createdAccount, err := form3Client.Accounts().Create(account)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created account: %+v\n", createdAccount)
}
