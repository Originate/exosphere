package config_test

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Config Helpers", func() {

	var _ = Describe("GetServiceBuiltDependencies", func() {
		var appContext *context.AppContext

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("rds")
			var err error
			appContext, err = context.GetAppContext(appDir)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should include both service and application dependencies", func() {
			builtDependencies := config.GetBuiltServiceDevelopmentDependencies(appContext.ServiceContexts["my-sql-service"].Config, appContext)
			_, exists := builtDependencies["mysql"]
			Expect(exists).To(Equal(true))
		})

	})

})
