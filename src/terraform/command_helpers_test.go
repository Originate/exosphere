package terraform_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompileVarFlags", func() {
	var _ = Describe("public service with no dependencies", func() {
		service1EnvVars := types.EnvVars{
			Default: map[string]string{
				"env1": "val1",
			},
			Secrets: []string{"secret1"},
		}
		service1Config := types.ServiceConfig{
			Type:        "public",
			Environment: service1EnvVars,
		}
		service2Config := types.ServiceConfig{
			Type: "public",
			Production: types.ServiceProductionConfig{
				Port: "80",
				URL:  "my-test-url.com",
			},
		}
		secrets := map[string]string{
			"secret1": "secret_value1",
		}
		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: types.AppConfig{
					Name: "my-app",
					Services: map[string]types.ServiceSource{
						"service1": types.ServiceSource{},
						"service2": types.ServiceSource{},
					},
				},
				ServiceContexts: map[string]*context.ServiceContext{
					"service1": {
						Config: service1Config,
					},
					"service2": {
						Config: service2Config,
					},
				},
			},
			BuildMode: types.BuildMode{
				Type:        types.BuildModeTypeDeploy,
				Environment: types.BuildModeEnvironmentProduction,
			},
		}
		imageMap := map[string]string{"service1": "dummy-image"}

		It("should compile the proper var flags", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, secrets, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(util.DoesStringArrayContain(vars, "secret1=secret_value1"))
			Expect(util.DoesStringArrayContain(vars, "service1_docker_image=dummy-image"))

			service1ExpectedValue := []map[string]string{
				{
					"name":  "ROLE",
					"value": "service1",
				},
				{
					"name":  "secret1",
					"value": "secret_value1",
				},
				{
					"name":  "env1",
					"value": "val1",
				},
				{
					"name":  "SERVICE2_EXTERNAL_ORIGIN",
					"value": "https://my-test-url.com",
				},
				{
					"name":  "SERVICE2_ORIGIN",
					"value": "http://service2.local",
				},
			}
			var service1ActualString string
			for _, varFlag := range vars {
				if strings.Contains(varFlag, "service1_env_vars") {
					service1ActualString = strings.Split(varFlag, "=")[1]
				}
			}
			actualValue := []map[string]string{}
			var escapedValue string
			err = json.Unmarshal([]byte(service1ActualString), &escapedValue)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedValue), &actualValue)
			Expect(err).NotTo(HaveOccurred())
			Expect(service1ExpectedValue).To(ConsistOf(actualValue))
		})
	})

	var _ = Describe("with exocom dependency", func() {
		It("compile the proper var flags", func() {
			deployConfig := deploy.Config{
				AppContext: &context.AppContext{
					Config: types.AppConfig{
						Remote: types.RemoteConfig{
							Dependencies: []types.RemoteDependency{
								{Name: "exocom"},
							},
						},
						Name: "my-app",
					},
					ServiceContexts: map[string]*context.ServiceContext{
						"service1": {},
					},
				},
			}
			imageMap := map[string]string{"service1": "dummy-image", "exocom": "originate/exocom:0.0.1"}

			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[2]).To(Equal("-var"))
			exocomVarFlag := strings.Split(vars[3], "=")[0]
			exocomDockerImageName := strings.Split(vars[3], "=")[1]
			Expect(exocomVarFlag).To(Equal("exocom_docker_image"))
			Expect(exocomDockerImageName).To(Equal("originate/exocom:0.0.1"))
			Expect(vars[4]).To(Equal("-var"))
			varName := strings.Split(vars[5], "=")[0]
			varVal := strings.Split(vars[5], "=")[1]
			var escapedVal string
			actualVal := []map[string]string{}
			expectedVal := []map[string]string{
				{
					"name":  "ROLE",
					"value": "service1",
				},
				{
					"name":  "EXOCOM_HOST",
					"value": "exocom.my-app.local",
				},
			}
			err = json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			Expect(varName).To(Equal("service1_env_vars"))
			Expect(expectedVal).To(ConsistOf(actualVal))
		})

		It("should compile the dependency terraform vars", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "simple")
			Expect(err).NotTo(HaveOccurred())
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := deploy.Config{
				AppContext: appContext,
			}

			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{}, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[4]).To(Equal("-var"))
			varFlagName := strings.Split(vars[7], "=")[0]
			varFlagValue := strings.Split(vars[7], "=")[1]
			var escapedFlagValue1 string
			var escapedFlagValue2 []map[string]string
			var actualValue string
			expectedValue := `[{"receives":["users.created"],"role":"web","sends":["users.create"]}]`

			err = json.Unmarshal([]byte(varFlagValue), &escapedFlagValue1)
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
		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: types.AppConfig{
					Remote: types.RemoteConfig{
						Dependencies: []types.RemoteDependency{},
					},
					Name: "my-app",
				},
				ServiceContexts: map[string]*context.ServiceContext{
					"service1": {
						Config: types.ServiceConfig{
							Production: types.ServiceProductionConfig{
								Dependencies: []types.RemoteDependency{
									{
										Config: types.RemoteDependencyConfig{
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
				},
			},
		}
		imageMap := map[string]string{"service1": "dummy-image"}

		It("should add the dependency service env vars to each service", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{"password-secret": "password123"}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[6]).To(Equal("-var"))
			varName := strings.Split(vars[7], "=")[0]
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
					"value": "service1",
				},
			}
			err = json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			Expect(varName).To(Equal("service1_env_vars"))
			Expect(expectedVal).To(ConsistOf(actualVal))
		})
	})
})
