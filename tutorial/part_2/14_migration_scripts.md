<table>
  <tr>
    <td><a href="13_add_search_service.md">&lt;&lt; add search service</a></td>
    <th>Migration Scripts</th>
    <td><a href="15_network_objects.md">network objects &gt;&gt;</a></td>
  </tr>
</table>


# Migration Scripts

<table>
  <tr>
    <td>
      <b><i>
        Status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


Our search now works for new records,
but not for old ones.
What we need is to add all the existing todo items to search.
This is done via __deployment scripts__.
These scripts help set up new services, or new versions of existing services,
as they are deployed.

This often means migrating data from an existing system into the new system.
Exosphere runs all versions of all services in parallel,
meaning if you deploy a new version of a service, the old version keeps running.
The new service is not brought online until its deployment scripts has run successfully.

A deployment script usually sends Exocom messages to other services to read their data,
and sends them to the new service to write data into it.
As such, it can be written in any language that has an Exorelay,
similar to services themselves.

In our case, we want to read all existing todo entries from the todo service,
and tell the todo-search service about them.
Once we are done, the search service will be kept up to date by the web server.

Since our script will run on production data,
we will TDD it.
Here is the spec for it:

```cucumber
Feature: deploying the todo-search service v0.0.1

  As a developer adding the search service
  I want that it gets prepopulated with all existing todo items
  So that my users can search for all their existing todos.

  Scenario: some todos exist
    Given a todos service with the entries:
      | TITLE |
      | one   |
      | two   |
    And a todos-search service
    When the deployment script runs
    Then the todos-search service knows about these documents:
      | TITLE |
      | one   |
      | two   |
    And a todos service still contains these entries:
      | TITLE |
      | one   |
      | two   |
```

The migration contains two scripts:
1. __offline migration:__
   this script runs before the service is taken online,
   i.e. connected to production traffic.
   It is used to make the service operable with ongoing production traffic,
   for example by populating it with data that is essential to even work.

2. __online migration:__
   this script runs after the service has been taken online.
   It is used to finalize it while it is already serving production traffic.

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
    <td><a href="14_migration_scripts.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
