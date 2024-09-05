//go:build e2e

package create_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModuleCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "'Create' Command Test Suite")
}
