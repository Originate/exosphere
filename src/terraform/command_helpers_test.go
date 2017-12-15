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

func ArrayHasStringMap(haystack []map[string]string, needle map[string]string) bool {
	for _, item := range haystack {
		if reflect.DeepEqual(item, needle) {
			return true
		}
	}
	return false
}

var _ = Describe("GetVarMap", func() {
	var _ = Describe("public service with no dependencies", func() {
		service1Config := types.ServiceConfig{
			Type: "public",
			Remote: types.ServiceRemoteConfig{
				Environment: map[string]string{
					"env1": "val1",
				},
				Secrets: []string{"secret1"},
			},
		}
		service2Config := types.ServiceConfig{
			Type: "public",
			Production: types.ServiceProductionConfig{
				Port: "80",
			},
			Remote: types.ServiceRemoteConfig{
				URL: "my-test-url.com",
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
		}
		imageMap := map[string]string{
			"service1": "dummy-image1",
			"service2": "dummy-image2",
		}

		It("should compile the proper var flags", func() {
			varMap, err := terraform.GetVarMap(deployConfig, secrets, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap["service1_docker_image"]).To(Equal("dummy-image1"))
			Expect(varMap["service2_docker_image"]).To(Equal("dummy-image2"))

			expectedService1EnvVars := []map[string]string{
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
					"name":  "SERVICE2_INTERNAL_ORIGIN",
					"value": "http://service2.local",
				},
			}
			actualService1Value := []map[string]string{}
			var escapedValue string
			err = json.Unmarshal([]byte(varMap["service1_env_vars"]), &escapedValue)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedValue), &actualService1Value)
			Expect(err).NotTo(HaveOccurred())
			for _, actualEnvVar := range actualService1Value {
				Expect(ArrayHasStringMap(expectedService1EnvVars, actualEnvVar)).To(BeTrue())
			}
		})
	})

	var _ = Describe("with exocom dependency", func() {
		It("compile the proper var flags", func() {
			deployConfig := deploy.Config{
				AppContext: &context.AppContext{
					Config: types.AppConfig{
						Remote: types.AppRemoteConfig{
							Dependencies: map[string]types.RemoteDependency{
								"exocom": types.RemoteDependency{
									Type:   "exocom",
									Config: types.RemoteDependencyConfig{},
								},
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

			varMap, err := terraform.GetVarMap(deployConfig, map[string]string{}, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap["service1_docker_image"]).To(Equal("dummy-image"))
			Expect(varMap["exocom_docker_image"]).To(Equal("originate/exocom:0.0.1"))
			expectedService1EnvVars := []map[string]string{
				{
					"name":  "ROLE",
					"value": "service1",
				},
				{
					"name":  "EXOCOM_HOST",
					"value": "exocom.my-app.local",
				},
			}
			actualService1Value := []map[string]string{}
			var escapedValue string
			err = json.Unmarshal([]byte(varMap["service1_env_vars"]), &escapedValue)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedValue), &actualService1Value)
			Expect(err).NotTo(HaveOccurred())
			for _, actualEnvVar := range actualService1Value {
				Expect(ArrayHasStringMap(expectedService1EnvVars, actualEnvVar)).To(BeTrue())
			}
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

			varMap, err := terraform.GetVarMap(deployConfig, map[string]string{}, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			_, ok := varMap["exocom_env_vars"]
			Expect(ok).To(BeTrue())
			var str string
			var actualDependencyVar []map[string]interface{}
			err = json.Unmarshal([]byte(varMap["exocom_env_vars"]), &str)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(str), &actualDependencyVar)
			Expect(err).NotTo(HaveOccurred())
			expectedValue := `{"web":{"receives":["users.created"],"sends":["users.create"]}}`
			Expect(reflect.DeepEqual(actualDependencyVar[0]["value"], expectedValue)).To(BeTrue())
		})
	})

	var _ = Describe("with service dependency", func() {
		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: types.AppConfig{
					Remote: types.AppRemoteConfig{
						Dependencies: map[string]types.RemoteDependency{},
					},
					Name: "my-app",
				},
				ServiceContexts: map[string]*context.ServiceContext{
					"service1": {
						Config: types.ServiceConfig{
							Remote: types.ServiceRemoteConfig{
								Environment: map[string]string{
									"TEST_APP_ENV": "TEST_APP_ENV_VAL",
								},
								Secrets: []string{"password-secret"},
								Dependencies: map[string]types.RemoteDependency{
									"postgres": types.RemoteDependency{
										Type: "rds",
										Config: types.RemoteDependencyConfig{
											Rds: types.RdsConfig{
												Engine:             "test-engine",
												EngineVersion:      "0.0.1",
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
									},
								},
							},
						},
					},
				},
			},
		}

		It("should add the dependency service env vars to each service", func() {
			varMap, err := terraform.GetVarMap(deployConfig, map[string]string{"password-secret": "password123"}, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			expectedService1EnvVars := []map[string]string{
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
				{
					"name":  "TEST_APP_ENV",
					"value": "TEST_APP_ENV_VAL",
				},
				{
					"name":  "password-secret",
					"value": "password123",
				},
			}
			actualService1Value := []map[string]string{}
			var escapedValue string
			err = json.Unmarshal([]byte(varMap["service1_env_vars"]), &escapedValue)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedValue), &actualService1Value)
			Expect(err).NotTo(HaveOccurred())
			for _, actualEnvVar := range actualService1Value {
				Expect(ArrayHasStringMap(expectedService1EnvVars, actualEnvVar)).To(BeTrue())
			}
		})
	})
})
