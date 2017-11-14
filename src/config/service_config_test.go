package config_test

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Service Config Helpers", func() {

	var _ = Describe("GetServiceConfigs", func() {
		var serviceConfigs map[string]types.ServiceConfig
		var appConfig types.AppConfig
		var appDir string

		var _ = BeforeEach(func() {
			appDir = helpers.GetTestApplicationDir("complex-setup-app")
			var err error
			appConfig, err = types.NewAppConfig(appDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not return an error when all service.yml files are valid", func() {
			var err error
			serviceConfigs, err = config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include all services", func() {
			for _, serviceRole := range appConfig.GetSortedServiceRoles() {
				_, exists := serviceConfigs[serviceRole]
				Expect(exists).To(Equal(true))
			}
		})

		It("should contain correct configuration for the internal service", func() {
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:        "public",
				Description: "dummy html service used for testing setup only - does not run",
				Author:      "test-author",
				ServiceMessages: types.ServiceMessages{
					Sends:    []string{"todo.create"},
					Receives: []string{"todo.created"},
				},
				Development: types.ServiceDevelopmentConfig{
					Scripts: map[string]string{
						"run": `echo "does not run"`,
					},
					Port: "80",
				},
			})
			Expect(err).ToNot(HaveOccurred())
			actual, err := yaml.Marshal(serviceConfigs["html-server"])
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

		It("should contain correct configuration for the external docker image", func() {
			serviceMessages := types.ServiceMessages{
				Sends:    []string{"users.list", "users.create"},
				Receives: []string{"users.listed", "users.created"},
			}
			development := types.ServiceDevelopmentConfig{
				Port: "5000",
				Scripts: map[string]string{
					"run": "node server.js",
				},
			}
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:            "public",
				Description:     "says hello to the world, ignores .txt files when file watching",
				Author:          "exospheredev",
				ServiceMessages: serviceMessages,
				Development:     development,
			})
			Expect(err).ToNot(HaveOccurred())
			actual, err := yaml.Marshal(serviceConfigs["external-service"])
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

	})

	var _ = Describe("GetServiceBuiltDependencies", func() {
		var appConfig types.AppConfig
		var appDir string

		var _ = BeforeEach(func() {
			appDir = helpers.GetTestApplicationDir("rds")
			var err error
			appConfig, err = types.NewAppConfig(appDir)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should include both service and application dependencies", func() {
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).ToNot(HaveOccurred())
			builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfigs["my-sql-service"], appConfig, appDir)
			_, exists := builtDependencies["mysql"]
			Expect(exists).To(Equal(true))
		})

	})

})
