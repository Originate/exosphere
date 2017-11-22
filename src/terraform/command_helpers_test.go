package terraform_test

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompileVarFlags", func() {
	var appContext types.AppContext

	BeforeEach(func() {
		appDir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		err = helpers.CheckoutApp(appDir, "simple")
		Expect(err).NotTo(HaveOccurred())
		appConfig, err := types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		appContext = types.AppContext{
			Location: appDir,
			Config:   appConfig,
		}
		application.GenerateComposeFiles(appContext)
	})

	var _ = Describe("a simple application", func() {
		var vars []string

		BeforeEach(func() {
			serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
			Expect(err).NotTo(HaveOccurred())
			deployConfig := types.DeployConfig{
				AppContext:       appContext,
				ServiceConfigs:   serviceConfigs,
				DockerComposeDir: path.Join(appContext.Location, "docker-compose"),
			}
			secrets := map[string]string{
				"API_KEY": "secret_api_key",
			}
			imageMap := map[string]string{"web": "dummy-image"}
			vars, err = terraform.CompileVarFlags(deployConfig, secrets, imageMap)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should compile the proper var flags", func() {
			Expect(vars[0]).To(Equal("-var"))
			Expect(vars[1]).To(Equal("API_KEY=secret_api_key"))
			Expect(vars[2]).To(Equal("-var"))
			Expect(vars[3]).To(Equal("web_docker_image=dummy-image"))
			Expect(vars[4]).To(Equal("-var"))

			varName := strings.Split(vars[7], "=")[0]
			Expect(varName).To(Equal("web_env_vars"))
			varVal := strings.Split(vars[7], "=")[1]
			var escapedVal string
			actualVal := []map[string]string{}
			expectedVal := []map[string]string{
				{
					"name":  "ROLE",
					"value": "web",
				},
				{
					"name":  "API_KEY",
					"value": "secret_api_key",
				},
				{
					"name":  "EXOCOM_HOST",
					"value": "exocom.simple.local",
				},
				{
					"name":  "ENV1",
					"value": "value1",
				},
			}
			err := json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			Expect(expectedVal).To(ConsistOf(actualVal))
		})

		It("should compile the dependency terraform vars", func() {
			varFlagName := strings.Split(vars[9], "=")[0]
			varFlagValue := strings.Split(vars[9], "=")[1]
			var escapedFlagValue1 string
			var escapedFlagValue2 []map[string]string
			var actualValue string
			expectedValue := `[{"receives":["users.created"],"role":"web","sends":["users.create"]}]`

			err := json.Unmarshal([]byte(varFlagValue), &escapedFlagValue1)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedFlagValue1), &escapedFlagValue2)
			Expect(err).NotTo(HaveOccurred())
			for k, v := range escapedFlagValue2[0] {
				if k == "value" {
					actualValue = v
				}
			}

			Expect(varFlagName).To(Equal("exocom_env_vars"))
			Expect(actualValue).To(Equal(expectedValue))
		})
	})

	var _ = Describe("with service dependency", func() {
		It("should add the dependency service env vars to each service", func() {
			imageMap := map[string]string{"web": "dummy-image"}
			deployConfig := types.DeployConfig{
				DockerComposeDir: path.Join(appContext.Location, "docker-compose"),
				AppContext: types.AppContext{
					Config: types.AppConfig{
						Production: types.AppProductionConfig{
							Dependencies: []types.ProductionDependencyConfig{},
						},
						Name: "my-app",
					},
				},
				ServiceConfigs: map[string]types.ServiceConfig{
					"web": {
						Production: types.ServiceProductionConfig{
							Dependencies: []types.ProductionDependencyConfig{
								{
									Config: types.ProductionDependencyConfigOptions{
										Rds: types.RdsConfig{
											Username:           "test-user",
											DbName:             "test-db",
											PasswordSecretName: "password-secret",
											ServiceEnvVarNames: types.ServiceEnvVarNames{
												DbName:   "DB_NAME",
												Username: "DB_USER",
												Password: "DB_PASS",
											},
										},
									},
									Name:    "postgres",
									Version: "0.0.1",
								},
							},
						},
					},
				},
			}
			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{"password-secret": "password123"}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[6]).To(Equal("-var"))
			varName := strings.Split(vars[7], "=")[0]
			Expect(varName).To(Equal("web_env_vars"))
			varVal := strings.Split(vars[7], "=")[1]
			var escapedVal string
			actualVal := []map[string]string{}
			expectedVal := []map[string]string{
				{
					"name":  "DB_PASS",
					"value": "password123",
				},
				{
					"name":  "POSTGRES",
					"value": "test-db.my-app.local",
				},
				{
					"name":  "DB_NAME",
					"value": "test-db",
				},
				{
					"name":  "DB_USER",
					"value": "test-user",
				},
				{
					"name":  "ROLE",
					"value": "web",
				},
			}
			err = json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			for _, expectedElement := range expectedVal {
				Expect(actualVal).To(ContainElement(expectedElement))
			}
		})
	})
})
