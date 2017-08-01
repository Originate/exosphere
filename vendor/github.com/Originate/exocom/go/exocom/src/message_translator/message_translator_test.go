package messageTranslator_test

import (
	"github.com/Originate/exocom/go/exocom/src/message_translator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageTranslator", func() {
	Describe("InternalMessageName", func() {
		It("translates the given message to the internal format of the given sender", func() {
			input := messageTranslator.GetInternalMessageNameOptions{
				Namespace:         "text-snippets",
				PublicMessageName: "tweets.create"}
			result := messageTranslator.GetInternalMessageName(&input)
			Expect(result).To(Equal("text-snippets.create"))
		})

		It("does not translate the given message if the internal namespace matches the internal namespace of the message", func() {
			input := messageTranslator.GetInternalMessageNameOptions{
				Namespace:         "users",
				PublicMessageName: "users.create"}
			result := messageTranslator.GetInternalMessageName(&input)
			Expect(result).To(Equal("users.create"))
		})

		It("does not translate the given message if it doesn't have a namespace", func() {
			input := messageTranslator.GetInternalMessageNameOptions{
				Namespace:         "users",
				PublicMessageName: "foo bar"}
			result := messageTranslator.GetInternalMessageName(&input)
			Expect(result).To(Equal("foo bar"))
		})

		It("does not translate the given message if no internal namespace is provided", func() {
			input := messageTranslator.GetInternalMessageNameOptions{
				Namespace:         "",
				PublicMessageName: "foo.bar"}
			result := messageTranslator.GetInternalMessageName(&input)
			Expect(result).To(Equal("foo.bar"))
		})
	})

	Describe("PublicMessageName", func() {
		It("does not convert messages that don't match the format", func() {
			input := messageTranslator.GetPublicMessageNameOptions{
				Namespace:           "text-snippets",
				ClientName:          "tweets",
				InternalMessageName: "foo bar"}
			result := messageTranslator.GetPublicMessageName(&input)
			Expect(result).To(Equal("foo bar"))
		})
		It("does not convert messages that have the same internal and external namespace", func() {
			input := messageTranslator.GetPublicMessageNameOptions{
				Namespace:           "users",
				ClientName:          "users",
				InternalMessageName: "users.create"}
			result := messageTranslator.GetPublicMessageName(&input)
			Expect(result).To(Equal("users.create"))
		})
		It("does not convert messages if the service has no internal namespace", func() {
			input := messageTranslator.GetPublicMessageNameOptions{
				Namespace:           "",
				ClientName:          "users",
				InternalMessageName: "users.create"}
			result := messageTranslator.GetPublicMessageName(&input)
			Expect(result).To(Equal("users.create"))
		})
		It("converts messages into the external namespace of the service", func() {
			input := messageTranslator.GetPublicMessageNameOptions{
				Namespace:           "text-snippets",
				ClientName:          "tweets",
				InternalMessageName: "text-snippets.create"}
			result := messageTranslator.GetPublicMessageName(&input)
			Expect(result).To(Equal("tweets.create"))
		})
	})
})
