package clientRegistry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestClientRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClientRegistry Suite")
}
