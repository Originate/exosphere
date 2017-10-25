package types_test

import (
	"path"

	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppConfig", func() {
	var appConfig types.AppConfig

	var _ = Describe("ValidateFields", func() {
		It("should throw an error when AppConfig is missing fields in production", func() {
			appConfig = types.AppConfig{
				Production: types.AppProductionConfig{
					URL:       "originate.com",
					AccountID: "123",
					Region:    "us-west-2",
				},
			}
			err := appConfig.Production.ValidateFields()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "application.yml missing required field 'production.SslCertificateArn'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("should not throw an error when AppConfig isn't missing fields", func() {
			appConfig = types.AppConfig{
				Production: types.AppProductionConfig{
					URL:               "originate.com",
					AccountID:         "123",
					Region:            "us-west-2",
					SslCertificateArn: "cert-arn",
				},
			}
			err := appConfig.Production.ValidateFields()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	var _ = Describe("NewAppConfig", func() {
		It("should throw and error if app name is invalid", func() {
			appDir := path.Join("..", "..", "example-apps", "invalid-app-name")
			_, err := types.NewAppConfig(appDir)
			Expect(err).To(HaveOccurred())
			expectedErrorString := "The 'name' field 'invalid app' in application.yml is invalid. Only lowercase alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-z0-9]+(-[a-z0-9]+)*$/"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("should throw and error if any service keys are invalid", func() {
			appDir := path.Join("..", "..", "example-apps", "invalid-app-service")
			_, err := types.NewAppConfig(appDir)
			Expect(err).To(HaveOccurred())
			expectedErrorString := "The 'services.public' key 'invalid-service!' in application.yml is invalid. Only alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$/"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})
	})

	var _ = Describe("GetAppConfig", func() {
		BeforeEach(func() {
			appDir := path.Join("..", "..", "example-apps", "complex-setup-app")
			var err error
			appConfig, err = types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include name, version and description", func() {
			Expect(appConfig.Name).To(Equal("complex-setup-app"))
			Expect(appConfig.Description).To(Equal("An app with complex setup used for testing"))
			Expect(appConfig.Version).To(Equal("0.0.1"))
		})

		It("should have all the dependencies", func() {
			Expect(appConfig.Development.Dependencies).To(Equal([]types.DevelopmentDependencyConfig{
				types.DevelopmentDependencyConfig{
					Name:    "exocom",
					Version: "0.26.1",
				},
				types.DevelopmentDependencyConfig{
					Name:    "mongo",
					Version: "3.4.0",
					Config: types.DevelopmentDependencyConfigOptions{
						Ports:                 []string{"4000:4000"},
						Volumes:               []string{"{{EXO_DATA_PATH}}:/data/db"},
						DependencyEnvironment: map[string]string{"DB_NAME": "test-db"},
						ServiceEnvironment:    map[string]string{"COLLECTION_NAME": "test-collection"},
					},
				},
			}))
		})

		It("should have all the services", func() {
			workerService := map[string]types.ServiceData{
				"todo-service": types.ServiceData{Location: "./todo-service"},
			}
			privateServices := map[string]types.ServiceData{
				"users-service": types.ServiceData{
					MessageTranslations: []types.MessageTranslation{
						types.MessageTranslation{
							Public:   "users create",
							Internal: "mongo create",
						},
					},
					Location: "./users-service"},
			}
			publicServices := map[string]types.ServiceData{
				"html-server":      types.ServiceData{Location: "./html-server"},
				"external-service": types.ServiceData{DockerImage: "originate/test-web-server"},
			}
			expected := types.Services{Private: privateServices, Public: publicServices, Worker: workerService}
			Expect(appConfig.Services).To(Equal(expected))
		})
	})

	var _ = Describe("GetDevelopmentDependencyNames", func() {
		It("should return the names of all application dependencies", func() {
			appConfig = types.AppConfig{
				Development: types.AppDevelopmentConfig{Dependencies: []types.DevelopmentDependencyConfig{
					{Name: "exocom"},
					{Name: "mongo"},
				},
				},
			}
			actual := appConfig.GetDevelopmentDependencyNames()
			expected := []string{"exocom", "mongo"}
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetSortedServiceRoles", func() {
		It("should return the names of all services in alphabetical order", func() {
			appConfig = types.AppConfig{
				Services: types.Services{
					Worker: map[string]types.ServiceData{
						"worker-service-1": types.ServiceData{},
					},
					Private: map[string]types.ServiceData{
						"private-service-1": types.ServiceData{},
					},
					Public: map[string]types.ServiceData{
						"public-service-1": types.ServiceData{},
						"public-service-2": types.ServiceData{},
					},
				},
			}
			actual := appConfig.GetSortedServiceRoles()
			expected := []string{"private-service-1", "public-service-1", "public-service-2", "worker-service-1"}
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetSilencedDevelopmentDependencyNames", func() {
		It("should return the names of all silenced dependencies", func() {
			appConfig = types.AppConfig{
				Development: types.AppDevelopmentConfig{
					Dependencies: []types.DevelopmentDependencyConfig{
						{Name: "exocom", Silent: true},
						{Name: "mongo"},
					},
				},
			}
			actual := appConfig.GetSilencedDevelopmentDependencyNames()
			expected := []string{"exocom"}
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetSilencedServiceRoles", func() {
		It("should return the names of all silenced services", func() {
			appConfig = types.AppConfig{
				Services: types.Services{
					Worker: map[string]types.ServiceData{
						"worker-service-1": types.ServiceData{Silent: true},
					},
					Private: map[string]types.ServiceData{
						"private-service-1": types.ServiceData{Silent: true},
						"private-service-2": types.ServiceData{},
					},
					Public: map[string]types.ServiceData{
						"public-service-1": types.ServiceData{},
						"public-service-2": types.ServiceData{Silent: true},
					},
				},
			}
			actual := appConfig.GetSilencedServiceRoles()
			expected := []string{"worker-service-1", "private-service-1", "public-service-2"}
			Expect(actual).To(ConsistOf(expected))
		})
	})

	var _ = Describe("VerifyServiceRoleDoesNotExist", func() {
		BeforeEach(func() {
			appConfig = types.AppConfig{
				Services: types.Services{
					Public: map[string]types.ServiceData{
						"public-service-1": types.ServiceData{},
					},
					Private: map[string]types.ServiceData{
						"private-service-1": types.ServiceData{},
					},
					Worker: map[string]types.ServiceData{
						"worker-service-1": types.ServiceData{},
					},
				},
			}
		})

		It("should return error when the given service already exists", func() {
			err := appConfig.VerifyServiceRoleDoesNotExist("public-service-1")
			Expect(err).To(HaveOccurred())
			err = appConfig.VerifyServiceRoleDoesNotExist("private-service-1")
			Expect(err).To(HaveOccurred())
			err = appConfig.VerifyServiceRoleDoesNotExist("worker-service-1")
			Expect(err).To(HaveOccurred())
		})

		It("should not return an error when the given service does not exist", func() {
			err := appConfig.VerifyServiceRoleDoesNotExist("new-service")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
