package appConfigHelpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppConfigHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppConfigHelpers Suite")
}
