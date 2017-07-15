package appConfigHelpers_test

import (
	"github.com/Originate/exosphere/exo-go/src/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAppConfig", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("GetEnvironmentVariables", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("GetDependencyNames", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("GetServiceNames", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("GetSilencedDependencyNames", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("GetSilencedServiceNames", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("UpdateAppConfig", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})

var _ = Describe("VerifyServiceDoesNotExist", func() {

	const role = "exo-run"

	It("should parse non-service log message correctly", func() {
		line := "Attaching to exocom0.21.8, web"
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		Expect(serviceName).To(Equal(role))
		Expect(serviceOutput).To(Equal(line))
	})

})
