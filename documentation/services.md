# Services

Exosphere makes it very easy to create applications consisting of lots of backend services,
and strongly encourages this pattern.


## Microservice-oriented Architectures

Microservices not only break up complex monolithic code bases,
but also complex monolithic databases
into dedicated, stand-alone, loosely coupled _network objects_.

This means the traditional way of modeling data as tables and relations
is no longer the most appropriate way of modeling your domain.

Rather, your services encapsulate data and behavior over that data.
You send them messages to modify their internal state or query it.
In summary, your application becomes a _distributed object-oriented NoSQL database_
with a domain-specific API automatically built into it.
For example, instead of the abstract SQL query:

```sql
UPDATE users
SET status=20
WHERE id=XXX
```

you now say

```livescript
send 'users.onboard', id: XXX
```

If you used a higher-level ORM that hides the complexity of the underlying data query API,
your API stays the same:

```javascript
user.onboard()
```


### Differences to relational databases

* Instead of tables you store data in services.
  Each service can be a table or even just a row within a table
* Services can (and should) store references (ids) of objects in other services.
* You have no JOIN (like in most scalable systems),
  but can still do subqueries
  (which are equally expressive as JOINs, and often easier to understand).
* Instead of SQL use
  [GraphQL](https://facebook.github.io/react/blog/2015/05/01/graphql-introduction.html)
  to query your service fleet on a high level.

  Instead of

  ```sql
  SELECT photos.*
  FROM users u INNER JOIN
       friends ON user.id=friends.user_id INNER JOIN
       photos ON photos.user_id=friends.friend_id
  WHERE users.id=1
  ```

  you now run

  ```graphql
  {
    user(id: 1) {
      friends {
        photo
      }
    }
  }
  ```


## Implementing services

### Full-Stack Services

Any code base can act as an Exoservice as long as it is able to send and receive
Exosphere messages.
Implementing this functionality is easy, since messages are transmitted via
simple HTTP requests.

The Exosphere SDK provides libraries called _communication relays_
that provide those communication facilities for all mainstream stacks:
* [ExoRelay-JS](https://github.com/Originate/exorelay-js) for Node.js
  code bases.


### Lambda Services

Lamba services are the easiest way to implement microservices.
They contain of 2 files:
* a configuration file that specifies the service
* a service file that provides handlers for all incoming message types

Exosphere provides frameworks to write lambda services in most popular languages:
* [exoservice-js](https://github.com/Originate/exoservice-js) for services in Node.js

Example lambda services:
* [users service](https://github.com/Originate/exosphere-users-service)
