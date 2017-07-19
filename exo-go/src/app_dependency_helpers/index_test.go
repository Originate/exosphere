package appDependencyHelpers_test

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppDependency", func() {
	var appConfig types.AppConfig
	var appDir string
	var exocom appDependencyHelpers.AppDependency
	var nats appDependencyHelpers.AppDependency
	var mongo appDependencyHelpers.AppDependency

	var _ = BeforeSuite(func() {
		appDir = path.Join("..", "..", "..", "exosphere-shared", "example-apps", "complex-setup-app")
		var err error
		appConfig, err = appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("Build", func() {
		It("should build each dependency correctly", func() {
			for _, dependency := range appConfig.Dependencies {
				builtDependency := appDependencyHelpers.Build(dependency, appConfig, appDir)
				switch dependency.Name {
				case "exocom":
					exocom = builtDependency
				case "mongo":
					mongo = builtDependency
				}
			}
			nats = appDependencyHelpers.Build(types.Dependency{
				Name:    "nats",
				Version: "0.9.6",
			}, appConfig, appDir)
		})
	})

	var _ = Describe("GetContainerName", func() {
		It("should be the concatenation of dependency name and version", func() {
			Expect(exocom.GetContainerName()).To(Equal("exocom0.22.1"))
			Expect(mongo.GetContainerName()).To(Equal("mongo3.4.0"))
			Expect(nats.GetContainerName()).To(Equal("nats0.9.6"))
		})
	})

	var _ = Describe("GetDokerConfig", func() {
		It("should return the correct docker config for exocom", func() {
			actual, err := exocom.GetDockerConfig()
			Expect(err).NotTo(HaveOccurred())
			expectedServiceRoutes := []string{
				`[{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]},{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]},{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]},{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}]`,
				`[{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]},{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]},{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]},{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]}]`,
				`[{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]},{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]},{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]},{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]}]`,
				`[{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]},{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]},{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]},{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}]`,
				`[{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]},{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]},{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]},{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}]`,
				`[{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]},{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]},{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]},{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}]`,
			}
			matched := false
			for _, serviceRoutes := range expectedServiceRoutes {
				expected := types.DockerConfig{
					Image:         "originate/exocom:0.22.1",
					Command:       "bin/exocom",
					ContainerName: "exocom0.22.1",
					Environment: map[string]string{
						"ROLE":           "exocom",
						"PORT":           "$EXOCOM_PORT",
						"SERVICE_ROUTES": serviceRoutes,
					},
				}
				if reflect.DeepEqual(actual, expected) {
					matched = true
				}
			}
			Expect(matched).To(Equal(true))
		})

		It("should return the correct docker config for nats", func() {
			actual, err := nats.GetDockerConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(types.DockerConfig{
				Image:         "nats:0.9.6",
				ContainerName: "nats0.9.6",
			}).To(Equal(actual))
		})

		It("should return the correct docker config for generic dependency", func() {
			actual, err := mongo.GetDockerConfig()
			Expect(err).NotTo(HaveOccurred())
			volumesRegex := regexp.MustCompile(`./\.exosphere/complex-setup-app/mongo/data:/data/db`)
			Expect(volumesRegex.MatchString(actual.Volumes[0])).To(Equal(true))
			actual.Volumes = nil
			fmt.Println(actual)
			Expect(types.DockerConfig{
				Image:         "mongo:3.4.0",
				ContainerName: "mongo3.4.0",
				Ports:         []string{"4000:4000"},
			}).To(Equal(actual))
		})

	})

	var _ = Describe("GetEnviromentVariables", func() {
		It("should set a default port for EXOCOM_PORT if it is not set", func() {
			expected := map[string]string{"EXOCOM_PORT": "80"}
			Expect(exocom.GetEnvVariables()).To(Equal(expected))
		})

		It("should return the EXOCOM_PORT as is set on the user's machine", func() {
			if err := os.Setenv("EXOCOM_PORT", "5000"); err != nil {
				panic(err)
			}
			expected := map[string]string{"EXOCOM_PORT": "5000"}
			Expect(exocom.GetEnvVariables()).To(Equal(expected))
		})

		It("should include the correct NATS_HOST for nats dependency", func() {
			expected := map[string]string{"NATS_HOST": "nats0.9.6"}
			Expect(nats.GetEnvVariables()).To(Equal(expected))
		})

		It("should return the correct environment variables for external dependencies", func() {
			expected := map[string]string{"DB_NAME": "test-db"}
			Expect(mongo.GetEnvVariables()).To(Equal(expected))
		})
	})

	var _ = Describe("GetOnlineText", func() {
		It("should return the correct online text for exocom", func() {
			expected := "ExoCom WebSocket listener online"
			Expect(exocom.GetOnlineText()).To(Equal(expected))
		})

		It("should return the correct online text for nats", func() {
			expected := "Listening for route connections"
			Expect(nats.GetOnlineText()).To(Equal(expected))
		})

		It("should include the correct online text for external dependencies", func() {
			expected := "waiting for connections"
			Expect(mongo.GetOnlineText()).To(Equal(expected))
		})
	})
})
