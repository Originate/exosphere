package config_test

import (
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Config Helpers", func() {

	var appConfig types.AppConfig
	var appDir string

	var _ = BeforeEach(func() {
		appDir = path.Join("..", "..", "example-apps", "complex-setup-app")
		var err error
		appConfig, err = types.NewAppConfig(appDir)
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("GetBuiltDevelopmentDependencies", func() {

		It("should include the dependencies of all services and of the app itself", func() {
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).ToNot(HaveOccurred())
			builtDependencies := config.GetBuiltDevelopmentDependencies(appConfig, serviceConfigs, appDir)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})

	var _ = Describe("GetBuiltAppDevelopmentDependencies", func() {

		It("should include the dependencies of the application", func() {
			builtDependencies := config.GetBuiltAppDevelopmentDependencies(appConfig, appDir)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})
})
