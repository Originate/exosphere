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

	var _ = Describe("GetAllBuiltDependencies", func() {

		It("should include the dependencies of all services and of the app itself", func() {
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).ToNot(HaveOccurred())
			builtDependencies := config.GetAllBuiltDependencies(appConfig, serviceConfigs, appDir, homeDir)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})

	var _ = Describe("GetAppBuiltDependencies", func() {

		It("should include the dependencies of the application", func() {
			builtDependencies := config.GetAppBuiltDependencies(appConfig, appDir, homeDir)
			dependencyNames := []string{"mongo", "exocom"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})

	var _ = Describe("GetEnvironmentVariables", func() {
		It("should return the environment variables of all dependencies", func() {
			appDir := path.Join("..", "..", "example-apps", "complex-setup-app")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			builtDependencies := config.GetAppBuiltDependencies(appConfig, appDir, homeDir)
			actual := config.GetEnvironmentVariables(builtDependencies)
			expected := map[string]string{"EXOCOM_PORT": "80", "DB_NAME": "test-db"}
			Expect(actual).To(Equal(expected))
		})
	})

})
