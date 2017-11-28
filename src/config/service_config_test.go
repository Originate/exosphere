package config_test

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Config Helpers", func() {

	var _ = Describe("GetServiceBuiltDependencies", func() {
		var appContext types.AppContext

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("rds")
			var err error
			appContext, err = types.GetAppContext(appDir)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should include both service and application dependencies", func() {
			serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
			Expect(err).ToNot(HaveOccurred())
			builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfigs["my-sql-service"], appContext)
			_, exists := builtDependencies["mysql"]
			Expect(exists).To(Equal(true))
		})

	})

})
