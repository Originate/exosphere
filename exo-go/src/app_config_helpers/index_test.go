package appConfigHelpers_test

import (
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	appDir  string
	homeDir string
)

var _ = BeforeSuite(func() {
	var err error
	homeDir, err = util.GetUserHomeDir()
	if err != nil {
		panic(err)
	}
})

var _ = Describe("GetEnvironmentVariables", func() {
	It("should return the environment variables of all dependencies", func() {
		appDir = path.Join("..", "..", "..", "exosphere-shared", "example-apps", "complex-setup-app")
		appConfig, err := types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		actual := appConfigHelpers.GetEnvironmentVariables(appConfig, appDir, homeDir)
		expected := map[string]string{"EXOCOM_PORT": "80", "DB_NAME": "test-db"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetAllDependencyNames", func() {
	It("should return the container names of all application and service dependencies", func() {
		appDir := path.Join("..", "..", "..", "exosphere-shared", "example-apps", "external-dependency")
		appConfig, err := types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		actual, err := appConfigHelpers.GetAllDependencyNames(appDir, appConfig)
		Expect(err).NotTo(HaveOccurred())
		expected := []string{"exocom0.22.1", "mongo3.4.0"}
		Expect(actual).To(Equal(expected))
	})
})
