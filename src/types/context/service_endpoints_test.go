package context_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceEndpoints", func() {
	var appContext *context.AppContext
	var serviceRole string

	var _ = BeforeEach(func() {
		appDir := helpers.GetTestApplicationDir("simple")
		var err error
		appContext, err = context.GetAppContext(appDir)
		Expect(err).NotTo(HaveOccurred())
		serviceRole = "web"
	})

	It("compiles the proper local development endpoints", func() {
		portReservation := types.NewPortReservation()
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentDevelopment,
		}
		s := context.NewServiceEndpoint(appContext, serviceRole, portReservation, buildMode)
		endpoints := s.GetEndpointMappings()
		Expect(endpoints["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
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
		s := context.NewServiceEndpoint(appContext, serviceRole, portReservation, buildMode)
		endpoints := s.GetEndpointMappings()
		Expect(endpoints["WEB_EXTERNAL_ORIGIN"]).To(Equal("http://localhost:3000"))
		mapping := s.GetPortMappings()
		Expect(mapping[0]).To(Equal("3000:80"))
	})
})
