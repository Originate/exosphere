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

var _ = Describe("Secrets methods", func() {

	It("converts a tf string to a Secrets map", func() {
		tfvars := `var1="val1"
var2="val2"
var3="val3"`
		expectedMap := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
		})
		Expect(types.NewSecrets(tfvars)).To(Equal(expectedMap))
	})

	It("converts Secrets to a tf string", func() {
		secrets := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
		})
		expectedTfvars := `var1="val1"
var2="val2"
var3="val3"`
		Expect(secrets.TfString()).To(Equal(expectedTfvars))
	})

	It("merges secrets", func() {
		existingTfVars := `var1="val1"
var2="val2"
var3="val3"`
		newSecrets := types.Secrets(map[string]string{
			"var4": "val4",
		})

		expectedTfString := `var1="val1"
var2="val2"
var3="val3"
var4="val4"`
		actualSecrets := types.NewSecrets(existingTfVars).Merge(newSecrets)

		Expect(expectedTfString).To(Equal(actualSecrets.TfString()))
	})
})

var _ = Describe("ECR helpers", func() {
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
			AppDir:    appDir,
			HomeDir:   homeDir,
			AppConfig: appConfig, //TODO
		}
		imageNames, err := aws.GetImageNames(deployConfig, "./tmp", dockerCompose)
		Expect(err).NotTo(HaveOccurred())
		expectedImages := map[string]string{
			"exocom":    "originate/exocom:0.24.0",
			"users":     "tmp_users",
			"dashboard": "tmp_dashboard",
			"web":       "tmp_web",
		}

		for k, v := range expectedImages {
			Expect(imageNames[k]).To(Equal(v))
		}
	})
})
