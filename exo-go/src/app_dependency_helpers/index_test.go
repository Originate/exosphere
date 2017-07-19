package appDependencyHelpers_test

import (
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppDependency", func() {

	var exocom appDependencyHelpers.AppDependency
	var mongo appDependencyHelpers.AppDependency
	var nats appDependencyHelpers.AppDependency

	var _ = BeforeSuite(func() {
		appConfig := types.AppConfig{}
		exocom = appDependencyHelpers.Build(types.Dependency{
			Name:    "exocom",
			Version: "0.22.1",
		}, appConfig)
		nats = appDependencyHelpers.Build(types.Dependency{
			Name:    "nats",
			Version: "0.9.6",
		}, appConfig)
		mongo = appDependencyHelpers.Build(types.Dependency{
			Name:    "mongo",
			Version: "3.4.0",
			Config: types.DependencyConfig{
				OnlineText:            "MongoDB connected",
				DependencyEnvironment: map[string]string{"DB_NAME": "test-db"},
			},
		}, appConfig)
	})

	var _ = Describe("GetContainerName", func() {
		It("should be the concatenation of dependency name and version", func() {
			Expect(exocom.GetContainerName()).To(Equal("exocom0.22.1"))
			Expect(mongo.GetContainerName()).To(Equal("mongo3.4.0"))
			Expect(nats.GetContainerName()).To(Equal("nats0.9.6"))
		})
	})

	var _ = Describe("GetEnviromentVariables", func() {
		It("should set a default port for EXOCOM_PORT if it is not set", func() {
			expected := map[string]string{"EXOCOM_PORT": "80"}
			Expect(exocom.GetEnvVariables()).To(Equal(expected))
		})

		It("should return the EXOCOM_PORT as is set on the user's machine", func() {
			if err := os.Setenv("EXOCOM_PORT", "5000"); err != nil {
				panic(err)
			}
			expected := map[string]string{"EXOCOM_PORT": "5000"}
			Expect(exocom.GetEnvVariables()).To(Equal(expected))
		})

		It("should include the correct NATS_HOST for nats dependency", func() {
			expected := map[string]string{"NATS_HOST": "nats0.9.6"}
			Expect(nats.GetEnvVariables()).To(Equal(expected))
		})

		It("should return the correct environment variables for external dependencies", func() {
			expected := map[string]string{"DB_NAME": "test-db"}
			Expect(mongo.GetEnvVariables()).To(Equal(expected))
		})
	})

	var _ = Describe("GetOnlineText", func() {
		It("should return the correct online text for exocom", func() {
			expected := "ExoCom WebSocket listener online"
			Expect(exocom.GetOnlineText()).To(Equal(expected))
		})

		It("should return the correct online text for nats", func() {
			expected := "Listening for route connections"
			Expect(nats.GetOnlineText()).To(Equal(expected))
		})

		It("should include the correct online text for external dependencies", func() {
			expected := "MongoDB connected"
			Expect(mongo.GetOnlineText()).To(Equal(expected))
		})
	})

})
