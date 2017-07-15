package serviceConfigHelpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServiceConfigHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ParseDockerComposeLog Suite")
}
