package appSetup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppConfigHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppSetup Suite")
}