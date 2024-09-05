//go:build e2e

package create_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_Create(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "'Create' Command Test Suite")
}
