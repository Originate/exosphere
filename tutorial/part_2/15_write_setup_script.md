<table>
  <tr>
    <td><a href="14_setup_scripts.md">&lt;&lt; setup scripts</a></td>
    <th>Write A Setup Script</th>
    <td><a href="16_network_objects.md">network objects &gt;&gt;</a></td>
  </tr>
</table>


# Writing the setup script for the search service

<table>
  <tr>
    <td>
      <b><i>
        Status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


To set up the search feature of our todo application,
we want to read all existing todo entries from the todo service
and send them to the todo-search service.
Since our script will touch production data,
we will TDD it.
Here is the spec for it:

```cucumber
Feature: v0.0.2 deployment

  As a developer launching the search feature of the todo application
  I want to prepopulate it with the existing todo items in the system
  So that my users can search for all their existing todos.

  Scenario: some todo entries exist
    Given the todos service has the entries:
      | TEXT          |
      | take medicine |
      | go party      |
    When the setup script for version 0.0.2 has run
    Then the todos-search service knows about these documents:
      | TEXT          |
      | take medicine |
      | go party      |
```

Let's run it to verify the tests fail:

```
$ exo setup
```

Exosphere runs the migration for the current version of the application.
It tells us that there is no setup script for version 0.0.2.
Let's create one:

__~/todo-app/setup-scripts/0.0.2.js__

```javascript
module.exports = {

  "todo-app-0.0.2-setting-up": (_, {reply}) => {
    send('todos.list', (todos) => {
      for todo in todos
        send('todo-search.add', todo);
    }
  }

};
```

When running our setup script again (`exo setup`),
the terminal shows its activities.
Now the web application can search for all todo items!


Takeaway:
> Exosphere provides a way to safely and seamlessly deploy
> new services into the existing service fleet.


<table>
  <tr>
    <td><a href="16_network_objects.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>

