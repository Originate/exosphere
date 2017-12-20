package localdependencies_test

import (
	"encoding/json"

	"github.com/Originate/exosphere/src/docker/composebuilder/localdependencies"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalDependency", func() {
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
				_ = localdependencies.NewLocalDependency(dependencyName, dependency, appContext)
			}
		})
	})

	var _ = Describe("exocom dev dependency", func() {
		var exocomDev *localdependencies.LocalDependency

		var _ = BeforeEach(func() {
			for dependencyName, dependency := range appContext.Config.Local.Dependencies {
				if dependencyName == "exocom" {
					exocomDev = localdependencies.NewLocalDependency(dependencyName, dependency, appContext)
					break
				}
			}
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for exocom", func() {
				actual, err := exocomDev.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				serviceData, err := json.Marshal(map[string]map[string]interface{}{
					"api-service":      {},
					"external-service": {},
					"html-server": {
						"receives": []interface{}{"todo.created"},
						"sends":    []interface{}{"todo.create"},
					},
					"todo-service": {
						"receives": []interface{}{"todo.create"},
						"sends":    []interface{}{"todo.created"},
					},
					"users-service": {
						"receives": []interface{}{"mongo.list", "mongo.create"},
						"sends":    []interface{}{"mongo.listed", "mongo.created"},
						"translations": []interface{}{
							map[string]interface{}{
								"internal": "mongo create",
								"public":   "users create",
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image: "originate/exocom:0.27.0",
					Environment: map[string]string{
						"SERVICE_DATA": string(serviceData),
					},
					Volumes: []string{},
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

	var _ = Describe("generic dependency", func() {
		var mongo *localdependencies.LocalDependency

		var _ = BeforeEach(func() {
			for dependencyName, dependency := range appContext.Config.Local.Dependencies {
				if dependencyName == "mongo" {
					mongo = localdependencies.NewLocalDependency(dependencyName, dependency, appContext)
					break
				}
			}
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for generic dependency", func() {
				actual, err := mongo.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				serviceData, err := json.Marshal(map[string]map[string]interface{}{
					"api-service":      {},
					"external-service": {},
					"html-server":      {},
					"todo-service":     {},
					"users-service":    {},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image:   "mongo:3.4.0",
					Volumes: []string{"mongo__data_db:/data/db"},
					Environment: map[string]string{
						"DB_NAME":      "test-db",
						"SERVICE_DATA": string(serviceData),
					},
					Restart: "on-failure",
				}).To(Equal(actual))
			})
		})

		var _ = Describe("GetServiceEnvVariables", func() {
			It("should return the correct service environment variables for generic dependencies", func() {
				expected := map[string]string{
					"COLLECTION_NAME": "test-collection",
					"MONGO_HOST":      "mongo",
				}
				Expect(mongo.GetServiceEnvVariables()).To(Equal(expected))
			})
		})
	})

	var _ = Describe("nats dependency", func() {
		var nats *localdependencies.LocalDependency

		var _ = BeforeEach(func() {
			nats = localdependencies.NewLocalDependency("nats", types.LocalDependency{
				Image: "nats:0.9.6",
			}, appContext)
		})

		var _ = Describe("GetDockerConfig", func() {
			It("should return the correct docker config for nats", func() {
				actual, err := nats.GetDockerConfig()
				Expect(err).NotTo(HaveOccurred())
				serviceData, err := json.Marshal(map[string]map[string]interface{}{
					"api-service":      {},
					"external-service": {},
					"html-server":      {},
					"todo-service":     {},
					"users-service":    {},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(types.DockerConfig{
					Image:   "nats:0.9.6",
					Restart: "on-failure",
					Volumes: []string{},
					Environment: map[string]string{
						"SERVICE_DATA": string(serviceData),
					},
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
})
