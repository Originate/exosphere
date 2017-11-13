package composebuilder_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAppSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "docker/composebuilder Suite")
}
