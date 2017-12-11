package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoteDependency", func() {

	Describe("validates required remote fields", func() {
		It("throws an error if db-name is not valid", func() {
			missingConfig := types.RemoteDependency{
				Type: "rds",
				Config: types.RemoteDependencyConfig{
					Rds: types.RdsConfig{
						Engine:             "postgres",
						EngineVersion:      "0.0.1",
						AllocatedStorage:   "10",
						DbName:             "test!",
						Username:           "test-user",
						PasswordSecretName: "TEST_PASSWORD",
						InstanceClass:      "db.t2.micro",
						StorageType:        "gp2",
						ServiceEnvVarNames: types.ServiceEnvVarNames{
							DbName:   "DB_NAME",
							Username: "DB_USER",
							Password: "DB_PASSWORD",
						},
					},
				},
			}
			err := missingConfig.ValidateFields()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "remote dependency rds has issues: only alphanumeric characters and hyphens allowed in 'rds.db-name'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("throws an error if remote fields are missing", func() {
			missingConfig := types.RemoteDependency{
				Type: "rds",
				Config: types.RemoteDependencyConfig{
					Rds: types.RdsConfig{
						Engine:             "postgres",
						EngineVersion:      "0.0.1",
						AllocatedStorage:   "10",
						Username:           "test-user",
						PasswordSecretName: "TEST_PASSWORD",
						InstanceClass:      "db.t2.micro",
						StorageType:        "gp2",
						ServiceEnvVarNames: types.ServiceEnvVarNames{
							DbName:   "DB_NAME",
							Username: "DB_USER",
							Password: "DB_PASSWORD",
						},
					},
				},
			}
			err := missingConfig.ValidateFields()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "remote dependency rds has issues: missing required field 'rds.db-name'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if production fields are valid", func() {
			goodConfig := types.RemoteDependency{
				Type: "rds",
				Config: types.RemoteDependencyConfig{
					Rds: types.RdsConfig{
						Engine:             "postgres",
						EngineVersion:      "0.0.1",
						AllocatedStorage:   "10",
						DbName:             "test",
						Username:           "test-user",
						PasswordSecretName: "TEST_PASSWORD",
						InstanceClass:      "db.t2.micro",
						StorageType:        "gp2",
						ServiceEnvVarNames: types.ServiceEnvVarNames{
							DbName:   "DB_NAME",
							Username: "DB_USER",
							Password: "DB_PASSWORD",
						},
					},
				},
			}
			err := goodConfig.ValidateFields()
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
