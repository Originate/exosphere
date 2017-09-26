package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductionDependencyConfig", func() {

	Describe("validates required production fields", func() {
		missingConfig := types.ProductionDependencyConfig{
			Name:    "postgres",
			Version: "0.0.1",
			Config: types.ProductionDependencyConfigOptions{
				Rds: types.RdsConfig{
					Username:       "test-user",
					PasswordEnvVar: "TEST_PASSWORD",
				},
			},
		}
		goodConfig := types.ProductionDependencyConfig{
			Name:    "postgres",
			Version: "0.0.1",
			Config: types.ProductionDependencyConfigOptions{
				Rds: types.RdsConfig{
					DbName:         "test",
					Username:       "test-user",
					PasswordEnvVar: "TEST_PASSWORD",
				},
			},
		}
		It("throws an error if production fields are missing", func() {
			err := missingConfig.ValidateFields()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "production dependency postgres:0.0.1 missing field(s): missing required field 'Rds.DbName'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if production fields are not missing", func() {
			err := goodConfig.ValidateFields()
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
