package terraform_test

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetInfrastructureVarMap", func() {
	var _ = Describe("with exocom dependency", func() {
		It("should compile the dependency terraform vars", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "simple")
			Expect(err).NotTo(HaveOccurred())
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			deployConfig := deploy.Config{
				AppContext: appContext,
				AwsConfig: types.AwsConfig{
					Profile:           "my_profile",
					Region:            "my_region",
					AccountID:         "123",
					SslCertificateArn: "456",
				},
				RemoteEnvironmentID: "qa",
			}

			varMap, err := terraform.GetInfrastructureVarMap(deployConfig, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap["aws_profile"]).To(Equal("my_profile"))
			Expect(varMap["aws_region"]).To(Equal("my_region"))
			Expect(varMap["aws_account_id"]).To(Equal("123"))
			Expect(varMap["aws_ssl_certificate_arn"]).To(Equal("456"))
			Expect(varMap["application_url"]).To(Equal("test-url.com"))
			Expect(varMap["env"]).To(Equal("qa"))
			Expect(varMap).To(HaveKey("exocom_env_vars"))
			var actualDependencyVar []map[string]interface{}
			err = json.Unmarshal([]byte(varMap["exocom_env_vars"]), &actualDependencyVar)
			Expect(err).NotTo(HaveOccurred())
			expectedValue := `{"web":{"receives":["users.created"],"sends":["users.create"]}}`
			Expect(actualDependencyVar[0]["name"]).To(Equal("SERVICE_DATA"))
			Expect(reflect.DeepEqual(actualDependencyVar[0]["value"], expectedValue)).To(BeTrue())
		})
	})
})
