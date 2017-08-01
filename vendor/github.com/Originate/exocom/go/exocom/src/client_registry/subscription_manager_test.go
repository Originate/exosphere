package clientRegistry_test

import (
	"github.com/Originate/exocom/go/exocom/src/client_registry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionManager", func() {
	var manager *clientRegistry.SubscriptionManager

	BeforeEach(func() {
		table := clientRegistry.Routes{
			"user": clientRegistry.Route{
				Receives: []string{"users.create"},
			},
			"tweet": clientRegistry.Route{
				InternalNamespace: "text-snippets",
				Receives:          []string{"text-snippets.create"},
			},
		}
		manager = clientRegistry.NewSubscriptionManager(table)
	})

	Describe("GetSubscribersFor", func() {
		Describe("no subscribers", func() {
			It("returns an empty array", func() {
				result := manager.GetSubscribersFor("users.create")
				Expect(result).To(BeEmpty())
			})
		})

		Describe("with a subscriber without an internal namespace", func() {
			BeforeEach(func() {
				manager.AddAll("user")
			})

			It("returns the subscriber", func() {
				result := manager.GetSubscribersFor("users.create")
				Expect(result).To(Equal([]clientRegistry.Subscriber{{ClientName: "user"}}))
			})
		})

		Describe("with a subscriber with an internal namespace", func() {
			It("returns the subscriber", func() {
				manager.AddAll("tweet")
				result := manager.GetSubscribersFor("tweet.create")
				Expect(result).To(Equal([]clientRegistry.Subscriber{{
					ClientName:        "tweet",
					InternalNamespace: "text-snippets",
				}}))
			})
		})
		Describe("with a subscriber that has been removed", func() {
			It("does not return the subscriber", func() {
				manager.AddAll("tweet")
				manager.RemoveAll("tweet")
				result := manager.GetSubscribersFor("tweet.create")
				Expect(result).To(BeEmpty())
			})
		})
		Describe("with multiple subscribers", func() {
			It("returns the subscribers", func() {

			})
		})

	})
})
