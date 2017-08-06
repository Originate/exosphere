package dockercompose_test

import (
	"github.com/Originate/exosphere/exo-go/src/dockercompose"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetServiceExitCode", func() {

	It("should return the correct exit code", func() {
		role := "tweets-service"
		dockerComposeLog := "tweets-service exited with code 1"
		exitCode, err := dockercompose.GetServiceExitCode(role, dockerComposeLog)
		Expect(err).NotTo(HaveOccurred())
		Expect(exitCode).To(Equal(1))
	})

	It("should return an error when exit code does not exist", func() {
		role := "tweets-service"
		dockerComposeLog := "tweets-service is running"
		_, err := dockercompose.GetServiceExitCode(role, dockerComposeLog)
		Expect(err).To(HaveOccurred())
	})

})
