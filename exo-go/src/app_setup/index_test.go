package appSetup_test

import (
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServiceDockerConfigs", func() {

	var setup *appSetup.AppSetup
	var appDir string

	var _ = BeforeSuite(func() {
		err := appSetup.CheckoutApp("tmp", "complex-setup-app")
		appDir = path.Join("tmp", "complex-setup-app")
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(err).NotTo(HaveOccurred())
		setup, err = appSetup.NewAppSetup(appConfig, logger.NewLogger([]string{}, []string{}), appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("GetServiceDockerConfigs", func() {

		It("should starts up the app sucessfully", func() {
			err := setup.StartSetup()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include the docker config for the service itself", func() {
		})

	})

})
