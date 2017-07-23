package dockerSetup_test

import (
	"io"
	"path"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServiceDockerConfigs", func() {
	var homeDir string

	var _ = BeforeSuite(func() {
		var err error
		homeDir, err = osHelpers.GetUserHomeDir()
		if err != nil {
			panic(err)
		}
	})

	var _ = Describe("GetServiceDockerConfigs", func() {

		var _ = Describe("unshared docker configs", func() {
			var dockerConfigs map[string]types.DockerConfig

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "exosphere-shared", "example-apps", "external-dependency")
				appConfig, err := appConfigHelpers.GetAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := serviceConfigHelpers.GetServiceData(appConfig.Services)
				serviceName := "mongo"
				_, pipeWriter := io.Pipe()
				setup := &dockerSetup.DockerSetup{
					AppConfig:     appConfig,
					ServiceConfig: serviceConfigs[serviceName],
					ServiceData:   serviceData[serviceName],
					Role:          serviceName,
					Logger:        logger.NewLogger([]string{}, []string{}, pipeWriter),
					AppDir:        appDir,
					HomeDir:       homeDir,
				}
				dockerConfigs, err = setup.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should include the docker config for the service itself", func() {
				dockerConfig, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(true))
				Expect(dockerConfig).To(Equal(types.DockerConfig{
					Build:         "../mongo",
					ContainerName: "mongo",
					Command:       "node_modules/exoservice/bin/exo-js",
					Links:         []string{"mongo3.4.0:mongo"},
					Environment: map[string]string{
						"ROLE":        "mongo",
						"EXOCOM_HOST": "exocom0.22.1",
						"EXOCOM_PORT": "$EXOCOM_PORT",
						"MONGO":       "mongo",
					},
					DependsOn: []string{"mongo3.4.0", "exocom0.22.1"},
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
			var dockerConfigs map[string]types.DockerConfig

			var _ = BeforeEach(func() {
				appDir := path.Join("..", "..", "..", "exosphere-shared", "example-apps", "complex-setup-app")
				appConfig, err := appConfigHelpers.GetAppConfig(appDir)
				Expect(err).NotTo(HaveOccurred())
				serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
				Expect(err).NotTo(HaveOccurred())
				serviceData := serviceConfigHelpers.GetServiceData(appConfig.Services)
				serviceName := "users-service"
				_, pipeWriter := io.Pipe()
				setup := &dockerSetup.DockerSetup{
					AppConfig:     appConfig,
					ServiceConfig: serviceConfigs[serviceName],
					ServiceData:   serviceData[serviceName],
					Role:          serviceName,
					Logger:        logger.NewLogger([]string{}, []string{}, pipeWriter),
					AppDir:        appDir,
					HomeDir:       homeDir,
				}
				dockerConfigs, err = setup.GetServiceDockerConfigs()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not include the docker config that is defined in application.yml", func() {
				_, exists := dockerConfigs["mongo"]
				Expect(exists).To(Equal(false))
			})
		})
	})
})
