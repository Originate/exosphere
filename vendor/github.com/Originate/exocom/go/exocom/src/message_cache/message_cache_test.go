package messageCache_test

import (
	"time"

	"github.com/Originate/exocom/go/exocom/src/message_cache"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageCache", func() {
	var (
		cache     *messageCache.MessageCache
		timestamp time.Time
	)
	BeforeEach(func() {
		cache = messageCache.NewMessageCache(time.Millisecond)
		timestamp = time.Now()
		cache.Cache = map[string]time.Time{}
		cache.Cache["message-name"] = timestamp
	})
	Describe("Get", func() {
		It("returns the timestamp of the message with the given messageId", func() {
			result, err := cache.Get("message-name")
			Expect(result).To(Equal(timestamp))
			Expect(err).To(BeNil())
		})
		It("returns nil if the messageid does not exist", func() {
			_, err := cache.Get("foo-bar")
			Expect(err).To(HaveOccurred())
		})
	})
	Describe("Set", func() {
		It("sets the messageId and timestamp of a message", func() {
			cache.Set("message1", timestamp)
			result, err := cache.Get("message1")
			Expect(err).To(BeNil())
			Expect(result).To(Equal(timestamp))

		})
	})
	Describe("clearCache", func() {
		It("clears the cache from messages older than one minute", func() {
			cache.Set("message1", time.Now().Add(time.Minute*-1))
			time.Sleep(time.Millisecond * 5)
			_, err := cache.Get("message1")
			Expect(err).To(HaveOccurred())
		})
		It("keeps messages in the cache that are younger than a minute", func() {
			cache.Set("message2", timestamp)
			time.Sleep(time.Millisecond * 5)
			result, err := cache.Get("message2")
			Expect(result).To(Equal(timestamp))
			Expect(err).To(BeNil())
		})
	})

})
