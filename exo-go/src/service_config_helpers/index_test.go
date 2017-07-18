package serviceConfigHelpers_test

import (
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var appConfig types.AppConfig
var appDir string

var _ = BeforeSuite(func() {
	appDir = path.Join("..", "..", "..", "exosphere-shared", "example-apps", "complex-setup-app")
	var err error
	appConfig, err = appConfigHelpers.GetAppConfig(appDir)
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("GetServiceConfigs", func() {
	var serviceConfigs map[string]types.ServiceConfig

	It("should not return an error when all service.yml files are valid", func() {
		var err error
		serviceConfigs, err = serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should include all services", func() {
		for _, serviceName := range appConfigHelpers.GetServiceNames(appConfig.Services) {
			_, exists := serviceConfigs[serviceName]
			Expect(exists).To(Equal(true))
		}
	})

	It("should contain correct configuration for the internal service", func() {
		startup := map[string]string{
			"command":     `echo "does not run"`,
			"online-text": "does not run",
		}
		expected, err := yaml.Marshal(types.ServiceConfig{
			Type:        "html-server",
			Description: "dummy html service used for testing setup only - does not run",
			Author:      "test-author",
			Setup:       "yarn install",
			Startup:     startup,
			ServiceMessages: types.ServiceMessages{
				Sends:    []string{"todo.create"},
				Receives: []string{"todo.created"},
			},
			Docker: map[string]interface{}{
				"ports": []string{"3000:3000"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		actual, err := yaml.Marshal(serviceConfigs["html-server"])
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(expected))
	})

	It("should contain correct configuration for the external docker image", func() {
		startup := map[string]string{
			"command":     "node server.js",
			"online-text": "web server running at port",
		}
		restart := map[string]interface{}{"ignore": []string{"**/*.txt"}}
		serviceMessages := types.ServiceMessages{
			Sends:    []string{"users.list", "users.create"},
			Receives: []string{"users.listed", "users.created"},
		}
		environmentVars := map[string]string{
			"EXTERNAL_SERVICE_HOST": "external-service0.1.2",
			"EXTERNAL_SERVICE_PORT": "$EXTERNAL_SERVICE_PORT",
		}
		docker := map[string]interface{}{
			"ports":       []string{"5000:5000"},
			"volumes":     []string{"{{EXO_DATA_PATH}}:/data/db"},
			"environment": environmentVars,
		}
		expected, err := yaml.Marshal(types.ServiceConfig{
			Type:            "external-service",
			Description:     "says hello to the world, ignores .txt files when file watching",
			Author:          "exospheredev",
			Setup:           "echo 'setting up ... done'",
			Startup:         startup,
			Restart:         restart,
			ServiceMessages: serviceMessages,
			Docker:          docker,
		})
		Expect(err).ToNot(HaveOccurred())
		actual, err := yaml.Marshal(serviceConfigs["external-service"])
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(expected))
	})

})
