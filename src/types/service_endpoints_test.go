package types_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceEndpoints", func() {
	var serviceRole string
	var serviceConfig types.ServiceConfig

	var _ = BeforeEach(func() {
		appDir := helpers.GetTestApplicationDir("simple")
		appContext, err := context.GetAppContext(appDir)
		Expect(err).NotTo(HaveOccurred())
		serviceRole = "web"
		serviceConfig = appContext.ServiceContexts[serviceRole].Config
	})

	It("compiles the proper local development endpoints", func() {
		portReservation := types.NewPortReservation()
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentDevelopment,
		}
		s := types.NewServiceEndpoint(serviceRole, serviceConfig, portReservation, buildMode)
		endpoints := s.GetEndpointMappings()
		Expect(endpoints["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
		Expect(endpoints["WEB_ORIGIN"]).To(Equal("http://localhost:8080"))
		mapping := s.GetPortMappings()
		Expect(mapping[0]).To(Equal("3000:8080"))
	})

	It("compiles the proper local production endpoints", func() {
		portReservation := types.NewPortReservation()
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentProduction,
		}
		s := types.NewServiceEndpoint(serviceRole, serviceConfig, portReservation, buildMode)
		endpoints := s.GetEndpointMappings()
		Expect(endpoints["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
		Expect(endpoints["WEB_ORIGIN"]).To(Equal("http://localhost:80"))
		mapping := s.GetPortMappings()
		Expect(mapping[0]).To(Equal("3000:80"))
	})
})
