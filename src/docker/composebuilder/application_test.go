package composebuilder_test

import (
	"io/ioutil"
	"strings"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("composebuilder", func() {
	var _ = Describe("GetApplicationDockerCompose", func() {
		It("should return the proper docker configs for deployment", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "rds")
			Expect(err).NotTo(HaveOccurred())
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())
			dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
				AppConfig: appConfig,
				AppDir:    appDir,
				BuildMode: composebuilder.BuildMode{
					Type: composebuilder.BuildModeTypeDeploy,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("ignoring rds dependencies")
			delete(dockerCompose.Services, "my-sql-service")
			Expect(len(dockerCompose.Services)).To(Equal(0))
		})

		It("should return the proper docker configs for development", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "complex-setup-app")
			Expect(err).NotTo(HaveOccurred())
			internalServices := []string{"html-server", "todo-service", "users-service"}
			externalServices := []string{"external-service"}
			internalDependencies := []string{"exocom0.26.1"}
			externalDependencies := []string{"mongo3.4.0"}
			allServices := util.JoinStringSlices(internalServices, externalServices, internalDependencies, externalDependencies)
			appConfig, err := types.NewAppConfig(appDir)
			Expect(err).NotTo(HaveOccurred())

			dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
				AppConfig: appConfig,
				AppDir:    appDir,
				BuildMode: composebuilder.BuildMode{
					Type:        composebuilder.BuildModeTypeLocal,
					Environment: composebuilder.BuildModeEnvironmentDevelopment,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			By("generate an image name for each dependency and external service")
			for _, serviceRole := range util.JoinStringSlices(internalDependencies, externalDependencies, externalServices) {
				Expect(len(dockerCompose.Services[serviceRole].Image)).ToNot(Equal(0))
			}

			By("should generate a container name for each service and dependency")
			for _, serviceRole := range allServices {
				Expect(len(dockerCompose.Services[serviceRole].ContainerName)).ToNot(Equal(0))
			}

			By("should have the correct build command for each internal service and dependency")
			for _, serviceRole := range internalServices {
				Expect(dockerCompose.Services[serviceRole].Command).To(Equal(`echo "does not run"`))
			}
			Expect(dockerCompose.Services["exocom0.26.1"].Command).To(Equal(""))

			By("should include 'exocom' in the dependencies of every service")
			for _, serviceRole := range append(internalServices, externalServices...) {
				exists := util.DoesStringArrayContain(dockerCompose.Services[serviceRole].DependsOn, "exocom0.26.1")
				Expect(exists).To(Equal(true))
			}

			By("should include external dependencies as dependencies")
			exists := util.DoesStringArrayContain(dockerCompose.Services["todo-service"].DependsOn, "mongo3.4.0")
			Expect(exists).To(Equal(true))

			By("should properly reserve ports for services")
			actualApiPort := dockerCompose.Services["api-service"].Ports
			expectedApiPort := []string{"3000:80"}
			Expect(actualApiPort).To(Equal(expectedApiPort))

			actualExternalServicePort := dockerCompose.Services["external-service"].Ports
			expectedExternalServicePort := []string{"3010:5000"}
			Expect(actualExternalServicePort).To(Equal(expectedExternalServicePort))

			actualHtmlPort := dockerCompose.Services["html-server"].Ports
			expectedHtmlPort := []string{"3020:80"}
			Expect(actualHtmlPort).To(Equal(expectedHtmlPort))

			By("should inject the proper service endpoint environment variables")
			expectedApiEndpointKey := "API_SERVICE_EXTERNAL_ORIGIN"
			expectedApiEndpointValue := "http://localhost:3000"
			expectedHtmlEndpointKey := "HTML_SERVER_EXTERNAL_ORIGIN"
			expectedHtmlEndpointValue := "http://localhost:3020"

			skipServices := []string{"api-service", "exocom0.26.1", "mongo3.4.0"}
			for serviceRole, dockerConfig := range dockerCompose.Services {
				if util.DoesStringArrayContain(skipServices, serviceRole) {
					continue
				}
				Expect(dockerConfig.Environment[expectedApiEndpointKey]).To(Equal(expectedApiEndpointValue))
			}
			skipServices = []string{"html-server", "exocom0.26.1", "mongo3.4.0"}
			for serviceRole, dockerConfig := range dockerCompose.Services {
				if util.DoesStringArrayContain(skipServices, serviceRole) {
					continue
				}
				Expect(dockerConfig.Environment[expectedHtmlEndpointKey]).To(Equal(expectedHtmlEndpointValue))
			}
			nonPublicServiceKeys := []string{
				"USERS_SERVICE_EXTERNAL_ORIGIN",
				"TODO_SERVICE_EXTERNAL_ORIGIN",
			}
			for _, dockerConfig := range dockerCompose.Services {
				for _, nonPublicKey := range nonPublicServiceKeys {
					Expect(dockerConfig.Environment[nonPublicKey]).To(Equal(""))
				}
			}

			By("should include the correct exocom environment variables")
			environment := dockerCompose.Services["exocom0.26.1"].Environment
			expectedServiceRoutes := []string{
				`{"receives":["todo.create"],"role":"todo-service","sends":["todo.created"]}`,
				`{"namespace":"mongo","receives":["mongo.list","mongo.create"],"role":"users-service","sends":["mongo.listed","mongo.created"]}`,
				`{"receives":["todo.created"],"role":"html-server","sends":["todo.create"]}`,
				`{"receives":["users.listed","users.created"],"role":"external-service","sends":["users.list","users.create"]}`,
			}
			for _, serviceRoute := range expectedServiceRoutes {
				Expect(strings.Contains(environment["SERVICE_ROUTES"], serviceRoute))
			}

			By("should include exocom environment variables in every services' environment")
			for _, serviceRole := range append(internalServices, externalServices...) {
				environment := dockerCompose.Services[serviceRole].Environment
				Expect(environment["EXOCOM_HOST"]).To(Equal("exocom0.26.1"))
			}

			By("should generate a volume path for an external dependency that mounts a volume")
			Expect(len(dockerCompose.Services["mongo3.4.0"].Volumes)).NotTo(Equal(0))

			By("should have the specified image and container names for the external service")
			serviceRole := "external-service"
			imageName := "originate/test-web-server:0.0.1"
			Expect(dockerCompose.Services[serviceRole].Image).To(Equal(imageName))
			Expect(dockerCompose.Services[serviceRole].ContainerName).To(Equal(serviceRole))

			By("should have the ports for the external dependency defined in application.yml")
			serviceRole = "mongo3.4.0"
			ports := []string{"4000:4000"}
			Expect(dockerCompose.Services[serviceRole].Ports).To(Equal(ports))
			Expect(len(dockerCompose.Services[serviceRole].Volumes)).NotTo(Equal(0))
		})
	})
})
