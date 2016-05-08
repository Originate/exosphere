<table>
  <tr>
    <td><a href="13_add_search_service.md">&lt;&lt; add search service</a></td>
    <th>Migration Scripts</th>
    <td><a href="15_write_migration_script.md">write a migration script &gt;&gt;</a></td>
  </tr>
</table>


# Setup Scripts

<table>
  <tr>
    <td>
      <b><i>
        Status: idea - not implemented yet
      </i></b>
    </td>
  </tr>
</table>


Our new search feature works for records that we create from now on,
but not for records that were already in our todo database
when we deployed search.
We need to tell the search engine about these existing records,
so that they get added to the search index as well.

Exosphere provides a mechanism called __setup scripts__ to do exactly this.
These scripts accompany a particular version of an application,
and run after this version is rolled out into an environment.
They typically set up the new services.
As an example,
the setup script for the new version of our todo application - which adds search -
would read all todo items from the todo service
and add them to the search service.

You can develop a setup script locally
and run it as often as you like
to make things work properly.
You can (and should) write tests for your script,
since they change production data and behavior.
Once the new application version is ready,
you deploy it somewhere (your _staging_ or _production_ environment),
and the setup script will run there automatically
after the deployment is done.

Upgrade scripts are implemented as separate micro services,
and can therefore be written in your language of choice,
the way you are used to.
Your setup service must respond to the
`<app name>-<version>-setting-up` message
from Exosphere's deployment system.
Setup scripts typically send out Exocom messages
to read data from some services
and send it to the service to be initialized.
When they are finished,
they send an `<app name>-<version>-set-up` message,
after which they get shut down by Exosphere.


## Release Strategies

Setup scripts allow different roll-out strategies for new features,
ranging from simple to sophisticated.

You can keep things simple and efficient by
rolling out a new feature all at once
and top it off with the latest data it once it is live.
This works if you are confident
that you can get everything set up
in a reasonably short amount of time.
In our todo example:
* the current application is v0.0.1
* we deploy a new version v0.0.2 that contains and uses the search service
* our upgrade script runs right after deployment
  and populates the search service with all todo items

There is a brief period right after deployment
where the search service doesn't contain all todo entries yet,
and search therefore doesn't return all results.
Since we can populate the search service
within a few seconds using our setup script,
this roll-out strategy is probably good enough here.

If it weren't, for example because we already have millions of todo entries,
lots of users and traffic,
and wanted to verify search works and scales in production
before making it available to the public,
we would roll out our search feature in two separate releases:
* the current application is v0.0.1
* we deploy a new version v0.0.2 that contains the search service
  but doesn't use it yet.
  In this version, the web server already sends updates to the search service,
  but doesn't show the search UI and doesn't send search queries.
* after deployment our setup script for v0.0.2 runs.
  It prepopulates the search service with all the existing data in the todo service.
  It is okay if this process takes some time,
  since the search service isn't visible to the public yet.
* once the setup script is done,
  and we have load tested and warmed up our search service,
  we deploy v0.0.3 of the application.
  It contains an updated web server that now includes search.


Takeaway:
> Setup scripts allow to safely and seamlessly add
> new services into an existing service fleet without interruptions of ongoing operations.


<table>
  <tr>
    <td><a href="15_write_migration_script.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
