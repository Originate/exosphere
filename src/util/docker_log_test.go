package util_test

import (
	"github.com/Originate/exosphere/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NormalizeDockerComposeLog", func() {
	It("should remove colors", func() {
		line := "\u001b[33mweb             |\u001b[0m web server running at port 4000"
		result := util.NormalizeDockerComposeLog(line)
		Expect(result).To(Equal("web             | web server running at port 4000"))
	})

	It("should remove VT100 control characters", func() {
		line := "\u001b[1A\u001b[2KCreating exoservice ..."
		result := util.NormalizeDockerComposeLog(line)
		Expect(result).To(Equal("Creating exoservice ..."))
	})

	It("should transform carriage return character to newlines", func() {
		line := "Creating exoservice ...\rdone"
		result := util.NormalizeDockerComposeLog(line)
		Expect(result).To(Equal("Creating exoservice ...\ndone"))
	})
})

var _ = Describe("ParseDockerComposeLog", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.26.1, web"
		serviceRole, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceRole).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

	It("should parse service log message correctly", func() {
		line := "web             | web server running at port 4000"
		serviceRole, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceRole).To(Equal("web"))
		Expect(serviceOutput).To(Equal("web server running at port 4000"))
	})

	It("should strip version from service name", func() {
		line := "exocom0.26.1    | ExoCom HTTP service online at port 80"
		serviceRole, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceRole).To(Equal("exocom"))
		Expect(serviceOutput).To(Equal("ExoCom HTTP service online at port 80"))
	})

})
