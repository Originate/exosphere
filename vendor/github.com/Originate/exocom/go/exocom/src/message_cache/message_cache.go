package messageCache

import (
	"errors"
	"sync"
	"time"
)

// MessageCache records the timestamp of each message.
// The cache automitically deletes any message older then 1 minute
type MessageCache struct {
	Cache map[string]time.Time
	mutex *sync.Mutex
}

// NewMessageCache returns a new MessageCache
// removing old messages with a frequency equal to the given duration
func NewMessageCache(cleanupInterval time.Duration) *MessageCache {
	result := new(MessageCache)
	result.mutex = &sync.Mutex{}
	go func() {
		for {
			time.Sleep(cleanupInterval)
			result.mutex.Lock()
			result.clearCache()
			result.mutex.Unlock()
		}
	}()
	return result
}

func (c *MessageCache) clearCache() {
	for id, timestamp := range c.Cache {
		if time.Since(timestamp) > time.Minute {
			delete(c.Cache, id)
		}
	}
	return
}

// Get returns the timestamp for the given messageId returning an error if
// no data is available
func (c *MessageCache) Get(messageID string) (time.Time, error) {
	c.mutex.Lock()
	timestamp, ok := c.Cache[messageID]
	if ok {
		c.mutex.Unlock()
		return timestamp, nil
	}
	c.mutex.Unlock()
	return timestamp, errors.New("MessageId does not exist")
}

// Set adds the given messageId and timestamp to the cache
func (c *MessageCache) Set(messageID string, timestamp time.Time) {
	c.mutex.Lock()
	c.Cache[messageID] = timestamp
	c.mutex.Unlock()
}
