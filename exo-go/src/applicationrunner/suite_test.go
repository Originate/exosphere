package applicationrunner_test

import (
	"testing"

	"github.com/Originate/exosphere/exo-go/src/osplus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApplicationRunner Suite")
}

var homeDir string

var _ = BeforeSuite(func() {
	var err error
	homeDir, err = osplus.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})
