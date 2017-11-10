package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceConfig", func() {

	Describe("validates required production fields", func() {
		publicConfig := types.ServiceConfig{
			Production: types.ServiceProductionConfig{
				URL: "originate.com",
			},
		}
		workerConfig := types.ServiceConfig{
			Production: types.ServiceProductionConfig{
				CPU:    "128",
				Memory: "128",
			},
		}

		It("throws an error if public production fields are missing", func() {
			err := publicConfig.Production.ValidateFields("./public-service", "public")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "./public-service/service.yml missing required field 'production.CPU'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if no worker production fields are missing", func() {
			err := workerConfig.Production.ValidateFields("./worker-service", "worker")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("compiles the correct set of environment variables", func() {
		serviceEnvVars := types.EnvVars{
			Default: map[string]string{
				"env1": "val1",
				"env2": "val2",
			},
			Development: map[string]string{
				"env1": "dev_val1",
				"env3": "dev_val3",
			},
			Production: map[string]string{
				"env1": "prod_val1",
			},
			Secrets: []string{"secret1", "secret2"},
		}
		serviceConfig := types.ServiceConfig{
			Environment: serviceEnvVars,
		}

		It("compiles development variables", func() {
			expectedVars := map[string]string{
				"env1": "dev_val1",
				"env2": "val2",
				"env3": "dev_val3",
			}
			expectedSecrets := []string{"secret1", "secret2"}
			envVars, secrets := serviceConfig.GetEnvVars("development")
			Expect(expectedVars).To(Equal(envVars))
			Expect(expectedSecrets).To(Equal(secrets))
		})

		It("compiles production variables", func() {
			expectedVars := map[string]string{
				"env1": "prod_val1",
				"env2": "val2",
			}
			expectedSecrets := []string{"secret1", "secret2"}
			envVars, secrets := serviceConfig.GetEnvVars("production")
			Expect(expectedVars).To(Equal(envVars))
			Expect(expectedSecrets).To(Equal(secrets))
		})
	})

})
