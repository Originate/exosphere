<table>
  <tr>
    <td><a href="13_generic_services.md">&lt;&lt; generic services</a></td>
    <th>Adding Search</th>
    <td><a href="15_setup_scripts.md">setup scripts &gt;&gt;</a></td>
  </tr>
</table>


# Adding Search

<table>
  <tr>
    <td>
      <b><i>
        Status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


So let's build our "todo-search-service"
using an already existing off-the-shelf search service.

```
cd ~/todo-app
exo link service
```

We enter into the wizard:

<table>
  <tr>
    <th>name of the service in our app</th>
    <td>todo-search-service</td>
  </tr>
  <tr>
    <th>URL of the service</th>
    <td>https://github.com/originate/fulltext-search-service</td>
  </tr>
  <tr>
    <th>version</th>
    <td>1.2</td>
  </tr>
</table>

Exosphere now adds the new service to the
[application configuration](03_app_config.md):

__~/todo-app/application.yml__

```yml
name: Todo application
...

services:
  - web:
    - location: web-server
  - todo:
    - location: todo-service
  - search:
    - location: https://github.com/originate/fulltext-search-service
    - version: 1.2
```

We add a search field into our UI:

> TODO: change view here

And hook up the search to our application.

Let's start this up.
We see that the web server correctly sends out search requests and receives responses,
but they are empty.
If we add new todo items, they can be found by search.
In the next chapter, we are going to fix this.


<table>
  <tr>
    <td><a href="15_setup_scripts.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
