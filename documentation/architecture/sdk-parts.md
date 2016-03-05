# SDK Architecture

The Exosphere SDK is built in a very modular way.


## Subprojects

The Exosphere SDK is comprised of the following subprojects:

<table>
  <tr>
    <th><a href="https://github.com/Originate/exosphere-cli">exosphere-cli</a></th>
    <td>CLI to start/deploy whole Exosphere applications</td>
  </tr>
  <tr>
    <th><a href="https://github.com/Originate/exoservice-js">exoservice-js</a></th>
    <td>
      a high-level, opinionated framework for creating microservices in Node.js.
      An example microservice created with ExoServiceJS is
      [here](https://github.com/Originate/exosphere-users-service).
    </td>
  </tr>
  <tr>
    <th><a href="https://github.com/Originate/generator-exoservice-livescript">LiveScript Exoservice Generator</a></th>
    <td>
      Scaffolds an
      <a href="https://github.com/Originate/exoservice-js">exoservice-js</a>
      instance in LiveScript.
    </td>
  </tr>
  <tr>
    <th><a href="https://github.com/Originate/exorelay-js">exorelay-js</a></th>
    <td>
      provides Exosphere communication services to larger Node.js applications
      like web or API servers
    </td>
  </tr>
  <tr>
    <th><a href="https://github.com/Originate/exocom-dev">exocom-dev</a></th>
    <td>
      a simple in-memory implementation of the Exosphere messaging infrastructure
      for local development.
    </td>
  </tr>
  <tr>
    <th><a href="https://github.com/Originate/exocom-mock-js">exocom-mock-js</a></th>
    <td>
      a mock implementation of
      [exocom-dev](https://github.com/Originate/exocom-dev)
      for testing ExoServices
    </td>
  </tr>
</table>


## exo-cli commands

Exo-CLI provides the following commands:

* [exo-install command](../../features/install.feature):
  Sets up a freshly cloned Exosphere application
  by running the setup scripts for each service.

* [exo-run command](../../features/run.feature):
  Starts a properly set up (using `exo-install`) Exosphere application
  by launching all of its services.
