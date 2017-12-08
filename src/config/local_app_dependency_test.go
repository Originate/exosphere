package config_test

import (
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalAppDependency", func() {
	var appContext *context.AppContext

	var _ = BeforeEach(func() {
		appDir := helpers.GetTestApplicationDir("complex-setup-app")
		var err error
		appContext, err = context.GetAppContext(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("Build", func() {
		It("should build each dependency successfully", func() {
			for dependencyName, dependency := range appContext.Config.Local.Dependencies {
				_ = config.NewLocalAppDependency(dependencyName, dependency, appContext)
			}
		})
	})

	var _ = Describe("exocom dev dependency", func() {
		var exocomDev config.LocalAppDependency

		var _ = BeforeEach(func() {
			for dependencyName, dependency := range appContext.Config.Local.Dependencies {
				if dependencyName == "exocom" {
					exocomDev = config.NewLocalAppDependency(dependencyName, dependency, appContext)
					break
				}
			}
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
					Image: "originate/exocom:0.26.1",
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
					"EXOCOM_HOST": "exocom",
				}
				Expect(exocomDev.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("exocom prod dependency", func() {
		var exocomProd config.RemoteAppDependency
		var _ = BeforeEach(func() {
			for dependencyName, dependency := range appContext.Config.Remote.Dependencies {
				if dependencyName == "exocom" {
					exocomProd = config.NewRemoteAppDependency(dependencyName, dependency, appContext)
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
		var mongo config.LocalAppDependency

		var _ = BeforeEach(func() {
			for dependencyName, dependency := range appContext.Config.Local.Dependencies {
				if dependencyName == "mongo" {
					mongo = config.NewLocalAppDependency(dependencyName, dependency, appContext)
					break
				}
			}
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for generic dependency", func() {
				actual, err := mongo.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image:       "mongo:3.4.0",
					Ports:       []string{"4000:4000"},
					Volumes:     []string{"mongo__data_db:/data/db"},
					Environment: map[string]string{"DB_NAME": "test-db"},
					Restart:     "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for generic dependencies", func() {
				expected := map[string]string{
					"COLLECTION_NAME": "test-collection",
					"MONGO":           "mongo",
				}
				Expect(mongo.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("nats dependency", func() {
		var nats config.LocalAppDependency

		var _ = BeforeEach(func() {
			nats = config.NewLocalAppDependency("nats", types.LocalDependency{
				Type:  "nats",
				Image: "nats:0.9.6",
			}, appContext)
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for nats", func() {
				actual, err := nats.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image:   "nats:0.9.6",
					Restart: "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should include the correct service environment variables for nats", func() {
				expected := map[string]string{"NATS_HOST": "nats"}
				Expect(nats.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("rds dependency", func() {
		var rds config.RemoteAppDependency
		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("rds")
			var err error
			appContext, err = context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			for dependencyName, dependency := range appContext.Config.Remote.Dependencies {
				if dependencyName == "postgres" {
					rds = config.NewRemoteAppDependency(dependencyName, dependency, appContext)
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
