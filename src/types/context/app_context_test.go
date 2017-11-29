package context_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("AppContext", func() {

	var _ = Describe("GetServiceContexts", func() {
		var appContext *context.AppContext

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("complex-setup-app")
			var err error
			appContext, err = context.GetAppContext(appDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should include all services", func() {
			for _, serviceRole := range appContext.Config.GetSortedServiceRoles() {
				_, exists := appContext.ServiceContexts[serviceRole]
				Expect(exists).To(Equal(true))
			}
		})

		It("should contain correct configuration for the internal service", func() {
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:        "public",
				Description: "dummy html service used for testing setup only - does not run",
				Author:      "test-author",
				ServiceMessages: types.ServiceMessages{
					Sends:    []string{"todo.create"},
					Receives: []string{"todo.created"},
				},
				Development: types.ServiceDevelopmentConfig{
					Scripts: map[string]string{
						"run": `echo "does not run"`,
					},
					Port: "80",
				},
			})
			Expect(err).ToNot(HaveOccurred())
			actual, err := yaml.Marshal(appContext.ServiceContexts["html-server"].Config)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

		It("should contain correct configuration for the external docker image", func() {
			serviceMessages := types.ServiceMessages{
				Sends:    []string{"users.list", "users.create"},
				Receives: []string{"users.listed", "users.created"},
			}
			development := types.ServiceDevelopmentConfig{
				Port: "5000",
				Scripts: map[string]string{
					"run": "node server.js",
				},
			}
			expected, err := yaml.Marshal(types.ServiceConfig{
				Type:            "public",
				Description:     "says hello to the world, ignores .txt files when file watching",
				Author:          "exospheredev",
				ServiceMessages: serviceMessages,
				Development:     development,
			})
			Expect(err).ToNot(HaveOccurred())
			actual, err := yaml.Marshal(appContext.ServiceContexts["external-service"].Config)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})

	})
})
