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

## Dependencies

Before installing Exosphere, ensure you have the following applications installed:
 * Git - [https://git-scm.com/]()
 * Docker - [https://www.docker.com]()

## General installation instructions

* Download and install the binary executable for your platform from the 
[GitHub release page](https://github.com/Originate/exosphere/releases/latest)
* Move the binary to a location in your path and name it `exo`

### Mac installation

_Note: Depending on the privileges of `/usr/local/bin` on your system these steps could require root privileges_

In a terminal, execute the following:
```bash
curl -L https://github.com/Originate/exosphere/releases/download/vX.Y.Z/exo-darwin-amd64 >/usr/local/bin/exo
chmod +x /usr/local/bin/exo
```

<table>
  <tr>
    <td><a href="../part_2/readme.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
