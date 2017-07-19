package terraformFileHelpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTerraformFileHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerraformFileHelpers Suite")
}
