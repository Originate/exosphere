package messageCache_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMessageCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MessageCache Suite")
}
