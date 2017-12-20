package localdependencies_test

import (
	"github.com/Originate/exosphere/src/docker/composebuilder/localdependencies"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Config Helpers", func() {

	var appContext *context.AppContext

	var _ = BeforeEach(func() {
		appDir := helpers.GetTestApplicationDir("complex-setup-app")
		var err error
		appContext, err = context.GetAppContext(appDir)
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("GetBuiltLocalAppDependencies", func() {

		It("should include the dependencies of the application", func() {
			builtDependencies := localdependencies.GetBuiltLocalAppDependencies(appContext)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})
	})

	var _ = Describe("GetServiceBuiltDependencies", func() {
		var appContext *context.AppContext

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("rds")
			var err error
			appContext, err = context.GetAppContext(appDir)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should include both service and application dependencies", func() {
			builtDependencies := localdependencies.GetBuiltLocalServiceDependencies(appContext.ServiceContexts["my-sql-service"].Config, appContext)
			_, exists := builtDependencies["mysql"]
			Expect(exists).To(Equal(true))
		})
	})
})
