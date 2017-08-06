package config_test

import (
	"path"

	"github.com/Originate/exosphere/exo-go/src/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppConfig", func() {
	var appConfig config.AppConfig

	var _ = Describe("GetAppConfig", func() {
		BeforeEach(func() {
			appDir := path.Join("..", "..", "..", "example-apps", "complex-setup-app")
			var err error
			appConfig, err = config.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should include name, version and description", func() {
			Expect(appConfig.Name).To(Equal("complex-setup-app"))
			Expect(appConfig.Description).To(Equal("An app with complex setup used for testing"))
			Expect(appConfig.Version).To(Equal("0.0.1"))
		})

		It("should have all the dependencies", func() {
			Expect(appConfig.Dependencies).To(Equal([]config.Dependency{
				config.Dependency{
					Name:    "exocom",
					Version: "0.22.1",
				},
				config.Dependency{
					Name:    "mongo",
					Version: "3.4.0",
					Config: config.DependencyConfig{
						Ports:                 []string{"4000:4000"},
						Volumes:               []string{"{{EXO_DATA_PATH}}:/data/db"},
						OnlineText:            "waiting for connections",
						DependencyEnvironment: map[string]string{"DB_NAME": "test-db"},
						ServiceEnvironment:    map[string]string{"COLLECTION_NAME": "test-collection"},
					},
				},
			}))
		})

		It("should have all the services", func() {
			privateServices := map[string]config.ServiceData{
				"todo-service":  config.ServiceData{Location: "./todo-service"},
				"users-service": config.ServiceData{NameSpace: "mongo", Location: "./users-service"},
			}
			publicServices := map[string]config.ServiceData{
				"html-server":      config.ServiceData{Location: "./html-server"},
				"external-service": config.ServiceData{DockerImage: "originate/test-web-server"},
			}
			expected := config.Services{Private: privateServices, Public: publicServices}
			Expect(appConfig.Services).To(Equal(expected))
		})
	})

	var _ = Describe("GetDependencyNames", func() {
		It("should return the names of all application dependencies", func() {
			appConfig = config.AppConfig{
				Dependencies: []config.Dependency{
					{Name: "exocom"},
					{Name: "mongo"},
				},
			}
			actual := appConfig.GetDependencyNames()
			expected := []string{"exocom", "mongo"}
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetServiceNames", func() {
		It("should return the names of all services", func() {
			appConfig = config.AppConfig{
				Services: config.Services{
					Private: map[string]config.ServiceData{
						"private-service-1": config.ServiceData{},
					},
					Public: map[string]config.ServiceData{
						"public-service-1": config.ServiceData{},
						"public-service-2": config.ServiceData{},
					},
				},
			}
			actual := appConfig.GetServiceNames()
			expected := []string{"private-service-1", "public-service-1", "public-service-2"}
			Expect(actual).To(ConsistOf(expected))
		})
	})

	var _ = Describe("GetSilencedDependencyNames", func() {
		It("should return the names of all silenced dependencies", func() {
			appConfig = config.AppConfig{
				Dependencies: []config.Dependency{
					{Name: "exocom", Silent: true},
					{Name: "mongo"},
				},
			}
			actual := appConfig.GetSilencedDependencyNames()
			expected := []string{"exocom"}
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetSilencedServiceNames", func() {
		It("should return the names of all silenced services", func() {
			appConfig = config.AppConfig{
				Services: config.Services{
					Private: map[string]config.ServiceData{
						"private-service-1": config.ServiceData{Silent: true},
						"private-service-2": config.ServiceData{},
					},
					Public: map[string]config.ServiceData{
						"public-service-1": config.ServiceData{},
						"public-service-2": config.ServiceData{Silent: true},
					},
				},
			}
			actual := appConfig.GetSilencedServiceNames()
			expected := []string{"private-service-1", "public-service-2"}
			Expect(actual).To(ConsistOf(expected))
		})
	})

	var _ = Describe("VerifyServiceDoesNotExist", func() {
		BeforeEach(func() {
			appConfig = config.AppConfig{
				Services: config.Services{
					Public: map[string]config.ServiceData{
						"public-service-1": config.ServiceData{},
					},
					Private: map[string]config.ServiceData{
						"private-service-1": config.ServiceData{},
					},
				},
			}
		})

		It("should return error when the given service already exists", func() {
			err := appConfig.VerifyServiceDoesNotExist("public-service-1")
			Expect(err).To(HaveOccurred())
			err = appConfig.VerifyServiceDoesNotExist("private-service-1")
			Expect(err).To(HaveOccurred())
		})

		It("should not return an error when the given service does not exist", func() {
			err := appConfig.VerifyServiceDoesNotExist("new-service")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
