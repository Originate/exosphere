package aws_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secrets methods", func() {

	It("merges secrets", func() {
		existingSecrets := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
		})
		newSecrets := types.Secrets(map[string]string{
			"var4": "val4",
		})
		expectedSecrets := types.Secrets(map[string]string{
			"var1": "val1",
			"var2": "val2",
			"var3": "val3",
			"var4": "val4",
		})
		actualSecrets := existingSecrets.Merge(newSecrets)

		Expect(actualSecrets).To(Equal(expectedSecrets))
	})
})
