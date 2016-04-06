# Reusability

Status: idea (not implemented yet)


Next let's add the ability to _search_ for todo items.
We could add a "search" action to our _todo-service_
and try to implement search using the database that stores them.

There are several reasons why this is not a good idea though:

* From an architectural perspective,
  searching through large quantities of data is a different functionality than storing them.
  In a micro-service world this means that we should have two separate services for both,
  since each service should provide only one functionality.

* From a technical perspective,
  our todo-service database is optimized for storing lots of data and pieces of it,
  not for searching through them at scale.
  For that we need a dedicated search engine.

So we will do this as a separate "todo-search-service".

We could build another exoservice, similar to the "todo-service".
Handling search engines is pretty complex though,
and the search functionality wouldn't be specific to our application in any way.
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
we see that it listens to "elasticsearch.store" and "elasticsearch.search" commands.
Our current messages have this nice domain-specific vocabulary,
where the commands tell us directly that they "create a todo-item" or list them.
We want to keep that domain-specific language,
and not mix it with vocabulary about the underlying technical implementation
like "perform a search using ElasticSearch".

Exosphere supports this by performing a translation of messages at runtime.
Our application has an instance of the "elasticsearch-service" named "todo-search-service".
If we send a command named "todo-search-service.store" to ExoCom,
it would send it as "elasticsearch.store" to the elasticsearch-service.
How does this happen, and why?

ExoCom expects messages to have this format:

```
[service name] [separator] [command name]
```

If ExoCom comes across a command call with this structure,
it replaces the _service name_ with the class of the service.

This allows not only to give services more domain-specific names,
but also to run several instances of a service in parallel under different names.
