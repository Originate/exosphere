package messageTranslator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMessageTranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Message Translator Suite")
}
