package composebuilder_test

import (
	"os"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/endpoints"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ComposeBuilder", func() {

	var _ = Describe("building external dependencies", func() {
		var dockerCompose *types.DockerCompose
		var serviceEndpoints map[string]*endpoints.ServiceEndpoint

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("external-dependency")
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceRole := "mongo"
			buildMode := types.BuildMode{
				Type:        types.BuildModeTypeLocal,
				Mount:       true,
				Environment: types.BuildModeEnvironmentDevelopment,
			}
			serviceEndpoints = map[string]*endpoints.ServiceEndpoint{
				"mongo": &endpoints.ServiceEndpoint{},
			}
			dockerCompose, err = composebuilder.GetServiceDockerCompose(appContext, serviceRole, buildMode, serviceEndpoints)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include the docker config for the service itself", func() {
			dockerConfig, exists := dockerCompose.Services["mongo"]
			Expect(exists).To(Equal(true))
			Expect(dockerConfig.DependsOn).To(Equal([]string{"exocom0.26.1", "mongo3.4.0"}))
			dockerConfig.DependsOn = nil
			Expect(dockerConfig).To(Equal(types.DockerConfig{
				Build: map[string]string{
					"dockerfile": "Dockerfile.dev",
					"context":    "${APP_PATH}/mongo",
				},
				ContainerName: "mongo",
				Command:       "node server.js",
				Ports:         []string{},
				Volumes:       []string{"${APP_PATH}/mongo:/mnt"},
				Environment: map[string]string{
					"ROLE":        "mongo",
					"EXOCOM_HOST": "exocom0.26.1",
					"MONGO":       "mongo3.4.0",
				},
				Restart: "on-failure",
			}))
		})

		It("should include the docker configs for the service's dependencies", func() {
			dockerConfig, exists := dockerCompose.Services["mongo3.4.0"]
			Expect(exists).To(Equal(true))
			Expect(dockerConfig).To(Equal(types.DockerConfig{
				Image:         "mongo:3.4.0",
				ContainerName: "mongo3.4.0",
				Ports:         []string{"27017:27017"},
				Volumes:       []string{"mongo__data_db:/data/db"},
				Restart:       "on-failure",
			}))
		})
	})

	var _ = Describe("compiling environment variables", func() {
		var serviceEndpoints map[string]*endpoints.ServiceEndpoint
		var appContext *context.AppContext
		var serviceRole string

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("complex-setup-app")
			err := os.Setenv("EXOSPHERE_SECRET", "exosphere-value")
			Expect(err).NotTo(HaveOccurred())
			appContext, err = context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceRole = "users-service"
			serviceEndpoints = map[string]*endpoints.ServiceEndpoint{
				"html-server":      &endpoints.ServiceEndpoint{},
				"api-service":      &endpoints.ServiceEndpoint{},
				"external-service": &endpoints.ServiceEndpoint{},
				"users-service":    &endpoints.ServiceEndpoint{},
				"todo-service":     &endpoints.ServiceEndpoint{},
			}
		})

		var _ = AfterEach(func() {
			err := os.Unsetenv("EXOSPHERE_SECRET")
			Expect(err).NotTo(HaveOccurred())
		})

		It("compiles development variables", func() {
			buildMode := types.BuildMode{
				Type:        types.BuildModeTypeLocal,
				Environment: types.BuildModeEnvironmentDevelopment,
			}
			dockerCompose, err := composebuilder.GetServiceDockerCompose(appContext, serviceRole, buildMode, serviceEndpoints)
			Expect(err).NotTo(HaveOccurred())
			expectedVars := map[string]string{
				"ENV1":             "value1",
				"ENV2":             "value2",
				"ENV3":             "dev_value3",
				"EXOSPHERE_SECRET": "exosphere-value",
			}
			actualVars := dockerCompose.Services["users-service"].Environment
			for k, v := range expectedVars {
				Expect(actualVars).Should(HaveKeyWithValue(k, v))
			}
		})
	})

	var _ = Describe("service specific dependency", func() {
		var dockerCompose *types.DockerCompose
		var serviceEndpoints map[string]*endpoints.ServiceEndpoint

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("service-specific-dependency")
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceRole := "postgres-service"
			buildMode := types.BuildMode{
				Type:        types.BuildModeTypeLocal,
				Environment: types.BuildModeEnvironmentDevelopment,
			}
			serviceEndpoints = map[string]*endpoints.ServiceEndpoint{
				"postgres-service": &endpoints.ServiceEndpoint{},
			}
			dockerCompose, err = composebuilder.GetServiceDockerCompose(appContext, serviceRole, buildMode, serviceEndpoints)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should pass dependency env variables to services", func() {
			postgresServiceDockerConfig, exists := dockerCompose.Services["postgres-service"]
			Expect(exists).To(Equal(true))
			Expect(postgresServiceDockerConfig.Environment["DB_NAME"]).To(Equal("my_db"))
		})
	})
})

var _ = Describe("building for local production", func() {
	var dockerCompose *types.DockerCompose
	var appDir string
	var serviceEndpoints map[string]*endpoints.ServiceEndpoint

	var _ = BeforeEach(func() {
		appDir = helpers.GetTestApplicationDir("simple")
		appContext, err := context.GetAppContext(appDir)
		Expect(err).NotTo(HaveOccurred())
		serviceRole := "web"
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentProduction,
		}
		serviceEndpoints = map[string]*endpoints.ServiceEndpoint{
			"web": &endpoints.ServiceEndpoint{},
		}
		dockerCompose, err = composebuilder.GetServiceDockerCompose(appContext, serviceRole, buildMode, serviceEndpoints)
		Expect(err).NotTo(HaveOccurred())
	})

	It("sets the correct dockerfile", func() {
		dockerConfig, exists := dockerCompose.Services["web"]
		Expect(exists).To(Equal(true))
		Expect(dockerConfig).To(Equal(types.DockerConfig{
			Build: map[string]string{
				"dockerfile": "Dockerfile.prod",
				"context":    "${APP_PATH}/web",
			},
			ContainerName: "web",
			Command:       "",
			Ports:         []string{},
			Volumes:       []string{"${APP_PATH}/web:/mnt"},
			Environment: map[string]string{
				"ROLE":        "web",
				"EXOCOM_HOST": "exocom0.26.1",
			},
			DependsOn: []string{"exocom0.26.1"},
			Restart:   "on-failure",
		}))
	})
})
