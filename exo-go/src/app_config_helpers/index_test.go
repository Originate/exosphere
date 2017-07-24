package appConfigHelpers_test

import (
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	appConfig types.AppConfig
	services  types.Services
	appDir    string
	homeDir   string
)

var _ = BeforeSuite(func() {
	dependencies := []types.Dependency{
		types.Dependency{
			Name:    "exocom",
			Version: "0.22.1",
			Silent:  true,
		},
		types.Dependency{
			Name:    "mongo",
			Version: "3.4.0",
			Config:  types.DependencyConfig{DependencyEnvironment: map[string]string{"DB_NAME": "test-db"}},
		},
	}
	todoService := types.ServiceData{Location: "./todo-service", NameSpace: "todo"}
	externalService := types.ServiceData{DockerImage: "originate/test-web-server", Silent: true}
	services = types.Services{
		Private: map[string]types.ServiceData{"todo-service": todoService},
		Public:  map[string]types.ServiceData{"external-service": externalService},
	}
	appConfig = types.AppConfig{
		Name:         "test ",
		Version:      "0.1.0",
		Services:     services,
		Dependencies: dependencies,
	}
	appDir = path.Join("..", "..", "..", "exosphere-shared", "example-apps", "complex-setup-app")
	var err error
	homeDir, err = osHelpers.GetUserHomeDir()
	if err != nil {
		panic(err)
	}
})

var _ = Describe("GetAppConfig", func() {
	var appConfig types.AppConfig

	BeforeEach(func() {
		var err error
		appConfig, err = appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should include name, version and description", func() {
		Expect(appConfig.Name).To(Equal("complex-setup-app"))
		Expect(appConfig.Description).To(Equal("An app with complex setup used for testing"))
		Expect(appConfig.Version).To(Equal("0.0.1"))
	})

	It("should have all the dependencies", func() {
		Expect(appConfig.Dependencies).To(Equal([]types.Dependency{
			types.Dependency{
				Name:    "exocom",
				Version: "0.22.1",
			},
			types.Dependency{
				Name:    "mongo",
				Version: "3.4.0",
				Config: types.DependencyConfig{
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
		privateServices := map[string]types.ServiceData{
			"todo-service":  types.ServiceData{Location: "./todo-service"},
			"users-service": types.ServiceData{NameSpace: "mongo", Location: "./users-service"},
		}
		publicServices := map[string]types.ServiceData{
			"html-server":      types.ServiceData{Location: "./html-server"},
			"external-service": types.ServiceData{DockerImage: "originate/test-web-server"},
		}
		expected := types.Services{Private: privateServices, Public: publicServices}
		Expect(appConfig.Services).To(Equal(expected))
	})

})

var _ = Describe("GetEnvironmentVariables", func() {
	It("should return the environment variables of all dependencies", func() {
		actual := appConfigHelpers.GetEnvironmentVariables(appConfig, appDir, homeDir)
		expected := map[string]string{"EXOCOM_PORT": "80", "DB_NAME": "test-db"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetDependencyNames", func() {
	It("should return the names of all application dependencies", func() {
		actual := appConfigHelpers.GetDependencyNames(appConfig)
		expected := []string{"exocom", "mongo"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetAllDependencyNames", func() {
	It("should return the container names of all application and service dependencies", func() {
		appDir := path.Join("..", "..", "..", "exosphere-shared", "example-apps", "external-dependency")
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		actual, err := appConfigHelpers.GetAllDependencyNames(appDir, appConfig)
		Expect(err).NotTo(HaveOccurred())
		expected := []string{"exocom0.22.1", "mongo3.4.0"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetServiceNames", func() {
	It("should return the names of all services", func() {
		actual := appConfigHelpers.GetServiceNames(services)
		expected := []string{"todo-service", "external-service"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetSilencedDependencyNames", func() {
	It("should return the names of all silenced dependencies", func() {
		actual := appConfigHelpers.GetSilencedDependencyNames(appConfig)
		expected := []string{"exocom"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("GetSilencedServiceNames", func() {
	It("should return the names of all silenced services", func() {
		actual := appConfigHelpers.GetSilencedServiceNames(services)
		expected := []string{"external-service"}
		Expect(actual).To(Equal(expected))
	})
})

var _ = Describe("VerifyServiceDoesNotExist", func() {
	It("should return error when the given service already exists", func() {
		err := appConfigHelpers.VerifyServiceDoesNotExist("todo-service", appConfigHelpers.GetServiceNames(services))
		Expect(err).To(HaveOccurred())
	})

	It("should not return an error when the given service does not exist", func() {
		err := appConfigHelpers.VerifyServiceDoesNotExist("user-service", appConfigHelpers.GetServiceNames(services))
		Expect(err).NotTo(HaveOccurred())
	})
})
