package appSetup_test

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_setup"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/Originate/exosphere/exo-go/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setup", func() {

	var dockerCompose types.DockerCompose
	var internalServices, externalServices, internalDependencies, externalDependencies, allServices []string

	var _ = BeforeSuite(func() {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = testHelpers.CheckoutApp(cwd, "complex-setup-app")
		Expect(err).NotTo(HaveOccurred())
		internalServices = []string{"html-server", "todo-service", "users-service"}
		externalServices = []string{"external-service"}
		internalDependencies = []string{"exocom0.22.1"}
		externalDependencies = []string{"mongo3.4.0"}
		allServices = util.JoinStringSlices(internalServices, externalServices, internalDependencies, externalDependencies)

		appDir := path.Join("tmp", "complex-setup-app")
		homeDir, err := osHelpers.GetUserHomeDir()
		if err != nil {
			panic(err)
		}
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		_, pipeWriter := io.Pipe()
		mockLogger := logger.NewLogger([]string{}, []string{}, pipeWriter)
		setup, err := appSetup.NewAppSetup(appConfig, mockLogger, appDir, homeDir)
		Expect(err).NotTo(HaveOccurred())
		err = setup.Setup()
		Expect(err).NotTo(HaveOccurred())
		expectedDockerComposePath := path.Join(appDir, "tmp", "docker-compose.yml")
		Expect(osHelpers.FileExists(expectedDockerComposePath)).To(Equal(true))
		dockerCompose, err = dockerHelpers.GetDockerCompose(expectedDockerComposePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should create docker-compose.yml version 3", func() {
		Expect(dockerCompose.Version).To(Equal("3"))
	})

	It("should list all services and dependencies under 'services'", func() {
		for _, serviceName := range allServices {
			_, exists := dockerCompose.Services[serviceName]
			Expect(exists).To(Equal(true))
		}
	})

	It("should generate an image name for each dependency and external service", func() {
		for _, serviceName := range util.JoinStringSlices(internalDependencies, externalDependencies, externalServices) {
			Expect(len(dockerCompose.Services[serviceName].Image)).ToNot(Equal(0))
		}
	})

	It("should generate a container name for each service and dependency", func() {
		for _, serviceName := range allServices {
			Expect(len(dockerCompose.Services[serviceName].ContainerName)).ToNot(Equal(0))
		}
	})

	It("should have the correct build command for each internal service and dependency", func() {
		for _, serviceName := range internalServices {
			Expect(dockerCompose.Services[serviceName].Command).To(Equal(`echo "does not run"`))
		}
		for _, serviceName := range internalDependencies {
			dependencyName := string(regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)$`).ReplaceAll([]byte(serviceName), []byte("")))
			Expect(dockerCompose.Services[serviceName].Command).To(Equal(fmt.Sprintf("bin/%s", dependencyName)))
		}
	})

	It("should include 'exocom' in the dependencies of every service", func() {
		for _, serviceName := range append(internalServices, externalServices...) {
			exists := util.DoesStringArrayContain(dockerCompose.Services[serviceName].DependsOn, "exocom0.22.1")
			Expect(exists).To(Equal(true))
		}
	})

	It("should include external dependencies as dependencies", func() {
		exists := util.DoesStringArrayContain(dockerCompose.Services["todo-service"].DependsOn, "mongo3.4.0")
		Expect(exists).To(Equal(true))
	})

	It("should include the correct exocom environment variables", func() {
		environment := dockerCompose.Services["exocom0.22.1"].Environment
		Expect(environment["PORT"]).To(Equal("$EXOCOM_PORT"))
		expectedServiceRoutes := []string{
			`{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]}`,
			`{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}`,
			`{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}`,
			`{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]}`,
		}
		for _, serviceRoute := range expectedServiceRoutes {
			Expect(strings.Contains(environment["SERVICE_ROUTES"], serviceRoute))
		}
	})

	It("should include exocom environment variables in internal services' environment", func() {
		for _, serviceName := range internalServices {
			environment := dockerCompose.Services[serviceName].Environment
			Expect(environment["EXOCOM_HOST"]).To(Equal("exocom0.22.1"))
			Expect(environment["EXOCOM_PORT"]).To(Equal("$EXOCOM_PORT"))
		}
	})

	It("should generate a volume path for an external dependency that mounts a volume", func() {
		Expect(len(dockerCompose.Services["mongo3.4.0"].Volumes)).NotTo(Equal(0))
	})

	It("should have the specified image and container names for the external service", func() {
		serviceName := "external-service"
		imageName := "originate/test-web-server"
		Expect(dockerCompose.Services[serviceName].Image).To(Equal(imageName))
		Expect(dockerCompose.Services[serviceName].ContainerName).To(Equal(serviceName))
	})

	It("should have the specified ports, volumes and environment variables for the external service", func() {
		serviceName := "external-service"
		ports := []string{"5000:5000"}
		Expect(dockerCompose.Services[serviceName].Ports).To(Equal(ports))
		Expect(len(dockerCompose.Services[serviceName].Volumes)).NotTo(Equal(0))
		Expect(dockerCompose.Services[serviceName].Environment["EXTERNAL_SERVICE_HOST"]).To(Equal("external-service0.1.2"))
		Expect(dockerCompose.Services[serviceName].Environment["EXTERNAL_SERVICE_PORT"]).To(Equal("$EXTERNAL_SERVICE_PORT"))
	})

	It("should have the ports for the external dependency defined in application.yml", func() {
		serviceName := "mongo3.4.0"
		ports := []string{"4000:4000"}
		Expect(dockerCompose.Services[serviceName].Ports).To(Equal(ports))
		Expect(len(dockerCompose.Services[serviceName].Volumes)).NotTo(Equal(0))
	})

})
