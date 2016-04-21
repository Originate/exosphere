<table>
  <tr>
    <td><a href="03_communication.md">&lt;&lt; message-oriented programming</a></td>
    <th>Reusability</th>
    <td><a href="07_nox.md">service-oriented data modeling &gt;&gt;</a></td>
  </tr>
</table>


# Reusability

<table>
  <tr>
    <td>
      Status: idea (not implemented yet)
    </td>
  </tr>
</table>


Next let's add the ability to search for todo items.
We could add a "search" action to our _todo-service_
and try to implement search using the existing todo database.

There are several reasons why this is not a good idea though:

* From a technical perspective,
  our todo-service database is optimized for storing lots of data and retrieving pieces of it,
  not for searching through large amounts of data.
  For the latter we need a dedicated search engine,
  and a developer team that knows about search and configures the engine according to search traffic.
  The search traffic might be quite different than the load-store traffic,
  so if both teams worked on the same database,
  they might end up disagreeing a lot and stepping on each others toes.
  Since search index updates are usually more expensive than normal database updates,
  we also might want to update the search engine only every couple of seconds
  with a batch of all updates in that timespan,
  which would be completely inappropriate behavior for the generic todo database.

* From an architectural perspective,
  searching through large quantities of data is also quite a different functionality than storing it.
  In a micro-service world,
  where each service provides only one functionality,
  this also means we should have two separate services.

So let's build this as a separate "todo-search-service".
We could build another exoservice, similar to the "todo-service",
which uses a search engine as its "database".
Dealing with search engines is too complex for this tutorial though,
and the search functionality wouldn't be specific to our application anyways.
We just store the text of each todo item together with some metadata,
and then search over both in various ways.
An off-the-shelf search service should be good enough to at least get us started here.
If we ever outgrow it, we can easily migrate to a more custom search service later,
since being implemented as its own service cleanly separates it
from all the other parts of our application.

```
cd ~/todo-app
exo add external-service todo-search-service https://github.com/originate/elastic-search-service
```

With this command we tell Exosphere:
Please add an already existing service to my current application.
Name it "todo-search-service".
You find its source code at https://github.com/originate/elasticsearch-service

Exosphere now downloads the source code for the service from the given URL,
and adds it to the [application configuration](01_app_config.md).


## Translating messages

How do we communicate with the new service?
If we look at the source code for this service,
we see that it listens to `elasticsearch.store` and `elasticsearch.search` commands.

This is a problem.
Our current messages have this nice domain-specific vocabulary like `todo.create`,
where the commands tell us directly what they do in application terminology
(it creates a todo entry).
We want to keep this domain-specific language,
and for example say `todo.search`,
instead of using vocabulary that leaks the underlying technical implementation
like `elasticsearch.search`.

Yet, the elasticsearch service correctly listens to `elasticsearch.search`,
and we don't want to change that either.

Exosphere supports this by performing a translation of messages at runtime.
Our application has an instance of the elasticsearch service named "todo-search-service".
If we send a command named `todo-search-service.store` to ExoCom,
it would send it as "elasticsearch.store" to the elasticsearch-service.
How does this happen, and why?

ExoCom expects messages to have this format:

```
[service name] [separator] [command name]
```

If ExoCom comes across a command call with this structure,
it replaces the _service name_ with the class of the service.

<img src="06_schema.png" width="514" height="469">

This allows not only to reuse off-the-shelf services with domain-specific names,
but also to run several instances of the same service in parallel under different names.
This comes in extremely handy when putting together prototypes quickly
out of generic components.
