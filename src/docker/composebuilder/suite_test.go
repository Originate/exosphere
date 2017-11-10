package composebuilder_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "docker/composebuilder Suite")
}

var cwd string

var _ = BeforeSuite(func() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
})
