package awsHelper_test

import (
	"github.com/Originate/exosphere/exo-go/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secrets methods", func() {

	It("converts a tf string to a Secrets map", func() {
		tfvars := `var1="val1"
var2="val2"
var3="val3"`
		expectedMap := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
		})
		Expect(types.NewSecrets(tfvars)).To(Equal(expectedMap))
	})

	It("converts Secrets to a tf string", func() {
		secrets := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
		})
		expectedTfvars := `var1="val1"
var2="val2"
var3="val3"`
		Expect(secrets.TfString()).To(Equal(expectedTfvars))
	})

	It("merges secrets", func() {
		existingTfVars := `var1="val1"
var2="val2"
var3="val3"`
		newSecrets := types.Secrets(map[string]string{
			"var4": "val4",
		})

		expectedTfString := `var1="val1"
var2="val2"
var3="val3"
var4="val4"`
		actualSecrets := types.NewSecrets(existingTfVars).Merge(newSecrets)

		Expect(expectedTfString).To(Equal(actualSecrets.TfString()))
	})
})
