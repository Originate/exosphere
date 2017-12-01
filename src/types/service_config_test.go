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

	Describe("compiles the correct set of environment variables", func() {
		serviceEnvVars := types.EnvVars{
			Default: map[string]string{
				"env1": "val1",
				"env2": "val2",
			},
			Local: map[string]string{
				"env1": "dev_val1",
				"env3": "dev_val3",
			},
			Remote: map[string]string{
				"env1": "prod_val1",
			},
			Secrets: []string{"secret1", "secret2"},
		}
		serviceConfig := types.ServiceConfig{
			Environment: serviceEnvVars,
		}

		It("compiles local variables", func() {
			expectedVars := map[string]string{
				"env1": "dev_val1",
				"env2": "val2",
				"env3": "dev_val3",
			}
			expectedSecrets := []string{"secret1", "secret2"}
			envVars, secrets := serviceConfig.GetEnvVars("local")
			Expect(expectedVars).To(Equal(envVars))
			Expect(expectedSecrets).To(Equal(secrets))
		})

		It("compiles remote variables", func() {
			expectedVars := map[string]string{
				"env1": "prod_val1",
				"env2": "val2",
			}
			expectedSecrets := []string{"secret1", "secret2"}
			envVars, secrets := serviceConfig.GetEnvVars("remote")
			Expect(expectedVars).To(Equal(envVars))
			Expect(expectedSecrets).To(Equal(secrets))
		})
	})

})
