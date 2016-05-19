<table>
  <tr>
    <td><a href="16_write_setup_scripts.md">&lt;&lt; write a setup script</a></td>
    <th>Service Oriented Data Modeling</th>
    <td><a href="18_things_to_look_out_for.md">things to look out for &gt;&gt;</a></td>
  </tr>
</table>


# Network objects

Let's take an overall look at our application architecture:

<img src="17_architecture.png" width="316" height="372">

If we added a few more features and services,
like for user accounts, login, and session management,
it will look something like this:

<img src="17_architecture_full.png" width="538" height="370">

Each of the services at the bottom has its own data store.
This means we have not only a distributed code base,
but also a distributed data base.
Or, more precisely, our services encapsulate
data and behavior around this data
into dedicated entities
that communicate via messages.
Does that sound familiar?
It matches the original definition of __object-orientation__!
* everything is an object (service)
* objects (services) communicate by sending and receiving messages
* objects (services) have their own memory (state)
* every object (service) is an instance of a class (service type)
* the class (service type) holds the shared behavior for its instances (services)

Services are __network objects__.
Services are in many ways more object-like than traditional OO classes.
For example,
they communicate using real messages instead of method calls,
and they run concurrently with each other.
At the same time,
they avoid a lot of the excess that gave traditional OO a bad reputation,
like inheritance.

Instead of modeling our application's behavior as a monolithic code base
that uses a monolithic database of tables and relations,
a micro-service oriented application
is a distributed, object-oriented fusion of code and data
with a built-in domain-specific API!
This API allows us to not simply query and update data,
but to interact with it on a higher level.

Our application only scales the parts of it that are used a lot,
and replaces individual services instances when they get old,
similar to how a biological organism exchanges its cells.
This means code in a service-oriented architecture
can for example leak some memory without that being so much of a problem,
and our application can satisfy many different types of traffic patterns at the same time.

Many of the structural design patterns for OO development can be applied to services:
* __Adapter:__ a service provides a different interface to another service
* __Decorator:__ a service provides the same interface as another service, but with additional functionality
* __Facade:__ a service provides a simplified interface to an entire subsystem of services
* __Proxy:__ a local service acts as an interface for a remote system

Be mindful when using these structural patterns though - they incur extra network trips!


## Differences to the relational data model

* instead of tables you store data in services
* each service can be a table, a part of a table, or even just a single column of a table
* services can (and should) store references (ids) of related objects in other services
* you have no JOIN (like in most scalable, distributed systems),
  but can do subqueries.
  The good news is that subqueries are equally expressive as JOINs,
  and often easier to understand.


Simple queries are pretty comparable. Instead of:

```sql
INSERT INTO tweets (name) VALUES "hello world"
```

you now say:

```livescript
send 'tweets.create', name: "hello world"
```

Complex queries are easier now if using a high-level service query language like
[GraphQL](https://facebook.github.io/react/blog/2015/05/01/graphql-introduction.html).
As an example, the SQL query to get photos of all friends of user 123:

```sql
SELECT photos.*
FROM users INNER JOIN
     friends ON user.id=friends.user_id INNER JOIN
     photos ON photos.user_id=friends.friend_id
WHERE users.id=1
```

would be represented in GraphQL as:

```graphql
{
  user(id: 1) {
    friends {
      photo
    }
  }
}
```

Takeaway:

> Working with a service-oriented code base
> can be as easy as working with monolithic code and databases,
> and with the right tooling it is often simpler.
> Its all about changing your way of thinking
> away from relational to service-oriented!


