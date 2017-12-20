package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoteDependency", func() {

	Describe("validates required remote fields", func() {
		It("throws an error if remote fields are missing", func() {
			missingConfig := types.RemoteDependency{
				Type: "rds",
				TemplateConfig: map[string]string{
					"engine":               "postgres",
					"engine-version":       "0.0.1",
					"allocated-storage":    "10",
					"username":             "test-user",
					"password-secret-name": "TEST_PASSWORD",
					"instance-class":       "db.t2.micro",
					"storage-type":         "gp2",
				},
			}
			err := missingConfig.ValidateFields()
			Expect(err).To(HaveOccurred())
			expectedErrorString := "remote dependency of type 'rds' missing required field 'template-config.db-name'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("does not throw an error if production fields are valid", func() {
			goodConfig := types.RemoteDependency{
				Type: "rds",
				TemplateConfig: map[string]string{
					"engine":               "postgres",
					"engine-version":       "0.0.1",
					"allocated-storage":    "10",
					"db-name":              "test",
					"username":             "test-user",
					"password-secret-name": "TEST_PASSWORD",
					"instance-class":       "db.t2.micro",
					"storage-type":         "gp2",
				},
			}
			err := goodConfig.ValidateFields()
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
