package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceConfig", func() {

	Describe("validates required production fields", func() {
		missingConfig := types.ServiceConfig{}

		publicConfig := types.ServiceConfig{
			Production: map[string]string{
				"url": "originate.com",
			},
		}
		privateConfig := types.ServiceConfig{
			Production: map[string]string{
				"memory": "128",
			},
		}
		workerConfig := types.ServiceConfig{
			Production: map[string]string{
				"cpu":    "128",
				"memory": "128",
			},
		}

		It("throws an error if the production field is missing", func() {
			err := missingConfig.ValidateProductionFields("./missing-service", "public")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "./missing-service/service.yml missing required section 'production'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("throws an error if public production fields are missing", func() {
			err := publicConfig.ValidateProductionFields("./public-service", "public")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "./public-service/service.yml missing required field 'production.cpu'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("throws an error if private production fields are missing", func() {
			err := privateConfig.ValidateProductionFields("./private-service", "private")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "./private-service/service.yml missing required field 'production.cpu'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if no worker production fields are missing", func() {
			err := workerConfig.ValidateProductionFields("./worker-service", "worker")
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
