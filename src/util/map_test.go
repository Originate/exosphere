package util_test

import (
	"github.com/Originate/exosphere/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map methods", func() {
	Describe("Merge", func() {
		It("merges maps", func() {
			existingMap := map[string]string{
				"var1": "val1",
				"var2": "val2",
				"var3": "val3",
			}
			newMap := map[string]string{
				"var4": "val4",
			}
			expectedMap := map[string]string{
				"var1": "val1",
				"var2": "val2",
				"var3": "val3",
				"var4": "val4",
			}
			util.Merge(existingMap, newMap)

			Expect(existingMap).To(Equal(expectedMap))
		})
	})
})
