package config_test

import (
	"path"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppDependency", func() {
	var appConfig types.AppConfig
	var appDir string

	var _ = BeforeEach(func() {
		appDir = path.Join("..", "..", "example-apps", "complex-setup-app")
		var err error
		appConfig, err = types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("Build", func() {
		It("should build each dependency successfully", func() {
			for _, dependency := range appConfig.Dependencies {
				_ = config.NewAppDependency(dependency, appConfig, appDir, homeDir)
			}
		})
	})

	var _ = Describe("exocom dependency", func() {
		var exocom config.AppDependency

		var _ = BeforeEach(func() {
			for _, dependency := range appConfig.Dependencies {
				if dependency.Name == "exocom" {
					exocom = config.NewAppDependency(dependency, appConfig, appDir, homeDir)
					break
				}
			}
		})

		var _ = Describe("GetContainerName", func() {
			It("should be the concatenation of dependency name and version", func() {
				Expect(exocom.GetContainerName()).To(Equal("exocom0.26.1"))
			})
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for exocom", func() {
				actual, err := exocom.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				expectedServiceRoutes := []string{
					`{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]}`,
					`{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]}`,
					`{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}`,
					`{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}`,
				}
				for _, serviceRoute := range expectedServiceRoutes {
					Expect(strings.Contains(actual.Environment["SERVICE_ROUTES"], serviceRoute))
				}
				actual.Environment["SERVICE_ROUTES"] = ""
				Expect(types.DockerConfig{
					Image:         "originate/exocom:0.26.1",
					ContainerName: "exocom0.26.1",
					Environment: map[string]string{
						"ROLE":           "exocom",
						"SERVICE_ROUTES": "",
					},
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetOnlineText", func() {
			It("should return the correct online text for exocom", func() {
				expected := "ExoCom online at port"
				Expect(exocom.GetOnlineText()).To(Equal(expected))
			})
		})

		var _ = Describe("GetDeploymentServiceEnvVariables", func() {
			It("should return the EXOCOM_HOST", func() {
				Expect(exocom.GetDeploymentServiceEnvVariables()).To(Equal(map[string]string{
					"EXOCOM_HOST": "exocom.complex-setup-app.local",
				}))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for exocom", func() {
				expected := map[string]string{
					"EXOCOM_HOST": "exocom0.26.1",
				}
				Expect(exocom.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("generic dependency", func() {
		var mongo config.AppDependency

		var _ = BeforeEach(func() {
			for _, dependency := range appConfig.Dependencies {
				if dependency.Name == "mongo" {
					mongo = config.NewAppDependency(dependency, appConfig, appDir, homeDir)
					break
				}
			}
		})

		var _ = Describe("GetContainerName", func() {
			It("should be the concatenation of dependency name and version", func() {
				Expect(mongo.GetContainerName()).To(Equal("mongo3.4.0"))
			})
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for generic dependency", func() {
				actual, err := mongo.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				volumesRegex := regexp.MustCompile(`./\.exosphere/complex-setup-app/mongo/data:/data/db`)
				Expect(volumesRegex.MatchString(actual.Volumes[0])).To(Equal(true))
				actual.Volumes = nil
				Expect(types.DockerConfig{
					Image:         "mongo:3.4.0",
					ContainerName: "mongo3.4.0",
					Ports:         []string{"4000:4000"},
					Environment:   map[string]string{"DB_NAME": "test-db"},
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetOnlineText", func() {
			It("should include the correct online text for generic dependencies", func() {
				expected := "waiting for connections"
				Expect(mongo.GetOnlineText()).To(Equal(expected))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for generic dependencies", func() {
				expected := map[string]string{"COLLECTION_NAME": "test-collection"}
				Expect(mongo.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("nats dependency", func() {
		var nats config.AppDependency

		var _ = BeforeEach(func() {
			nats = config.NewAppDependency(types.DependencyConfig{
				Name:    "nats",
				Version: "0.9.6",
			}, appConfig, appDir, homeDir)
		})

		var _ = Describe("GetContainerName", func() {
			It("should be the concatenation of dependency name and version", func() {
				Expect(nats.GetContainerName()).To(Equal("nats0.9.6"))
			})
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for nats", func() {
				actual, err := nats.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image:         "nats:0.9.6",
					ContainerName: "nats0.9.6",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetOnlineText", func() {
			It("should return the correct online text for nats", func() {
				expected := "Listening for route connections"
				Expect(nats.GetOnlineText()).To(Equal(expected))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should include the correct service environment variables for nats", func() {
				expected := map[string]string{"NATS_HOST": "nats0.9.6"}
				Expect(nats.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

})
