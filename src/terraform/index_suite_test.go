package terraform_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTerraformFileHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terraform Suite")
}

var appDir string

var _ = BeforeSuite(func() {
	var err error
	appDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
})
