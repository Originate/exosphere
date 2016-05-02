<table>
  <tr>
    <td><a href="14_migration_scripts.md">&lt;&lt; migration scripts</a></td>
    <th>Service Oriented Data Modeling</th>
    <td><a href="14_migration_scripts.md">service-oriented data modeling &gt;&gt;</a></td>
  </tr>
</table>


# Network objects

Let's take an overall look at our application architecture:

<img src="15_architecture.png" width="316" height="359">

If we added a few more features and services,
like for user accounts, login, and session management,
it will look like this:

<img src="15_architecture_full.png" width="538" height="314">

Each of the services at the bottom has its own data store.
This means we have not only a distributed code base,
but also a distributed data base.
Or, more precisely, our services encapsulate
data and behavior around this data
into entities that communicate via messages.
Does that sound familiar?
It matches the original definition of __object-orientation__!
* everything is an object (service)
* objects (services) communicate by sending and receiving messages
* objects (services) have their own memory
* every object (service) is an instance of a class (service type)
* the class (service type) holds the shared behavior for its instances (services)

Services are network objects,
and instead of a separate code and database
our program has become a distributed, object-oriented NoSQL data store
with a domain-specific API automatically built into it!

This means the traditional way of modeling data
as tables and relations
is no longer the most appropriate way of modeling your domain.
You need to start thinking object-oriented.


## Differences to relational databases

* instead of tables you store data in services
* each service can be a table, a part of a table, or even just a single row of a table
* services can (and should) store references (ids) of related objects in other services
* you have no JOIN (like in most scalable systems),
  but can do subqueries
  (which are equally expressive as JOINs)


## Single queries

Instead of

```sql
INSERT INTO tweets (name) VALUES "hello world"
```

you now say:

```livescript
send 'tweets.create', name: "hello world"
```

## Complex queries

Answering complex queries can require querying several services.
This can quickly become cumbersome.
For those cases,
consider using a high-level service query language
like
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

Working with a service-oriented code base
can be as easy as working with monolithic code and databases,
and with the right tooling it is often simpler!
