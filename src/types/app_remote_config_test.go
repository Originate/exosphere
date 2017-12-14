package types_test

import (
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppRemoteConfig", func() {
	var appRemoteConfig types.AppRemoteConfig

	var _ = Describe("ValidateFields", func() {
		It("should throw an error when AppConfig is missing fields in remote", func() {
			appRemoteConfig = types.AppRemoteConfig{
				URL:       "originate.com",
				AccountID: "123",
				Region:    "us-west-2",
			}
			err := appRemoteConfig.ValidateFields("qa")
			Expect(err).To(HaveOccurred())
			expectedErrorString := "application.yml missing required field 'remote.qa.SslCertificateArn'"
			Expect(err.Error()).To(ContainSubstring(expectedErrorString))
		})

		It("should not throw an error when AppConfig isn't missing fields", func() {
			appRemoteConfig = types.AppRemoteConfig{
				URL:               "originate.com",
				AccountID:         "123",
				Region:            "us-west-2",
				SslCertificateArn: "cert-arn",
			}
			err := appRemoteConfig.ValidateFields()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
