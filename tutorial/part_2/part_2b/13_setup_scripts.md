<table>
  <tr>
    <td><a href="12_add_search_service.md">&lt;&lt; add search service</a></td>
    <th>Setup Scripts</th>
    <td><a href="14_write_setup_script.md">write a setup script &gt;&gt;</a></td>
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

The Exosphere runtime provides a mechanism called __setup scripts__ to do exactly this.
These scripts accompany a particular version of an application,
and run after this version is deployed into an environment.
They typically set up the new services.
As an example,
the setup script for the new version of our todo application - which adds search -
would read all todo items from the todo service
and add them to the search service.

In development mode,
you can run setup scripts locally as often as you like.
You can (and should) write tests for them,
since they change production data and behavior.
Once the new application version is deployed,
for example into your _staging_ or _production_ environment,
the setup script will run automatically.

Upgrade scripts are implemented as separate micro services,
and can therefore be written in your language of choice,
the way you write normal services.
Your setup script must wait for the
`<app name>-<version>-setting-up` message
from Exosphere's deployment system.
Setup scripts typically send out Exocom messages
to read data from some services
and send it to the service to be initialized.
When they are finished,
they send an `<app name>-<version>-set-up` message,
after which they get shut down by the Exosphere runtime.


## Release Strategies

Setup scripts allow different roll-out strategies for new features,
ranging from simple to sophisticated.

You can keep things easy and efficient and
simply roll out a new feature all at once
and make everything work once it is live.
This works if you are confident
that you can set up the new services
in a reasonably short amount of time.
An example is the search feature in our todo application:
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
> new features and services into a running application
> without interruptions of ongoing operations.


<table>
  <tr>
    <td><a href="14_write_migration_script.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
