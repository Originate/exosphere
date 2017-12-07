package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceConfig", func() {

	Describe("validates service.yml fields", func() {
		wrongType := types.ServiceConfig{Type: "wrong-type"}
		rightType := types.ServiceConfig{Type: "public"}

		It("throws an error if the service type is unsupported", func() {
			err := wrongType.ValidateServiceConfig()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "Invalid value 'wrong-type' in service.yml field 'type'. Must be one of: public, worker"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if the service type is valid", func() {
			err := rightType.ValidateServiceConfig()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("validates required production fields", func() {
		publicConfig := types.ServiceConfig{
			Remote: types.ServiceRemoteConfig{
				URL: "originate.com",
			},
		}
		workerConfig := types.ServiceConfig{
			Remote: types.ServiceRemoteConfig{
				CPU:    "128",
				Memory: "128",
			},
		}

		It("throws an error if public deployment fields are missing", func() {
			err := publicConfig.ValidateDeployFields("./public-service", "public")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "./public-service/service.yml missing required field 'production.Port'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if no worker production fields are missing", func() {
			err := workerConfig.ValidateDeployFields("./worker-service", "worker")
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
