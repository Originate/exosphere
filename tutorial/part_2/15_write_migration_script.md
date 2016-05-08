<table>
  <tr>
    <td><a href="14_setup_scripts.md">&lt;&lt; setup scripts</a></td>
    <th>Migration Scripts</th>
    <td><a href="16_network_objects.md">network objects &gt;&gt;</a></td>
  </tr>
</table>


# Writing a migration script for our search service

<table>
  <tr>
    <td>
      <b><i>
        Status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


In our case, we want to read all existing todo entries from the todo service
and send them to the todo-search service.
Once we are done, the search service will be automatically kept up to date
by the web server.
Since our script will touch production data,
we have to TDD it.
Here is the spec for it:

```cucumber
Feature: todo-search service v0.0.1 deployment

  As a developer adding the search to the todo application
  I want that it gets prepopulated with all existing todo items in the system
  So that my users can search for all their existing todos.

  Scenario: some todo entries exist
    Given a todos service with the entries:
      | TEXT |
      | one  |
      | two  |
    And a currently deploying todos-search service at version 0.0.1
    When its deployment script runs
    Then the todos-search service knows about these documents:
      | TEXT |
      | one  |
      | two  |
    And a todos service still contains these entries:
      | TEXT |
      | one  |
      | two  |
```

In our case,
the search service is _operable_ with production traffic right away.
It can search for things, and update newly created todo entries.
We just need to fill it up with existing todo items
to make the search more complete.
So in our case we only need the _online migration_.


```javascript
module.exports = {

  "deployment-todo-search-0.0.1-offline-start": (_, {reply}) => {
    reply('deployment-todo-search-0.0.1-offline-done');
  },

  "deployment-todo-search-0.0.1-online-start": (_, {reply}) => {
    send('todos.list', (todos) => {
      for todo in todos
        send('todo-search.add', todo);
      reply("deployment-todo-search-0.0.1-online-done");
    });
  }

};
```

Let's run it by starting the new version of Exosphere with the todo-search service activated!

```
$ exo run
```

Exosphere keeps track about which deployment scripts it has run already.
Since it didn't run the one we just created, it does so now.
The terminal shows log entries for the activities done by the deployment script.
Now we can search for all todo items!


Takeaway:
> Exosphere provides a way to safely and seamlessly deploy
> new services into the existing service fleet.


<table>
  <tr>
    <td><a href="16_network_objects.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>

