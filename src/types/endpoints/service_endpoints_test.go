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
		serviceEndpoints := endpoints.NewServiceEndpoints(appContext, buildMode)
		envVars := serviceEndpoints.GetServiceEndpointEnvVars("users")
		Expect(envVars["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
		Expect(envVars["WEB_ORIGIN"]).To(Equal("http://localhost:8080"))
		mapping := serviceEndpoints.GetServicePortMappings("web")
		Expect(mapping[0]).To(Equal("3000:4000"))
	})

	It("compiles the proper local production endpoints", func() {
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentProduction,
		}
		serviceEndpoints := endpoints.NewServiceEndpoints(appContext, buildMode)
		envVars := serviceEndpoints.GetServiceEndpointEnvVars("users")
		Expect(envVars["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
		Expect(envVars["WEB_ORIGIN"]).To(Equal("http://localhost:80"))
		mapping := serviceEndpoints.GetServicePortMappings("web")
		Expect(mapping[0]).To(Equal("3000:80"))
	})
})
