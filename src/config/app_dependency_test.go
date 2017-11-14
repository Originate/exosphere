package config_test

import (
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppDevelopmentDependency", func() {
	var appConfig types.AppConfig
	var appDir string

	var _ = BeforeEach(func() {
		appDir = helpers.GetTestApplicationDir("complex-setup-app")
		var err error
		appConfig, err = types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("Build", func() {
		It("should build each dependency successfully", func() {
			for _, dependency := range appConfig.Development.Dependencies {
				_ = config.NewAppDevelopmentDependency(dependency, appConfig, appDir)
			}
		})
	})

	var _ = Describe("exocom dev dependency", func() {
		var exocomDev config.AppDevelopmentDependency

		var _ = BeforeEach(func() {
			for _, dependency := range appConfig.Development.Dependencies {
				if dependency.Name == "exocom" {
					exocomDev = config.NewAppDevelopmentDependency(dependency, appConfig, appDir)
					break
				}
			}
		})

		var _ = Describe("GetContainerName", func() {
			It("should be the concatenation of dependency name and version", func() {
				Expect(exocomDev.GetContainerName()).To(Equal("exocom0.26.1"))
			})
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for exocom", func() {
				actual, err := exocomDev.GetDockerConfig()
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
					Restart: "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for exocom", func() {
				expected := map[string]string{
					"EXOCOM_HOST": "exocom0.26.1",
				}
				Expect(exocomDev.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("exocom prod dependency", func() {
		var exocomProd config.AppProductionDependency
		var _ = BeforeEach(func() {
			for _, dependency := range appConfig.Production.Dependencies {
				if dependency.Name == "exocom" {
					exocomProd = config.NewAppProductionDependency(dependency, appConfig, appDir)
					break
				}
			}
		})

		var _ = Describe("GetDeploymentServiceEnvVariables", func() {
			It("should return the EXOCOM_HOST", func() {
				Expect(exocomProd.GetDeploymentServiceEnvVariables(types.Secrets{})).To(Equal(map[string]string{
					"EXOCOM_HOST": "exocom.complex-setup-app.local",
				}))
			})
		})

		var _ = Describe("GetDeploymentConfig", func() {
			It("should return the correct deployment config for exocom", func() {
				actual, err := exocomProd.GetDeploymentConfig()
				Expect(err).NotTo(HaveOccurred())
				expectedServiceRoutes := []string{
					`{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]}`,
					`{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]}`,
					`{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}`,
					`{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}`,
				}
				for _, serviceRoute := range expectedServiceRoutes {
					Expect(strings.Contains(actual["serviceRoutes"], serviceRoute))
				}
				Expect(actual["version"]).To(Equal("0.26.1"))
				Expect(actual["dnsName"]).To(Equal("originate.com"))
			})
		})
	})

	var _ = Describe("generic dependency", func() {
		var mongo config.AppDevelopmentDependency

		var _ = BeforeEach(func() {
			for _, dependency := range appConfig.Development.Dependencies {
				if dependency.Name == "mongo" {
					mongo = config.NewAppDevelopmentDependency(dependency, appConfig, appDir)
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
				Expect(types.DockerConfig{
					Image:         "mongo:3.4.0",
					ContainerName: "mongo3.4.0",
					Ports:         []string{"4000:4000"},
					Volumes: []string{
						"${APP_PATH}/.exosphere/data/mongo:/data/db",
					},
					Environment: map[string]string{"DB_NAME": "test-db"},
					Restart:     "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for generic dependencies", func() {
				expected := map[string]string{
					"COLLECTION_NAME": "test-collection",
					"MONGO":           "mongo3.4.0",
				}
				Expect(mongo.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("nats dependency", func() {
		var nats config.AppDevelopmentDependency

		var _ = BeforeEach(func() {
			nats = config.NewAppDevelopmentDependency(types.DevelopmentDependencyConfig{
				Name:    "nats",
				Version: "0.9.6",
			}, appConfig, appDir)
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
					Restart:       "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should include the correct service environment variables for nats", func() {
				expected := map[string]string{"NATS_HOST": "nats0.9.6"}
				Expect(nats.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("rds dependency", func() {
		var rds config.AppProductionDependency
		var _ = BeforeEach(func() {
			appDir = helpers.GetTestApplicationDir("rds")
			var err error
			appConfig, err = types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			for _, dependency := range appConfig.Production.Dependencies {
				if dependency.Name == "postgres" {
					rds = config.NewAppProductionDependency(dependency, appConfig, appDir)
					break
				}
			}
		})

		var _ = Describe("GetDeploymentServiceEnvVariables", func() {
			It("should return the required service env vars", func() {
				secrets := types.Secrets{
					"POSTGRES_PASSWORD": "password123",
				}
				Expect(rds.GetDeploymentServiceEnvVariables(secrets)).To(Equal(map[string]string{
					"POSTGRES":          "my-db.rds.local",
					"DATABASE_NAME":     "my-db",
					"DATABASE_USERNAME": "originate-user",
					"DATABASE_PASSWORD": "password123",
				}))
			})
		})
	})

})
