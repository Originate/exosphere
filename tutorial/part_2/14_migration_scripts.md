<table>
  <tr>
    <td><a href="13_add_search_service.md">&lt;&lt; add search service</a></td>
    <th>Migration Scripts</th>
    <td><a href="15_write_migration_script.md">write a migration script &gt;&gt;</a></td>
  </tr>
</table>


# Init Scripts

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
We need to tell the search engine about them,
so that they get added to the search index as well.

Exosphere allows to do this via __init scripts__,
which set up a freshly deployed new service
(or a new version of an existing service) into your application.
An init script consists of two parts:
* The __offline__ part runs before the new service is brought online.
  It typically pre-populates it with data that is essential for it to work,
  and can also be used to warm the service instance up to production traffic.

* The __online__ part runs after the new service is online.
  It tops it off with the most recent updates.

The init script is implemented as a separate micro service.
The respective script parts get triggered by
`<script name>-<version>-init.start-<phase>` messages
from Exosphere's deployment system
to this script, where _phase_ is either "offline" or "online".
Deployment scripts typically send out Exocom messages
that read data from other services,
and then messages to send that data to the service to be initialized.
When they are finished,
they send a `<script name>.offline-done` and `<script name>.online-done` message,
after which they get shut down by Exosphere.


Takeaway:
> Exosphere provides a way to safely and seamlessly add
> new services into an existing service fleet without interruptions of ongoing operations.


<table>
  <tr>
    <td><a href="15_write_migration_script.md"><b>&gt;&gt;</b></a></td>
  </tr>
</table>
