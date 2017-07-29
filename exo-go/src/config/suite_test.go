package config_test

import (
	"testing"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
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
	homeDir, err = osHelpers.GetUserHomeDir()
	if err != nil {
		panic(err)
	}
})
