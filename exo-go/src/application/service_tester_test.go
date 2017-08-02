package application_test

import (
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/application"
	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/Originate/exosphere/exo-go/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initializer", func() {
	It("should create a docker.compose.yml", func() {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = testHelpers.CheckoutApp(cwd, "tests-complex")
		Expect(err).NotTo(HaveOccurred())
		appDir := path.Join("tmp", "tests-complex")
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		appConfig, err := types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
		Expect(err).NotTo(HaveOccurred())
		serviceData := config.GetServiceData(appConfig.Services)
		// _, pipeWriter := io.Pipe()
		mockLogger := application.NewLogger([]string{}, []string{}, os.Stdout)

		By("should output no error to the channel when all tests pass")
		serviceName := "tweets-service"
		validateTestStatus(serviceName, false, appConfig, serviceConfigs[serviceName], serviceData[serviceName], mockLogger, appDir, homeDir)

		By("should output an error to the channel when tests fail")
		serviceName = "users-service"
		validateTestStatus(serviceName, true, appConfig, serviceConfigs[serviceName], serviceData[serviceName], mockLogger, appDir, homeDir)

		By("should run tests with complex dependencies successfully")
		serviceName := "exoservice"
		validateTestStatus(serviceName, false, appConfig, serviceConfigs[serviceName], serviceData[serviceName], mockLogger, appDir, homeDir)
	})
})

func validateTestStatus(serviceName string, expectToFail bool, appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, logger *application.Logger, appDir, homeDir string) {
	ErrChannel := make(chan error)
	builtDependencies := config.GetServiceBuiltDependencies(serviceConfig, appConfig, appDir, homeDir)
	initializer, err := application.NewInitializer(appConfig, logger, appDir, homeDir)
	Expect(err).NotTo(HaveOccurred())
	runner, err := application.NewRunner(appConfig, logger, appDir, homeDir)
	Expect(err).NotTo(HaveOccurred())
	serviceTester, err := application.NewServiceTester(serviceName, serviceConfig, builtDependencies, appDir, serviceData.Location, initializer, runner)
	Expect(err).NotTo(HaveOccurred())
	done := make(chan bool)
	err = nil
	go func() {
		err = <-ErrChannel
		done <- true
	}()
	serviceTester.Run(ErrChannel)
	<-done
	if expectToFail {
		Expect(err).To(HaveOccurred())
	} else {
		Expect(err).ToNot(HaveOccurred())
	}
}
