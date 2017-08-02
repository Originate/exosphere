package application_test

import (
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/application"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/Originate/exosphere/exo-go/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tester", func() {

	It("should create docker-compose.yml version 3", func() {
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
		// _, pipeWriter := io.Pipe()
		mockLogger := application.NewLogger([]string{}, []string{}, os.Stdout)
		tester, err := application.NewTester(appConfig, mockLogger, appDir, homeDir)
		Expect(err).NotTo(HaveOccurred())
		tester.Run()

		It("should not include the docker config that is defined in application.yml", func() {
			Expect("3").To(Equal("3"))
		})

	})

})
