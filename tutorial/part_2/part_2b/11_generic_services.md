<table>
  <tr>
    <td><a href="10_integration_into_web_server.md">&lt;&lt; integration into web server</a></td>
    <th>Generic Services</th>
    <td><a href="12_add_search_service.md">adding the search service &gt;&gt;</a></td>
  </tr>
</table>


# Generic Services

<table>
  <tr>
    <td>
      <b><i>
        Status: alpha - some parts implemented
      </i></b>
    </td>
  </tr>
</table>


Next we add the ability to search for todo items.
We could add a "search" endpoint to our exising service that stores todo items
and implement search using its internal MongoDB database.
There are several reasons why this is not a good idea though:

* We might want to use more search capabilities
  than our database provides (fuzzy search, finding synonyms, etc).

* The Exosphere runtime provides an off-the-shelf fulltext-search service.
  It took a while to build, secure, and scale it.
  We don't want to re-do all this work in each of our services that needs search.

* Search traffic patterns and availability requirements
  are often different for search than for CRUD traffic.
  Search index updates are usually more expensive than normal database updates.
  We might want to batch-update the search engine only every couple of seconds.
  This would be completely inappropriate behavior for the main database.
  Optimizing a single database for both CRUD and search traffic will be difficult at scale.

So let's use a dedicated search service here.
As already mentioned, the Exosphere framework provides a number of off-the-shelf services.
They are ready to plug into existing applications
no matter what language the rest of the system is written in.
All of these services are robust, secure, and scalable.
Here is a small selection of them:

<table>
  <tr>
    <th>users service</th>
    <td>stores user accounts (name, email, etc)</td>
  </tr>
  <tr>
    <th>password auth service</th>
    <td>
      allows user accounts to authenticate via password,
      allows to reset forgotten passwords, etc
    </td>
  </tr>
  <tr>
    <th>API token auth service</th>
    <td>
      allows user accounts to authenticate via API tokens,
      generates new random tokens, etc
    </td>
  </tr>
  <tr>
    <th>session service</th>
    <td>
      caches a small amount of session data in the backend
      and retrieves it extremely quickly
    </td>
  </tr>
  <tr>
    <th>fulltext search service</th>
    <td>provides fulltext search capabilities</td>
  </tr>
</table>

Building Exosphere applications is quite different than building monolithic applications:
Especially in the beginning,
one mostly plugs generic services together like [Lego](http://www.lego.com) bricks.
This allows to rapidly prototype fully functional and scalable applications,
and only spend time engineering application-specific components
when outgrowing the abilities of the generic building blocks<sup>1</sup>.
Replacing services in Exosphere applications is easy,
since they are naturally loosely coupled from each other
and communicate only via well-defined, wholistically type-checked APIs.


## Translating messages

In the [documentation for the search service]()
we see that it expects amongst other things
the command `fts.add` to add data to the search engine
and `fts.search` to perform a search.
When building a service that searches todo items,
we want to use domain-specific commands though,
like `todo-search.add` and `todo-search.search`.
So that searching for todos
is intuitive via domain-specific terminology,
and we can build several search services into our application.
For example, when we add search for user accounts,
we want to be able to spin up another fulltext search engine,
and communicate with it via `user-search.add` and `user-search.search` messages.

Exocom supports this by translating messages at runtime:
* the fulltext-search service expects for example the command `fts.add` for adding data
* our application contains a fulltext-search service instance called `todo-search`
* exocom subscribes this instance to `todo-search.add` commands
* when the web server sends out a `todo-search.add` command,
  exocom sends it to the fulltext-search service instance as `fts.add`
* when the search sends out an `fts.added` reply,
  exocom translates it back to `todo-search.added`,
  and sends the latter back to the web server

<img src="12_schema.png" width="514" height="354">

Thanks to this message translation feature of exocom,
our search service uses its own vocabulary internally,
and is accessible via its name from the outside.
This doesn't just make it more intuitive to use.
If we would add another service to search user accounts,
we would communicate with it via `user-search.add` and `user-search.search` commands.
Both search services will be nicely separated from each other.

Putting all this together, Exocom has two layers of messaging functionality:

1. __generic messaging:__
   This is the foundation for all messaging within Exosphere.
   On this layer, a message can have any arbitrary title,
   and exocom happily forwards it to whoever is subscribed to this string,
   without modifying anything.
   This is what we use to talk to the `todo` service.

2. __NOS messaging:__
   More opinionated messaging functionality
   built on top of the generic messaging layer.
   It's name stands for _network objects_,
   since it allows to create more than one instance of a service in your application.
   NOS improves your application's architecture<sup>2</sup>.

NOS messages have this format (called _namespaced message format_):

```
[service name] [separator] [command name]
```

The separator is `.` by default.
You can define additional separators,
like for example `: ` or `/`.
This allows commands to look like
`todo-search: add` or `todo-search/add`
if you prefer one of these formats.

NOS is activated if a service defines a messaging namespace in its configuration file.
If enabled, Exocom changes all messages going in and out of the service
by replacing the _service name_ in the message title
with the _namespace_ defined by the service
(and vice versa in the other direction).


Takeaway:
> The Exosphere framework provides various layers of communication infrastructure,
> ranging from generic and flexible to opinionated and efficient.
> This allows to prototype capable applications out of prefabricated components,
> which can be evolved into fully custom-built production applications later.


Next we are going to add the service to search todo-items.


<table>
  <tr>
    <td><a href="13_add_search_service.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>


<hr>

<sup>1</sup>
Usually at this time the product idea and usage patterns have solidified substantially.
This means "real" development effort happens with a lot less risk in those areas,
and building them can happen in a much more straight line compared to
writing custom code right from the start of the project.

<sup>2</sup>
NOS is also a term for using
[nitrous oxide](https://en.wikipedia.org/wiki/Nitrous_oxide_engine)
to increase car engine performance.
Now you can also tune up your app with NOS! :)

<sup>3</sup>
