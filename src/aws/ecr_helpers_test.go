package aws_test

import (
	"os"
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ECR helpers", func() {

	var _ = Describe("Get image names", func() {
		It("compiles the list of image names", func() {
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			err = testHelpers.CheckoutApp(cwd, "test")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join("tmp", "test")
			homeDir, err := util.GetHomeDirectory()
			if err != nil {
				panic(err)
			}
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())

			dockerCompose := types.DockerCompose{
				Services: map[string]types.DockerConfig{
					"web":       types.DockerConfig{},
					"users":     types.DockerConfig{},
					"dashboard": types.DockerConfig{},
				},
			}
			deployConfig := types.DeployConfig{
				AppDir:                   appDir,
				HomeDir:                  homeDir,
				AppConfig:                appConfig,
				DockerComposeProjectName: "appname",
			}
			imageNames, err := aws.GetImageNames(deployConfig, "./tmp", dockerCompose)
			Expect(err).NotTo(HaveOccurred())
			expectedImages := map[string]string{
				"exocom":    "originate/exocom:0.24.0",
				"users":     "appname_users",
				"dashboard": "appname_dashboard",
				"web":       "appname_web",
			}

			for k, v := range expectedImages {
				Expect(imageNames[k]).To(Equal(v))
			}
		})
	})
})
