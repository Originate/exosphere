require! {
  'wait' : {repeat}
  'nanoseconds'
}


# Caches timestamp information of messages
#
# Used for calculating the processing time of clients
# by calculating the difference between the time
# when a message was sent to the client and when it sent the reply.
class MessageCache

  (@cleanup-interval = 60_000) ~>
    @cache = {}
    repeat @cleanup-interval, @cleanup


  # removes all messages older than @cleanup-interval
  cleanup: ~>
    now = nanoseconds process.hrtime!
    for id, timestamp of @cache when (timestamp - now) >= @cleanup-interval * 1e9
      @remove id


  get-original-timestamp: (message-id) ~>
    @cache[message-id]


  push: (message-id, timestamp) ~>
    @cache[message-id] = timestamp


  remove: (message-id) ~>
    delete @cache[message-id]



module.exports = MessageCache
