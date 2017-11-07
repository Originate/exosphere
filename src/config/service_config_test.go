package config_test

import (
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Service Config Helpers", func() {
	var appConfig types.AppConfig
	var appDir string

	var _ = BeforeEach(func() {
		appDir = path.Join("..", "..", "example-apps", "complex-setup-app")
		var err error
		appConfig, err = types.NewAppConfig(appDir)
		Expect(err).ToNot(HaveOccurred())
	})

	var _ = Describe("GetServiceData", func() {

		It("should join the public, private and worker services into a single map", func() {
			actual := appConfig.GetServiceData()
			Expect(map[string]types.ServiceData{
				"todo-service": types.ServiceData{
					Location: "./todo-service",
				},
				"users-service": types.ServiceData{
					Location: "./users-service",
					MessageTranslations: []types.MessageTranslation{
						types.MessageTranslation{
							Public:   "users create",
							Internal: "mongo create",
						},
					},
				},
				"external-service": types.ServiceData{
					DockerImage: "originate/test-web-server",
				},
				"html-server": types.ServiceData{
					Location: "./html-server",
				},
				"api-service": types.ServiceData{
					Location: "./api-service",
				},
			}).To(Equal(actual))
		})
	})

	var _ = Describe("GetServiceConfigs", func() {
		var serviceConfigs map[string]types.ServiceConfig

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
				Type:        "html-server",
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
			environmentVars := map[string]string{
				"EXTERNAL_SERVICE_HOST": "external-service0.1.2",
				"EXTERNAL_SERVICE_PORT": "$EXTERNAL_SERVICE_PORT",
			}
			docker := types.DockerConfig{
				Ports:       []string{"5000:5000"},
				Volumes:     []string{"{{EXO_DATA_PATH}}:/data/db"},
				Environment: environmentVars,
			}
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:            "external-service",
				Description:     "says hello to the world, ignores .txt files when file watching",
				Author:          "exospheredev",
				ServiceMessages: serviceMessages,
				Docker:          docker,
			})
			Expect(err).ToNot(HaveOccurred())
			actual, err := yaml.Marshal(serviceConfigs["external-service"])
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

	})

	var _ = Describe("GetInternalServiceConfigs", func() {
		var serviceConfigs map[string]types.ServiceConfig

		It("should not return an error when all service.yml files are valid", func() {
			var err error
			serviceConfigs, err = config.GetInternalServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include all internal services", func() {
			internalServiceNames := []string{"todo-service", "users-service", "html-server"}
			for _, serviceRole := range internalServiceNames {
				_, exists := serviceConfigs[serviceRole]
				Expect(exists).To(Equal(true))
			}
		})

		It("should contain correct configuration for each internal service", func() {
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:        "html-server",
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

	})

	var _ = Describe("GetServiceBuiltDependencies", func() {

		It("should include both service and application dependencies", func() {
			serviceConfigs, err := config.GetInternalServiceConfigs(appDir, appConfig)
			Expect(err).ToNot(HaveOccurred())
			builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfigs["todo-service"], appConfig, appDir, homeDir)
			dependencyNames := []string{"mongo"}
			for _, dependencyName := range dependencyNames {
				_, exists := builtDependencies[dependencyName]
				Expect(exists).To(Equal(true))
			}
		})

	})

})
