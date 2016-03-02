# Service Oriented Data Modeling

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

followed by code somewhere else to send a welcome email to the user,
you now say:

```livescript
send 'users.onboard', id: XXX
```

If you used a higher-level ORM (or better OSM)
that hides the complexity of the underlying data query API,
your API stays the same:

```javascript
user.onboard()
```


### Differences to relational databases

* Instead of tables you store data in services.
  Each service can be a table or even just a row within a table
* Services can (and should) store references (ids) of related objects in other services.
* You have no JOIN (like in most scalable systems),
  but can still do subqueries
  (which are actually equally expressive as JOINs, and often easier to understand).
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
