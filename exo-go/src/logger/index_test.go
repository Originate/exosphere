package logger_test

import (
	"github.com/Originate/exosphere/exo-go/src/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	It("work as expected", func() {
		logger := logger.NewLogger([]string{"role1"}, []string{"silencedRole1"})
		Expect(len(logger.Roles)).To(Equal(1))
	})
})
