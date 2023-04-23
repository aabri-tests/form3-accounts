package form3

import "github.com/aabri-assignments/form3-accounts/v1/accounts"

// Form3 is a struct that represents a client for the Form3 Accounts API.
type Form3 struct {
	options accounts.Options
}

// NewForm3 creates a new Form3 client with the specified base URL.
func NewForm3(options accounts.Options) (*Form3, error) {
	return &Form3{
		options: options,
	}, nil
}

func (f *Form3) Accounts() accounts.Client {
	return accounts.New(f.options)
}
