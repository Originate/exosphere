# ExoCom Taxonomies

ExoCom applications contain a central taxonomy
that defines how each message that is sent over the bus looks like,
what it means,
and what responses can happen in response to it.
This taxonomy is located in a dedicated file in the root of an Exosphere application.

Here is an example for the taxonomy of a simple Twitter clone.


__taxonomy.yml__
```yml
- "user signed up":
  meaning: a user just completed the sign-up process
  senders:
    HTML service: when the user has finished the sign-up form in the web UI
  listeners:
    user service: creates a user account record for this user
    welcome email service: sends a welcome email to the user's email address
  TTL: infinite
- "user sent tweet":
  meaning: a user has just created a new tweet.
  senders:
    - HTML service: when the user submitted a new tweet via the web UI
    - API service: when the user submitted a new tweet via the API
  listeners:
    - tweets service: stores the tweet
    - newsfeed service: adds the tweet to all the newsfeeds of the user's followers
- "user favorited tweet":
  meaning: a user has just favorited a tweet
  senders:
    - HTML service: when a user favorites a tweet via the web UI
    - API service: when a user favorites a tweet via the API
  listeners:
    - favorites-notification-service: sends an email notification to the creator of the tweet
    - newsfeed service: bumps up the relevance of this tweet in some of their newsfeeds
```

`TTL` defines the persistence of the message.
If an infinite TTL is given,
the bus guarantees at-least-once delivery of the message.
This is for messages that represent new important insights into the domain
that should be kept around.
For example "user has liked a tweet".

If a non-infinite TTL is given,
the bus only makes a best-effort delivery of the message.
This is for messages that are a part of an interactive session with the user.
Example: "I need a list of all current user accounts".
There is no point in archiving and replaying those messages at later times,
beyond best-effort in the moment because the interactive activity will time out in 100ms anyways
and show an error message.

Persistent messages are delivered via

All these services can operate in realtime or in batch mode.
In batch mode they wait for some period of time,
for example the end of the day,
and then process all accumulated messages in one batch.
In realtime mode they process each message right when it occurs.
