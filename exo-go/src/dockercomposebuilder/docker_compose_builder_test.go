package dockercomposebuilder_test

import (
	"path"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/dockercompose"
	"github.com/Originate/exosphere/exo-go/src/dockercomposebuilder"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ComposeBuilder", func() {
	var _ = Describe("GetServiceDockerConfigs", func() {
		var _ = Describe("unshared docker configs", func() {
			var dockerConfigs dockercompose.DockerConfigs

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "example-apps", "external-dependency")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "mongo"
				homeDir, err := util.GetHomeDirectory()
				Expect(err).NotTo(HaveOccurred())
				dockerComposeBuilder := dockercomposebuilder.NewComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the docker config for the service itself", func() {
				dockerConfig, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(true))
				Expect(dockerConfig).To(Equal(dockercompose.DockerConfig{
					Build:         map[string]string{"context": "../mongo"},
					ContainerName: "mongo",
					Command:       "node_modules/exoservice/bin/exo-js",
					Links:         []string{"mongo3.4.0:mongo"},
					Environment: map[string]string{
						"ROLE":        "mongo",
						"EXOCOM_HOST": "exocom0.22.1",
						"EXOCOM_PORT": "$EXOCOM_PORT",
						"MONGO":       "mongo",
					},
					DependsOn: []string{"exocom0.22.1", "mongo3.4.0"},
				}))
			})

			It("should include the docker configs for the service's dependencies", func() {
				dockerConfig, exists := dockerConfigs["mongo3.4.0"]
				Expect(exists).To(Equal(true))
				volumesRegex := regexp.MustCompile(`./\.exosphere/Exosphere-application-with-a-third-party-dependency/mongo/data:/data/db`)
				Expect(volumesRegex.MatchString(dockerConfig.Volumes[0])).To(Equal(true))
				dockerConfig.Volumes = nil
				Expect(dockerConfig).To(Equal(dockercompose.DockerConfig{
					Image:         "mongo:3.4.0",
					ContainerName: "mongo3.4.0",
					Ports:         []string{"27017:27017"},
				}))
			})
		})

		var _ = Describe("shared docker configs", func() {
			var dockerConfigs dockercompose.DockerConfigs

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "example-apps", "complex-setup-app")
				appConfig, err := types.NewAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := appConfig.GetServiceData()
				serviceName := "users-service"
				homeDir, err := util.GetHomeDirectory()
				Expect(err).NotTo(HaveOccurred())
				dockerComposeBuilder := dockercomposebuilder.NewComposeBuilder(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, appDir, homeDir)
				dockerConfigs, err = dockerComposeBuilder.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not include the docker config that is defined in application.yml", func() {
				_, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(false))
			})
		})
	})
})
