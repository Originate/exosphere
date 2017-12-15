package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppRemoteConfig", func() {
	var appRemoteEnvironment types.AppRemoteEnvironment

	var _ = Describe("ValidateFields", func() {
		It("should throw an error when AppConfig is missing fields in remote", func() {
			appRemoteEnvironment = types.AppRemoteEnvironment{
				URL:       "originate.com",
				AccountID: "123",
				Region:    "us-west-2",
			}
			err := appRemoteEnvironment.ValidateFields("qa")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "application.yml missing required field 'remote.environments.qa.SslCertificateArn'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("should not throw an error when AppConfig isn't missing fields", func() {
			appRemoteEnvironment = types.AppRemoteEnvironment{
				URL:               "originate.com",
				AccountID:         "123",
				Region:            "us-west-2",
				SslCertificateArn: "cert-arn",
			}
			err := appRemoteEnvironment.ValidateFields("qa")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
