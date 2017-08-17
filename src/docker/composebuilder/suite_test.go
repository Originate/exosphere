package composebuilder_test

import (
	"testing"

	"github.com/Originate/exosphere/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "docker/composebuilder Suite")
}

var homeDir string

var _ = BeforeSuite(func() {
	var err error
	homeDir, err = util.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})
