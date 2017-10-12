package terraform_test

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompileVarFlags", func() {
	var _ = Describe("no dependencies", func() {
		service1EnvVars := types.EnvVars{
			Default: map[string]string{
				"env1": "val1",
			},
			Secrets: []string{"secret1"},
		}
		service1Config := types.ServiceConfig{
			Environment: service1EnvVars,
		}
		secrets := map[string]string{
			"secret1": "secret_value1",
		}
		deployConfig := types.DeployConfig{
			ServiceConfigs: map[string]types.ServiceConfig{
				"service1": service1Config,
			},
		}
		imageMap := map[string]string{"service1": "dummy-image"}

		It("should compile the proper var flags", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, secrets, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[0]).To(Equal("-var"))
			Expect(vars[1]).To(Equal("secret1=secret_value1"))
			Expect(vars[2]).To(Equal("-var"))
			Expect(vars[3]).To(Equal("service1_docker_image=dummy-image"))
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
					"name":  "secret1",
					"value": "secret_value1",
				},
				{
					"name":  "env1",
					"value": "val1",
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

	var _ = Describe("with exocom dependency", func() {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		It("should add the dependency service env vars to each service", func() {
			deployConfig := types.DeployConfig{
				AppConfig: types.AppConfig{
					Production: types.AppProductionConfig{
						Dependencies: []types.ProductionDependencyConfig{
							{Name: "exocom"},
						},
					},
					Name: "my-app",
				},
				ServiceConfigs: map[string]types.ServiceConfig{
					"service1": {},
				},
			}
			imageMap := map[string]string{"service1": "dummy-image"}

			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[2]).To(Equal("-var"))
			varName := strings.Split(vars[3], "=")[0]
			varVal := strings.Split(vars[3], "=")[1]
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
			err := testHelpers.CheckoutApp(cwd, "simple")
			Expect(err).NotTo(HaveOccurred())
			appDir := path.Join("tmp", "simple")
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppConfig:      appConfig,
				ServiceConfigs: serviceConfigs,
				AppDir:         appDir,
			}

			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{}, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[4]).To(Equal("-var"))
			varName := strings.Split(vars[5], "=")[0]
			varVal := strings.Split(vars[5], "=")[1]
			actualVal := `[{"receives":["users.created"],"role":"web","sends":["users.create"]}]`
			Expect(varName).To(Equal("exocom_service_routes"))
			Expect(varVal).To(Equal(actualVal))
		})
	})

	var _ = Describe("with service dependency", func() {
		deployConfig := types.DeployConfig{
			AppConfig: types.AppConfig{
				Production: types.AppProductionConfig{
					Dependencies: []types.ProductionDependencyConfig{},
				},
				Name: "my-app",
			},
			ServiceConfigs: map[string]types.ServiceConfig{
				"service1": {
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
		imageMap := map[string]string{"service1": "dummy-image"}

		It("should add the dependency service env vars to each service", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{"password-secret": "password123"}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[0]).To(Equal("-var"))
			varName := strings.Split(vars[5], "=")[0]
			varVal := strings.Split(vars[5], "=")[1]
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
