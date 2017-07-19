package dockerSetup_test

import (
	"path"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServiceDockerConfigs", func() {

	var setup *dockerSetup.DockerSetup
	var appDir string

	var _ = BeforeSuite(func() {
		appDir = path.Join("..", "..", "..", "exosphere-shared", "example-apps", "external-dependency")
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
		Expect(err).NotTo(HaveOccurred())
		serviceData := serviceConfigHelpers.GetServiceData(appConfig.Services)
		serviceName := "mongo"
		setup = dockerSetup.NewDockerSetup(appConfig, serviceConfigs[serviceName], serviceData[serviceName], serviceName, logger.NewLogger([]string{}, []string{}), appDir)
	})

	var _ = Describe("GetServiceDockerConfigs", func() {
		var dockerConfigs map[string]types.DockerConfig

		It("should get docker configs sucessfully", func() {
			var err error
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

		It("should include the docker config for the service's dependencies", func() {
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

})
