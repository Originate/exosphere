# Message Translation

Exocom provides the ability to translate messages at runtime.
This allows for reusable services.
This is best illustrated by an example:

Let's say we need to store user accounts.
We decide to store them in MongoDB initially.
Also, initially we only need the standard CRUD operations for user accounts.
We use the "MongoStorage" service as the initial implementation of the users service.
This service replies to messages like `store record` or `list records`.
In our application, we want to use messages like `new user` and `list users`, though,
because they are described in the domain languege.
To use the MongoStorage service with the domain language,
we define a translation section in the configuration file of the service:

__application.yml__
```yml
services:
  private
    users-service:
      docker_image: xxx
      message-translation:
        - domain: "new user"
          internal: "store record"
        - domain: "user created"
          internal: "record stored"
        - domain: "list users"
          internal: "list records"
```

Now, each time ExoCom sees a `new user` message,
it sends it as a `create record` message to instances of the `users-service` service.
Replies from the service, for example `record created`, also get translated,
in this case into `user created`.
Translation does not change the payload.

This allows to use a generic MongoStorage service as if it were a real users service.

Using this technique, we can can use MongoStorage service instances
for other services, like for example a ProductService.
