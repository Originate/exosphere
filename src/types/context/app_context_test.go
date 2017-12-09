package context_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			expected := types.ServiceConfig{
				Type:        "public",
				Description: "dummy html service used for testing setup only - does not run",
				Author:      "test-author",
				DependencyData: types.ServiceDependencyData{
					"exocom": {
						"receives": []interface{}{"todo.created"},
						"sends":    []interface{}{"todo.create"},
					},
				},
				Development: types.ServiceDevelopmentConfig{
					Scripts: map[string]string{
						"run": `echo "does not run"`,
					},
					Port: "80",
				},
			}
			actual := appContext.ServiceContexts["html-server"].Config
			Expect(actual).To(Equal(expected))
		})

		It("should contain correct configuration for the external docker image", func() {
			development := types.ServiceDevelopmentConfig{
				Port: "5000",
				Scripts: map[string]string{
					"run": "node server.js",
				},
			}
			expected := types.ServiceConfig{
				Type:        "public",
				Description: "says hello to the world, ignores .txt files when file watching",
				Author:      "exospheredev",
				Development: development,
			}
			actual := appContext.ServiceContexts["external-service"].Config
			Expect(actual).To(Equal(expected))
		})
	})

	var _ = Describe("GetDependencyData", func() {
		var appContext *context.AppContext

		var _ = BeforeEach(func() {
			appDir := helpers.GetTestApplicationDir("complex-setup-app")
			var err error
			appContext, err = context.GetAppContext(appDir)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should correctly merge all the dependency data", func() {
			expected := map[string]map[string]interface{}{
				"api-service":      {},
				"external-service": {},
				"html-server": {
					"receives": []interface{}{"todo.created"},
					"sends":    []interface{}{"todo.create"},
				},
				"todo-service": {
					"receives": []interface{}{"todo.create"},
					"sends":    []interface{}{"todo.created"},
				},
				"users-service": {
					"receives": []interface{}{"mongo.list", "mongo.create"},
					"sends":    []interface{}{"mongo.listed", "mongo.created"},
					"translations": []interface{}{
						map[string]interface{}{
							"internal": "mongo create",
							"public":   "users create",
						},
					},
				},
			}
			actual := appContext.GetDependencyServiceData("exocom")
			Expect(actual).To(Equal(expected))
		})
	})
})
