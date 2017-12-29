package endpoints_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/endpoints"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceEndpoints", func() {
	var appContext *context.AppContext

	var _ = BeforeEach(func() {
		var err error
		appDir := helpers.GetTestApplicationDir("running")
		appContext, err = context.GetAppContext(appDir)
		Expect(err).NotTo(HaveOccurred())
	})

	It("compiles the proper local development endpoints", func() {
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentDevelopment,
		}
		serviceEndpoints := endpoints.NewServiceEndpoints(appContext, buildMode, "")
		envVars := serviceEndpoints.GetServiceEndpointEnvVars("users")
		Expect(envVars["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3010"))
		Expect(envVars["WEB_INTERNAL_ORIGIN"]).To(Equal("http://web:4000"))
		mapping := serviceEndpoints.GetServicePortMappings("web")
		Expect(mapping[0]).To(Equal("3010:4000"))

		envVars = serviceEndpoints.GetServiceEndpointEnvVars("web")
		Expect(envVars["USERS_HOST"]).To(Equal("http://users:2121"))
		mapping = serviceEndpoints.GetServicePortMappings("users")
		Expect(mapping[0]).To(Equal("3000:2121"))
	})

	It("compiles the proper local production endpoints", func() {
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentProduction,
		}
		serviceEndpoints := endpoints.NewServiceEndpoints(appContext, buildMode, "qa")
		envVars := serviceEndpoints.GetServiceEndpointEnvVars("users")
		Expect(envVars["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3010"))
		Expect(envVars["WEB_INTERNAL_ORIGIN"]).To(Equal("http://web:80"))
		mapping := serviceEndpoints.GetServicePortMappings("web")
		Expect(mapping[0]).To(Equal("3010:80"))

		envVars = serviceEndpoints.GetServiceEndpointEnvVars("web")
		Expect(envVars["USERS_HOST"]).To(Equal("http://users:21"))
		mapping = serviceEndpoints.GetServicePortMappings("users")
		Expect(mapping[0]).To(Equal("3000:21"))
	})
})
