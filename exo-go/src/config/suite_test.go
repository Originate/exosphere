package config_test

import (
	"testing"

	"github.com/Originate/exosphere/exo-go/src/ostools"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppConfigHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var homeDir string

var _ = BeforeSuite(func() {
	var err error
	homeDir, err = ostools.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
})
