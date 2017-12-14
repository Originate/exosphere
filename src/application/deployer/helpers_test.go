package deployer_test

import (
	"io/ioutil"

	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
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
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())

			dockerCompose := types.DockerCompose{
				Services: map[string]types.DockerConfig{
					"web":       types.DockerConfig{},
					"users":     types.DockerConfig{},
					"dashboard": types.DockerConfig{},
				},
			}
			deployConfig := deploy.Config{
				AppContext: appContext,
			}
			imageNames, err := deployer.GetImageNames(deployConfig, dockerCompose)
			Expect(err).NotTo(HaveOccurred())
			expectedImages := map[string]string{
				"exocom":    "originate/exocom:0.27.0",
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
