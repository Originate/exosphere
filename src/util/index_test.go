package util_test

import (
	"github.com/Originate/exosphere/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServiceExitCode", func() {

	It("should return the correct exit code", func() {
		role := "tweets-service"
		dockerComposeLog := "tweets-service exited with code 1"
		exitCode, err := util.GetServiceExitCode(role, dockerComposeLog)
		Expect(err).NotTo(HaveOccurred())
		Expect(exitCode).To(Equal(1))
	})

	It("should return an error when exit code does not exist", func() {
		role := "tweets-service"
		dockerComposeLog := "tweets-service is running"
		_, err := util.GetServiceExitCode(role, dockerComposeLog)
		Expect(err).To(HaveOccurred())
	})

})

var _ = Describe("ParseDockerComposeLog", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.24.0, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

	It("should parse service log message correctly", func() {
		line := "\u001b[33mweb             |\u001b[0m web server running at port 4000"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal("web"))
		Expect(serviceOutput).To(Equal("web server running at port 4000"))
	})

	It("should strip version from service name", func() {
		line := "\u001b[36mexocom0.24.0    |\u001b[0m ExoCom HTTP service online at port 80"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal("exocom"))
		Expect(serviceOutput).To(Equal("ExoCom HTTP service online at port 80"))
	})

})

var _ = Describe("Map methods", func() {

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
