package accounts

import "github.com/onsi/ginkgo/v2"

// DescribeAccount annotates the test with the label.
func DescribeAccount(text string, body func()) bool {
	return ginkgo.Describe("[accounts] "+text, body)
}
