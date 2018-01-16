package terraform_test

import (
	"encoding/json"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServicesVarMap", func() {
	var _ = Describe("no dependencies", func() {
		service1Config := types.ServiceConfig{
			Type: "public",
			Remote: types.ServiceRemoteConfig{
				Environments: map[string]types.ServiceRemoteEnvironment{
					"qa": {
						EnvironmentVariables: map[string]string{
							"env1": "val1",
						},
						Secrets: []string{"secret1"},
						URL:     "service1.example.com",
					},
				},
			},
		}
		service2Config := types.ServiceConfig{
			Type: "public",
			Production: types.ServiceProductionConfig{
				Port: "80",
			},
			Remote: types.ServiceRemoteConfig{
				Environments: map[string]types.ServiceRemoteEnvironment{
					"qa": {
						URL: "service2.example.com",
					},
				},
			},
		}
		secrets := map[string]string{
			"secret1": "secret_value1",
		}
		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: &types.AppConfig{
					Name: "my-app",
					Services: map[string]types.ServiceSource{
						"service1": types.ServiceSource{},
						"service2": types.ServiceSource{},
					},
					Remote: types.AppRemoteConfig{
						Environments: map[string]types.AppRemoteEnvironment{
							"qa": {
								URL: "app.example.com",
							},
						},
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
			AwsConfig: types.AwsConfig{
				Profile:           "my_profile",
				Region:            "my_region",
				AccountID:         "123",
				SslCertificateArn: "456",
			},
			RemoteEnvironmentID: "qa",
		}
		imageMap := map[string]string{
			"service1": "dummy-image1",
			"service2": "dummy-image2",
		}

		It("should compile the proper var flags", func() {
			varMap, err := terraform.GetServicesVarMap(deployConfig, secrets, imageMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap["service1_docker_image"]).To(Equal("dummy-image1"))
			Expect(varMap["service2_docker_image"]).To(Equal("dummy-image2"))
			Expect(varMap["aws_profile"]).To(Equal("my_profile"))
			Expect(varMap["aws_region"]).To(Equal("my_region"))
			Expect(varMap["aws_account_id"]).To(Equal("123"))
			Expect(varMap["aws_ssl_certificate_arn"]).To(Equal("456"))
			Expect(varMap["application_url"]).To(Equal("app.example.com"))
			Expect(varMap["env"]).To(Equal("qa"))
			Expect(varMap["service1_url"]).To(Equal("service1.example.com"))
			Expect(varMap["service2_url"]).To(Equal("service2.example.com"))
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
					"value": "https://service2.example.com",
				},
				{
					"name":  "SERVICE2_INTERNAL_ORIGIN",
					"value": "http://service2.qa-my-app.local",
				},
			}
			actualService1Value := []map[string]string{}
			err = json.Unmarshal([]byte(varMap["service1_env_vars"]), &actualService1Value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualService1Value).To(ConsistOf(expectedService1EnvVars))
		})
	})

	var _ = Describe("with service dependency", func() {
		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: &types.AppConfig{
					Remote: types.AppRemoteConfig{
						Dependencies: map[string]types.RemoteDependency{},
					},
					Name: "my-app",
				},
				ServiceContexts: map[string]*context.ServiceContext{
					"service1": {
						Config: types.ServiceConfig{
							Remote: types.ServiceRemoteConfig{
								Environments: map[string]types.ServiceRemoteEnvironment{
									"qa": {
										EnvironmentVariables: map[string]string{
											"RDS_HOST": "rds.qa-my-app.local",
											"DB_NAME":  "test-db",
											"DB_USER":  "test-user",
										},
										Secrets: []string{"DB_PASS"},
									},
								},
								Dependencies: map[string]types.RemoteDependency{
									"postgres": types.RemoteDependency{
										Type: "rds",
									},
								},
							},
						},
					},
				},
			},
			RemoteEnvironmentID: "qa",
		}

		It("should add the dependency service env vars to each service", func() {
			varMap, err := terraform.GetServicesVarMap(deployConfig, map[string]string{"DB_PASS": "password123"}, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			expectedService1EnvVars := []map[string]string{
				{
					"name":  "RDS_HOST",
					"value": "rds.qa-my-app.local",
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
					"name":  "DB_PASS",
					"value": "password123",
				},
			}
			actualService1Value := []map[string]string{}
			err = json.Unmarshal([]byte(varMap["service1_env_vars"]), &actualService1Value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualService1Value).To(ConsistOf(expectedService1EnvVars))
		})
	})
})
