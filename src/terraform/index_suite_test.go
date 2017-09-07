package terraform_test

import (
	"os"

	"github.com/Originate/exosphere/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTerraformFileHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terraform Suite")
}

var appDir string
var homeDir string

var _ = BeforeSuite(func() {
	var err error
	appDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	homeDir, err = util.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})
