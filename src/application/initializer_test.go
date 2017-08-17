package application_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initializer", func() {
	It("should create a docker.compose.yml", func() {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = testHelpers.CheckoutApp(cwd, "complex-setup-app")
		Expect(err).NotTo(HaveOccurred())
		internalServices := []string{"html-server", "todo-service", "users-service"}
		externalServices := []string{"external-service"}
		internalDependencies := []string{"exocom0.24.0"}
		externalDependencies := []string{"mongo3.4.0"}
		allServices := util.JoinStringSlices(internalServices, externalServices, internalDependencies, externalDependencies)

		appDir := path.Join("tmp", "complex-setup-app")
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		appConfig, err := types.NewAppConfig(appDir)
		Expect(err).NotTo(HaveOccurred())
		mockLogger := application.NewLogger([]string{}, []string{}, ioutil.Discard)
		initializer, err := application.NewInitializer(appConfig, mockLogger.GetLogChannel(""), "exo-run", appDir, homeDir)
		Expect(err).NotTo(HaveOccurred())
		err = initializer.Initialize()
		Expect(err).NotTo(HaveOccurred())
		expectedDockerComposePath := path.Join(appDir, "tmp", "docker-compose.yml")
		Expect(util.DoesFileExist(expectedDockerComposePath)).To(Equal(true))
		dockerCompose, err := docker.GetDockerCompose(expectedDockerComposePath)
		Expect(err).NotTo(HaveOccurred())

		By("set the version to 3")
		Expect(dockerCompose.Version).To(Equal("3"))

		By("list all services and dependencies under 'services'")
		for _, serviceName := range allServices {
			_, exists := dockerCompose.Services[serviceName]
			Expect(exists).To(Equal(true))
		}

		By("generate an image name for each dependency and external service")
		for _, serviceName := range util.JoinStringSlices(internalDependencies, externalDependencies, externalServices) {
			Expect(len(dockerCompose.Services[serviceName].Image)).ToNot(Equal(0))
		}

		By("should generate a container name for each service and dependency")
		for _, serviceName := range allServices {
			Expect(len(dockerCompose.Services[serviceName].ContainerName)).ToNot(Equal(0))
		}

		By("should have the correct build command for each internal service and dependency")
		for _, serviceName := range internalServices {
			Expect(dockerCompose.Services[serviceName].Command).To(Equal(`echo "does not run"`))
		}
		Expect(dockerCompose.Services["exocom0.24.0"].Command).To(Equal(""))

		By("should include 'exocom' in the dependencies of every service")
		for _, serviceName := range append(internalServices, externalServices...) {
			exists := util.DoesStringArrayContain(dockerCompose.Services[serviceName].DependsOn, "exocom0.24.0")
			Expect(exists).To(Equal(true))
		}

		By("should include external dependencies as dependencies")
		exists := util.DoesStringArrayContain(dockerCompose.Services["todo-service"].DependsOn, "mongo3.4.0")
		Expect(exists).To(Equal(true))

		By("should include the correct exocom environment variables")
		environment := dockerCompose.Services["exocom0.24.0"].Environment
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

		By("should include exocom environment variables in internal services' environment")
		for _, serviceName := range internalServices {
			environment := dockerCompose.Services[serviceName].Environment
			Expect(environment["EXOCOM_HOST"]).To(Equal("exocom0.24.0"))
			Expect(environment["EXOCOM_PORT"]).To(Equal("$EXOCOM_PORT"))
		}

		By("should generate a volume path for an external dependency that mounts a volume")
		Expect(len(dockerCompose.Services["mongo3.4.0"].Volumes)).NotTo(Equal(0))

		By("should have the specified image and container names for the external service")
		serviceName := "external-service"
		imageName := "originate/test-web-server"
		Expect(dockerCompose.Services[serviceName].Image).To(Equal(imageName))
		Expect(dockerCompose.Services[serviceName].ContainerName).To(Equal(serviceName))

		By("should have the specified ports, volumes and environment variables for the external service")
		serviceName = "external-service"
		ports := []string{"5000:5000"}
		Expect(dockerCompose.Services[serviceName].Ports).To(Equal(ports))
		Expect(len(dockerCompose.Services[serviceName].Volumes)).NotTo(Equal(0))
		Expect(dockerCompose.Services[serviceName].Environment["EXTERNAL_SERVICE_HOST"]).To(Equal("external-service0.1.2"))
		Expect(dockerCompose.Services[serviceName].Environment["EXTERNAL_SERVICE_PORT"]).To(Equal("$EXTERNAL_SERVICE_PORT"))

		By("should have the ports for the external dependency defined in application.yml")
		serviceName = "mongo3.4.0"
		ports = []string{"4000:4000"}
		Expect(dockerCompose.Services[serviceName].Ports).To(Equal(ports))
		Expect(len(dockerCompose.Services[serviceName].Volumes)).NotTo(Equal(0))
	})
})
