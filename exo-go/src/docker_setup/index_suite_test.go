package dockerSetup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDockerSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DockerSetup Suite")
}
