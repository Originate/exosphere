package deployer_test

import (
	"io/ioutil"

	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployer helpers", func() {

	var _ = Describe("GetImageNames", func() {
		It("compiles the list of image names", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "test")
			Expect(err).NotTo(HaveOccurred())
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
				AppContext: types.AppContext{
					Location: appDir,
					Config:   appConfig,
				},
				DockerComposeProjectName: "appname",
			}
			imageNames, err := deployer.GetImageNames(deployConfig, "./tmp", dockerCompose)
			Expect(err).NotTo(HaveOccurred())
			expectedImages := map[string]string{
				"exocom":    "originate/exocom:0.26.1",
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
