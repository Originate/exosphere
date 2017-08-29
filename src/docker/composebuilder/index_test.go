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

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "example-apps", "external-dependency")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "mongo"
				dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the docker config for the service itself", func() {
				dockerConfig, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(true))
				Expect(dockerConfig).To(Equal(types.DockerConfig{
					Build:         map[string]string{"context": "../mongo"},
					ContainerName: "mongo",
					Command:       "node_modules/exoservice/bin/exo-js",
					Links:         []string{"mongo3.4.0:mongo"},
					Environment: map[string]string{
						"ROLE":        "mongo",
						"EXOCOM_HOST": "exocom0.24.0",
						"EXOCOM_PORT": "$EXOCOM_PORT",
						"MONGO":       "mongo",
					},
					DependsOn: []string{"exocom0.24.0", "mongo3.4.0"},
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
				appDir := path.Join("..", "..", "..", "example-apps", "complex-setup-app")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "users-service"
				dockerComposeBuilder := composebuilder.NewDockerComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not include the docker config that is defined in application.yml", func() {
				_, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(false))
			})
		})
	})

	var _ = Describe("compiles the docker compose project name properly", func() {
		expected := "spacetweet123"

		It("converts all characters to lowercase", func() {
			actual := cmd.GetDockerComposeProjectName("SpaceTweet123")
			Expect(actual).To(Equal(expected))
		})

		It("strips non-alphanumeric characters", func() {
			actual := cmd.GetDockerComposeProjectName("$Space-Tweet_123")
			Expect(actual).To(Equal(expected))
		})

		It("strips whitespace characters", func() {
			actual := cmd.GetDockerComposeProjectName("Space   Tweet  123")
			Expect(actual).To(Equal(expected))
		})
	})
})
