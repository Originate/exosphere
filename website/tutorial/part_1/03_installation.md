<table>
  <tr>
    <td><a href="02_architecture.md">&lt;&lt; architecture</a></td>
    <th>Installation</th>
    <td><a href="../part_2/readme.md">Part II &gt;&gt;</a></td>
  </tr>
</table>


# Installation

Exosphere is not just a simple application,
but a framework for large-scale software development of micro-service based applications.
Developing applications consisting of many different code bases requires
a __package management system__ for installing/updating the various pieces of platforms and tools as needed.
  You don't want to get into the business of having to install several programming languages,
  frameworks, and other dependencies of your polyglot applications manually.

* Download and install the binary executable for your platform from the 
[GitHub release page](https://github.com/Originate/exosphere/releases/latest)
* Rename the executable to `exo`
* Check and update your `PATH` environmental variable to include its install location, if required

## Managing multiple versions of Exosphere

A common pattern for managing executables on a system is to install the executable with
the version in the filename and symlinking the desired command name to the executable.
Using this pattern, old versions are still available in your `PATH`. In the event that you need
to run an old version of Exosphere it will still be available as `exo-X.Y.Z` unless it is explicitly 
removed.

### Example: Installing version X.Y.Z of Exosphere on Mac
_Note: Depending on where you install these steps could require root privileges_

```bash
curl -L https://github.com/Originate/exosphere/releases/download/vX.Y.Z/exo-darwin-amd64 >/usr/local/bin/exo-X.Y.Z
chmod +x /usr/local/bin/exo-X.Y.Z
ln -fns /usr/local/bin/exo-X.Y.Z /usr/local/bin/exo
```

<table>
  <tr>
    <td><a href="../part_2/readme.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
