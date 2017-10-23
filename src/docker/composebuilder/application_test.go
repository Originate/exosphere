package composebuilder_test

import (
	"path"
	"runtime"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("composebuilder", func() {
	var _ = Describe("GetApplicationDockerConfigs", func() {
		var filePath string

		BeforeEach(func() {
			_, filePath, _, _ = runtime.Caller(0)
		})

		It("should create a docker-compose.yml for production", func() {
			err := testHelpers.CheckoutApp(cwd, "rds")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join(path.Dir(filePath), "tmp", "rds")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
				AppConfig: appConfig,
				AppDir:    appDir,
				BuildMode: composebuilder.BuildModeDeployProduction,
				HomeDir:   homeDir,
			})
			Expect(err).NotTo(HaveOccurred())

			By("ignoring rds dependencies")
			delete(dockerConfigs, "my-sql-service")
			Expect(len(dockerConfigs)).To(Equal(0))
		})
	})
})
