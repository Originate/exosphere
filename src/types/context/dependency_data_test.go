package context_test

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDependencyData", func() {
	It("should include data from all services", func() {
		appContext := &context.AppContext{
			Config: &types.AppConfig{
				Services: map[string]types.ServiceSource{
					"serviceB": {},
					"serviceA": {},
				},
			},
			ServiceContexts: map[string]*context.ServiceContext{
				"serviceA": &context.ServiceContext{
					Config: types.ServiceConfig{
						DependencyData: types.ServiceDependencyData{
							"exocom": {
								"key3": "value3",
							},
						},
					},
					Source: &types.ServiceSource{
						DependencyData: types.ServiceDependencyData{
							"exocom": {
								"key1": []interface{}{"a", "b"},
							},
						},
					},
				},
				"serviceB": &context.ServiceContext{
					Config: types.ServiceConfig{
						DependencyData: types.ServiceDependencyData{
							"exocom": {
								"key4": "value4",
							},
						},
					},
					Source: &types.ServiceSource{
						DependencyData: types.ServiceDependencyData{
							"exocom": {
								"key2": map[string]interface{}{
									"c": "d",
									"e": "f",
								},
							},
						},
					},
				},
			},
		}
		result := appContext.GetDependencyServiceData("exocom")
		Expect(result).To(Equal(map[string]map[string]interface{}{
			"serviceA": map[string]interface{}{
				"key1": []interface{}{"a", "b"},
				"key3": "value3",
			},
			"serviceB": map[string]interface{}{
				"key2": map[string]interface{}{
					"c": "d",
					"e": "f",
				},
				"key4": "value4",
			},
		}))
	})
})
