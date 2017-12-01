package config_test

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Config Helpers", func() {

	var appContext *types.AppContext

	var _ = BeforeEach(func() {
		appDir := helpers.GetTestApplicationDir("complex-setup-app")
		var err error
		appContext, err = types.GetAppContext(appDir)
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("GetBuiltLocalAppDependencies", func() {

		It("should include the dependencies of the application", func() {
			builtDependencies := config.GetBuiltLocalAppDependencies(appContext)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})
})
