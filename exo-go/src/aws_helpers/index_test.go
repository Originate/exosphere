package awsHelper_test

import (
	"github.com/Originate/exosphere/exo-go/src/aws_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AwsHelpers", func() {

	var _ = Describe("TFString methods", func() {
		It("converts TFString to a Secrets map", func() {
			tfvars := types.TFString(`var1="val1"
var2="val2"
var3="val3"`)
			expectedMap := types.Secrets(map[string]string{
				"var1": "val1",
				"var2": "val2",
				"var3": "val3",
			})
			Expect(tfvars.ToMap()).To(Equal(expectedMap))
		})
	})

	var _ = Describe("Secrets methods", func() {
		It("converts Secrets to a TFString", func() {
			secrets := types.Secrets(map[string]string{
				"var1": "val1",
				"var2": "val2",
				"var3": "val3",
			})
			expectedTfvars := types.TFString(`var1="val1"
var2="val2"
var3="val3"`)
			Expect(secrets.ToTfString()).To(Equal(expectedTfvars))
		})
	})

	var _ = Describe("validating and merging secrets", func() {
		It("throws an error if existing key is created", func() {
			existingTfVars := types.TFString(`var1="val1"
var2="val2"
var3="val3"`)
			newSecrets := types.Secrets(map[string]string{
				"var1": "val1",
			})

			_, err := awsHelper.ValidateAndMergeSecrets(existingTfVars, newSecrets)
			Expect(err).To(HaveOccurred())
		})

		It("merges secrets if there are no conflicting keys", func() {
			existingTfVars := types.TFString(`var1="val1"
var2="val2"
var3="val3"`)
			newSecrets := types.Secrets(map[string]string{
				"var4": "val4",
			})

			expectedTfString := types.TFString(`var1="val1"
var2="val2"
var3="val3"
var4="val4"`)
			actualTfString, err := awsHelper.ValidateAndMergeSecrets(existingTfVars, newSecrets)

			Expect(err).NotTo(HaveOccurred())
			Expect(expectedTfString).To(Equal(actualTfString))
		})
	})
})
