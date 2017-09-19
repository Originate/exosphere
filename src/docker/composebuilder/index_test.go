package composebuilder_test

import (
	"path"
	"regexp"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ComposeBuilder", func() {
	var _ = Describe("GetServiceDockerConfigs", func() {
		var _ = Describe("unshared docker configs", func() {
			var dockerConfigs types.DockerConfigs
			var appDir string

			var _ = BeforeEach(func() {
				appDir = path.Join(cwd, "..", "..", "..", "example-apps", "external-dependency")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "mongo"
				dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir, false)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the docker config for the service itself", func() {
				dockerConfig, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(true))
				Expect(dockerConfig.DependsOn).To(ConsistOf([]string{"exocom0.26.1", "mongo3.4.0"}))
				dockerConfig.DependsOn = nil
				Expect(dockerConfig).To(Equal(types.DockerConfig{
					Build: map[string]string{
						"dockerfile": "Dockerfile.dev",
						"context":    path.Join(appDir, "mongo"),
					},
					ContainerName: "mongo",
					Command:       "node server.js",
					Links:         []string{"mongo3.4.0:mongo"},
					Volumes:       []string{path.Join(appDir, "mongo") + ":/mnt"},
					Environment: map[string]string{
						"ROLE":        "mongo",
						"EXOCOM_HOST": "exocom0.26.1",
						"MONGO":       "mongo",
					},
				}))
			})

			It("should include the docker configs for the service's dependencies", func() {
				dockerConfig, exists := dockerConfigs["mongo3.4.0"]
				Expect(exists).To(Equal(true))
				volumesRegex := regexp.MustCompile(`./\.exosphere/Exosphere-application-with-a-third-party-dependency/mongo/data:/data/db`)
				Expect(volumesRegex.MatchString(dockerConfig.Volumes[0])).To(Equal(true))
				dockerConfig.Volumes = nil
				Expect(dockerConfig).To(Equal(types.DockerConfig{
					Image:         "mongo:3.4.0",
					ContainerName: "mongo3.4.0",
					Ports:         []string{"27017:27017"},
				}))
			})
		})

		var _ = Describe("shared docker configs", func() {
			var dockerConfigs types.DockerConfigs

			var _ = BeforeEach(func() {
				appDir := path.Join(cwd, "..", "..", "..", "example-apps", "complex-setup-app")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "users-service"
				dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir, false)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not include the docker config that is defined in application.yml", func() {
				_, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(false))
			})
		})

		var _ = Describe("service specific dependency", func() {
			var dockerConfigs types.DockerConfigs

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "example-apps", "service-specific-dependency")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "postgres-service"
				dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir, false)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should pass dependency env variables to services", func() {
				postgresServiceDockerConfig, exists := dockerConfigs["postgres-service"]
				Expect(exists).To(Equal(true))
				Expect(postgresServiceDockerConfig.Environment["DB_NAME"]).To(Equal("my_db"))
			})
		})
	})

	var _ = Describe("building for production", func() {
		var dockerConfigs types.DockerConfigs
		var appDir string

		var _ = BeforeEach(func() {
			appDir = path.Join("..", "..", "..", "example-apps", "simple")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())
			serviceData := appConfig.GetServiceData()
			serviceName := "web"
			dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir, true)
			dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets the correct dockerfile", func() {
			dockerConfig, exists := dockerConfigs["web"]
			Expect(exists).To(Equal(true))
			Expect(dockerConfig).To(Equal(types.DockerConfig{
				Build: map[string]string{
					"dockerfile": "Dockerfile.prod",
					"context":    path.Join(appDir, "web"),
				},
				ContainerName: "web",
				Command:       "",
				Links:         []string{},
				Volumes:       []string{},
				Environment: map[string]string{
					"ROLE":        "web",
					"EXOCOM_HOST": "exocom0.26.1",
				},
				DependsOn: []string{"exocom0.26.1"},
			}))
		})
	})

	var _ = Describe("compiles the docker compose project name properly", func() {
		expected := "spacetweet123"

		It("converts all characters to lowercase", func() {
			actual := composebuilder.GetDockerComposeProjectName("SpaceTweet123")
			Expect(actual).To(Equal(expected))
		})

		It("strips non-alphanumeric characters", func() {
			actual := composebuilder.GetDockerComposeProjectName("$Space-Tweet_123")
			Expect(actual).To(Equal(expected))
		})

		It("strips whitespace characters", func() {
			actual := composebuilder.GetDockerComposeProjectName("Space   Tweet  123")
			Expect(actual).To(Equal(expected))
		})
	})
})
