package dockercomposebuilder_test

import (
	"github.com/Originate/exosphere/exo-go/src/ostools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDockercomposebuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockercomposebuilder Suite")
}

var homeDir string

var _ = BeforeSuite(func() {
	var err error
	homeDir, err = ostools.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})
