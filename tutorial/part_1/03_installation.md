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
Developing applications consisting of many different code bases requires:
* a __package management system__ for installing/updating the various pieces of platforms and tools as needed.
  You don't want to get into the business of having to install several programming languages,
  frameworks, and other dependencies of your polyglot applications manually.
* a __compiler__ to be used by the package manager and applications running inside Exosphere

The installation instructions below include
making these infrastructure components available on your system.


## MacOS (OS X)

The easiest way to install Exosphere on OS X is to run its installer app:

```
curl https://raw.githubusercontent.com/Originate/exosphere-sdk/master/bin/install_osx | sh
```

This script will check whether all needed components are present and working on your system,
and install missing ones.
It confirms with you before making any changes to your system, so running it is safe.
You can exit it at any time via __Ctrl-C__,
and re-run it at any time to verify proper installation.


### Manual installation on OS X

To install the SDK manually:

* install [Xcode](https://itunes.apple.com/us/app/xcode/id497799835) via the Apple App Store
* install [Homebrew](http://brew.sh): `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* install [ZeroMQ](http://brewformulas.org/zeromq): `brew install zeromq`
* install [pkg-config](http://brewformulas.org/pkg-config): `brew install pkg-config`
* install [Node.js](https://nodejs.org) version 6 or above: `brew install node`
* install Exosphere: `npm i -g exosphere-sdk`
* verify that you can run exosphere commands: `exo version`


## Windows

### Windows 7

Some of the commands below have to be run in an __administrative command shell__.
Instructions for how to open such a shell are available
for [windows 8-10](http://www.howtogeek.com/194041/how-to-open-the-command-prompt-as-administrator-in-windows-8.1)
and [Windows 7](http://www.howtogeek.com/howto/windows-vista/run-a-command-as-administrator-from-the-windows-vista-run-box)).

* install the package manager
  * in an administrative shell, install [Chocolatey](https://chocolatey.org/install)
  * when done, close the current shell and open a new one to load the environment changes prepared by the installer
* install the compiler
  * in another administrative shell, install [Node.js](http://nodejs.org) version 6 or above: `choco install nodejs.install -y`
  * in a normal shell, install npm@next: `npm i -g npm@next`
  * in an administrative shell, install the [compilers](https://github.com/felixrieseberg/windows-build-tools): `npm install --global --production windows-build-tools`
* install Exosphere:
  * in a normal shell: `npm install --global exosphere-sdk`
  * close this shell, open a normal one, and run `exo` to make sure it works
* install other tools you will need
  * in an administrative shell: [Git](https://git-scm.com) (`choco install git.install -y`
  * optionally [Github Desktop](https://desktop.github.com) (install manually, the choco package is broken)


### Contributing

If you want to become an Exosphere contributor and you are on Windows,
you need some additional infrastructure:
* [Git](https://git-scm.com) for the Git command line and Git Bash: `choco install git.install -y`
* you also need to set up SSH keys for Github.
  The easiest way is via [Github Desktop](https://desktop.github.com): `choco install github -y`
  (Note: this install is currently broken, you might want to install this manually)

Please perform all



## Linux

Installation instructions given for [Ubuntu](http://www.ubuntu.com),
please adapt them to your distro as needed:
* install [ZeroMQ](http://zeromq.org)
* install [Node.js](https://nodejs.org) version 6 or above: `sudo apt-get install node`


<table>
  <tr>
    <td><a href="../part_2/readme.md"><b>&gt;&gt;</b></td>
  </tr>
</table>
