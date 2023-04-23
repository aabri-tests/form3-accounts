package accounts

import (
	"math/rand"
	"time"

	"github.com/aabri-assignments/form3-accounts/v1/accounts"
	"github.com/aabri-assignments/form3-accounts/v1/accounts/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	testBaseURL = "http://accountapi:8080"
)

var _ = DescribeAccount("[Create]", func() {
	var client accounts.Client
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     10 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
	})

	args := func(country, baseCurrency, bankID, bankIDCode, bic string, name []string, accountNumber, iban string, expectError bool, classification string) {
		version := int64(0)
		account := &models.AccountData{
			Type:           "accounts",
			ID:             uuid.New().String(),
			OrganisationID: uuid.New().String(),
			Attributes: &models.AccountAttributes{
				Country:       &country,
				BaseCurrency:  baseCurrency,
				BankID:        bankID,
				BankIDCode:    bankIDCode,
				Bic:           bic,
				Name:          name,
				AccountNumber: accountNumber,
				Iban:          iban,
			},
			Version: &version,
		}
		if classification != "" {
			account.Attributes.AccountClassification = &classification
		}
		createdAccount, err := client.Create(account)

		if expectError {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(err).NotTo(HaveOccurred())
			Expect(createdAccount.ID).NotTo(BeEmpty())
			Expect(*createdAccount.Attributes.Country).To(Equal(country))
			Expect(createdAccount.Attributes.BaseCurrency).To(Equal(baseCurrency))
			Expect(createdAccount.Attributes.BankID).To(Equal(bankID))
			Expect(createdAccount.Attributes.BankIDCode).To(Equal(bankIDCode))
			Expect(createdAccount.Attributes.Bic).To(Equal(bic))
			Expect(createdAccount.Attributes.Name).To(Equal(name))
			Expect(createdAccount.Attributes.AccountNumber).To(Equal(accountNumber))
			Expect(createdAccount.Attributes.Iban).To(Equal(iban))
			if classification != "" {
				Expect(*createdAccount.Attributes.AccountClassification).To(Equal(classification))
			}
		}
		if createdAccount != nil && createdAccount.Version != nil {
			err := client.Delete(createdAccount.ID, *createdAccount.Version)
			Expect(err).NotTo(HaveOccurred())
		}
	}
	DescribeTable("Create Account", args,
		Entry("should create a UK Account With Confirmation of Payee",
			"GB", "GBP", "123456", "GBDSC", "EXMPLGB2XXX", []string{"Test Account"}, "12345678", "GB28NWBK12345678901234", false, "Personal"),
		Entry("should Create a UK Account Without Confirmation of Payee",
			"GB", "GBP", "123456", "GBDSC", "EXMPLGB2XXX", []string{"Test Account"}, "12345678", "GB28NWBK12345678901234", false, ""),
		Entry("should fail with invalid account number",
			"GB", "GBP", "123456", "GBDSC", "EXMPLGB2XXX", []string{"Test Account"}, "invalid", "GB28NWBK12345678901234", true, ""),
		Entry("should fail with invalid IBAN",
			"GB", "GBP", "123456", "GBDSC", "EXMPLGB2XXX", []string{"Test Account"}, "12345678", "invalid", true, ""),
	)
})
var _ = DescribeAccount("[Create With Retries]", func() {
	var client accounts.Client
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     10 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
	})
	Describe("Create with Retry", func() {
		var accountData *models.AccountData
		BeforeEach(func() {
			version := int64(0)
			country := "GB"
			accountData = &models.AccountData{
				Type:           "accounts",
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Attributes: &models.AccountAttributes{
					Country:                 &country,
					BaseCurrency:            "GBP",
					BankID:                  "123456",
					BankIDCode:              "GBDSC",
					Bic:                     "EXMPLGB2XXX",
					AccountNumber:           "12345678",
					AlternativeNames:        []string{"AltName1", "AltName2"},
					Iban:                    "GB11NWBK56000123456789",
					Name:                    []string{"John Smith"},
					SecondaryIdentification: "1234",
				},
				Version: &version,
			}
		})
		It("should create an account with retries", func() {
			// Simulate a random temporary server failure by adjusting the BaseURL
			if rand.Intn(2) == 0 {
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      "http://invalid-url",
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())
				// Wait for a random time between 0 and 1 seconds
				time.Sleep(time.Duration(rand.Float64()) * time.Second)

				// Reset the transport to the correct BaseURL
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      testBaseURL,
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())

				createdAccount, err := client.Create(accountData)
				Expect(err).NotTo(HaveOccurred())

				if createdAccount != nil && createdAccount.Version != nil {
					err := client.Delete(createdAccount.ID, *createdAccount.Version)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})

var _ = DescribeAccount("[Fetch]", func() {
	var client accounts.Client
	var accountData *models.AccountData
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     5 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
		orgId := uuid.New().String()
		Id := uuid.New().String()
		version := int64(0)
		country := "GB"
		accountData = &models.AccountData{
			Type:           "accounts",
			ID:             Id,
			OrganisationID: orgId,
			Attributes: &models.AccountAttributes{
				Country:                 &country,
				BaseCurrency:            "GBP",
				BankID:                  "123456",
				BankIDCode:              "GBDSC",
				Bic:                     "EXMPLGB2XXX",
				AccountNumber:           "12345678",
				AlternativeNames:        []string{"AltName1", "AltName2"},
				Iban:                    "GB11NWBK56000123456789",
				Name:                    []string{"John Smith"},
				SecondaryIdentification: "1234",
			},
			Version: &version,
		}
	})
	Describe("Fetch Account", func() {
		var createdAccount *models.AccountData
		BeforeEach(func() {
			var err error
			createdAccount, err = client.Create(accountData)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			if createdAccount != nil && createdAccount.Version != nil {
				err := client.Delete(createdAccount.ID, *createdAccount.Version)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should fetch an account", func() {
			fetchedAccount, err := client.Fetch(createdAccount.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedAccount).To(Equal(createdAccount))
		})
		It("should return an error when fetching with an invalid ID", func() {
			invalidID := "invalid-id"
			_, err := client.Fetch(invalidID)
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = DescribeAccount("[Fetch With Retries]", func() {
	var client accounts.Client
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     10 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
	})

	Describe("Fetch with Retry", func() {
		var accountData *models.AccountData
		BeforeEach(func() {
			version := int64(0)
			country := "GB"
			accountData = &models.AccountData{
				Type:           "accounts",
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Attributes: &models.AccountAttributes{
					Country:                 &country,
					BaseCurrency:            "GBP",
					BankID:                  "123456",
					BankIDCode:              "GBDSC",
					Bic:                     "EXMPLGB2XXX",
					AccountNumber:           "12345678",
					AlternativeNames:        []string{"AltName1", "AltName2"},
					Iban:                    "GB11NWBK56000123456789",
					Name:                    []string{"John Smith"},
					SecondaryIdentification: "1234",
				},
				Version: &version,
			}
		})

		It("should fetch an account with retries", func() {
			// Simulate a random temporary server failure by adjusting the BaseURL
			if rand.Intn(2) == 0 {
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      "http://invalid-url",
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())
				// Wait for a random time between 0 and 1 seconds
				time.Sleep(time.Duration(rand.Float64()) * time.Second)

				// Reset the transport to the correct BaseURL
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      testBaseURL,
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())

				createdAccount, err := client.Create(accountData)
				Expect(err).NotTo(HaveOccurred())

				fetchedAccount, err := client.Fetch(createdAccount.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(fetchedAccount.ID).To(Equal(createdAccount.ID))

				if createdAccount != nil && createdAccount.Version != nil {
					err := client.Delete(createdAccount.ID, *createdAccount.Version)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})

var _ = DescribeAccount("[Delete]", func() {
	var client accounts.Client
	var accountData *models.AccountData
	var version int64
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     5 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
		orgId := uuid.New().String()
		Id := uuid.New().String()
		version = int64(0)
		country := "GB"
		accountData = &models.AccountData{
			Type:           "accounts",
			ID:             Id,
			OrganisationID: orgId,
			Attributes: &models.AccountAttributes{
				Country:                 &country,
				BaseCurrency:            "GBP",
				BankID:                  "123456",
				BankIDCode:              "GBDSC",
				Bic:                     "EXMPLGB2XXX",
				AccountNumber:           "12345678",
				AlternativeNames:        []string{"AltName1", "AltName2"},
				Iban:                    "GB11NWBK56000123456789",
				Name:                    []string{"John Smith"},
				SecondaryIdentification: "1234",
			},
			Version: &version,
		}
	})
	Describe("Delete Account", func() {
		var createdAccount *models.AccountData
		BeforeEach(func() {
			var err error
			createdAccount, err = client.Create(accountData)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete an account", func() {
			err := client.Delete(createdAccount.ID, version)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return an error when deleting with an invalid ID", func() {
			invalidID := "invalid-id"
			err := client.Delete(invalidID, version)
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = DescribeAccount("[Delete With Retries]", func() {
	var client accounts.Client
	BeforeEach(func() {
		client = accounts.New(accounts.Options{
			BaseURL:      testBaseURL,
			Duration:     10 * time.Second,
			Retries:      3,
			InitialDelay: 1 * time.Second,
			Multiplier:   2,
			Factor:       1.5,
			LogLevel:     0,
		})
	})

	Describe("Delete with Retry", func() {
		var accountData *models.AccountData
		var version int64
		BeforeEach(func() {
			version = int64(0)
			country := "GB"
			accountData = &models.AccountData{
				Type:           "accounts",
				ID:             uuid.New().String(),
				OrganisationID: uuid.New().String(),
				Attributes: &models.AccountAttributes{
					Country:                 &country,
					BaseCurrency:            "GBP",
					BankID:                  "123456",
					BankIDCode:              "GBDSC",
					Bic:                     "EXMPLGB2XXX",
					AccountNumber:           "12345678",
					AlternativeNames:        []string{"AltName1", "AltName2"},
					Iban:                    "GB11NWBK56000123456789",
					Name:                    []string{"John Smith"},
					SecondaryIdentification: "1234",
				},
				Version: &version,
			}
		})

		It("should delete an account with retries", func() {
			// Simulate a random temporary server failure by adjusting the BaseURL
			if rand.Intn(2) == 0 {
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      "http://invalid-url",
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())
				// Wait for a random time between 0 and 1 seconds
				time.Sleep(time.Duration(rand.Float64()) * time.Second)

				// Reset the transport to the correct BaseURL
				client.SetTransport(accounts.New(accounts.Options{
					BaseURL:      testBaseURL,
					Duration:     5 * time.Second,
					Retries:      3,
					InitialDelay: 1 * time.Second,
					Multiplier:   2,
					Factor:       1.5,
				}).GetTransport())

				createdAccount, err := client.Create(accountData)
				Expect(err).NotTo(HaveOccurred())

				err = client.Delete(createdAccount.ID, version)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
