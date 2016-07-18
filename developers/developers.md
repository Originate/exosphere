# Exosphere Developer Documentation

Exosphere consists of these sub-projects:

<table>
  <tr>
    <td><a href="https://github.com/Originate/exosphere-sdk">Exosphere SDK</a></td>
    <td>
      <a href="https://circleci.com/gh/Originate/exosphere-sdk">
        <img src="https://circleci.com/gh/Originate/exosphere-sdk.svg?style=shield&circle-token=fc8148ed828cc81e6ca44920672af8f773106795">
      </a>
      <a href="https://david-dm.org/originate/exosphere-sdk">
        <img src="https://david-dm.org/originate/exosphere-sdk.svg">
      </a>
      <a href="https://david-dm.org/originate/exosphere-sdk#info=devDependencies">
        <img src="https://david-dm.org/originate/exosphere-sdk/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exocom-dev">ExoCom-Dev</a></td>
    <td>
      <a href="https://github.com/Originate/exocom-dev">
        <img src="https://circleci.com/gh/Originate/exocom-dev.svg?style=shield&circle-token=0f68f90da677a3c5bffc88d9d41910c00f10b81e">
      </a>
      <a href="https://david-dm.org/originate/exocom-dev">
        <img src="https://david-dm.org/originate/exocom-dev.svg">
      </a>
      <a href="https://david-dm.org/originate/exocom-dev#info=devDependencies">
        <img src="https://david-dm.org/originate/exocom-dev/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exoservice-js">Exoservice-JS</a></td>
    <td>
      <a href="https://circleci.com/gh/Originate/exoservice-js">
        <img src="https://circleci.com/gh/Originate/exoservice-js.svg?style=shield&circle-token=33fbf4fc2b0c128479443c5e8bff337815205ec7">
      </a>
      <a href="https://david-dm.org/originate/exoservice-js">
        <img src="https://david-dm.org/originate/exoservice-js.svg">
      </a>
      <a href="https://david-dm.org/originate/exoservice-js#info=devDependencies">
        <img src="https://david-dm.org/originate/exoservice-js/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exorelay-js">Exorelay-JS</a></td>
    <td>
      <a href="https://github.com/Originate/exorelay-js">
        <img src="https://circleci.com/gh/Originate/exorelay-js.svg?style=shield&circle-token=012a2c6405c702e0a8271de804eed0c4c179772f">
      </a>
      <a href="https://david-dm.org/originate/exorelay-js">
        <img src="https://david-dm.org/originate/exorelay-js.svg">
      </a>
      <a href="https://david-dm.org/originate/exorelay-js#info=devDependencies">
        <img src="https://david-dm.org/originate/exorelay-js/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exocom-mock-js">ExoCom-Mock-JS</a></td>
    <td>
      <a href="https://github.com/Originate/exocom-mock-js">
        <img src="https://circleci.com/gh/Originate/exocom-mock-js.svg?style=shield&circle-token=4f522d83e80f98f58b30cd1c9ad7f2e24f8e0b58">
      </a>
      <a href="https://david-dm.org/originate/exocom-mock-js">
        <img src="https://david-dm.org/originate/exocom-mock-js.svg">
      </a>
      <a href="https://david-dm.org/originate/exocom-mock-js#info=devDependencies">
        <img src="https://david-dm.org/originate/exocom-mock-js/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exosphere-mongodb-service">MongoDB service</a></td>
    <td>
      <a href="https://github.com/Originate/exosphere-mongodb-service">
        <img src="https://circleci.com/gh/Originate/exosphere-mongodb-service.svg?style=shield&circle-token=389739b88cceec7155d0253e1560339a8409fd98">
      </a>
      <a href="https://david-dm.org/originate/exosphere-mongodb-service">
        <img src="https://david-dm.org/originate/exosphere-mongodb-service.svg">
      </a>
      <a href="https://david-dm.org/originate/exosphere-mongodb-service#info=devDependencies">
        <img src="https://david-dm.org/originate/exosphere-mongodb-service/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exosphere-users-service">Users Service</a></td>
    <td>
      <a href="https://github.com/Originate/exosphere-users-service">
        <img src="https://circleci.com/gh/Originate/exosphere-users-service.svg?style=shield&circle-token=b8da91b53c5b269eeb2460e344f521461ffe9895">
      </a>
      <a href="https://david-dm.org/originate/exosphere-users-service">
        <img src="https://david-dm.org/originate/exosphere-users-service.svg">
      </a>
      <a href="https://david-dm.org/originate/exosphere-users-service#info=devDependencies">
        <img src="https://david-dm.org/originate/exosphere-users-service/dev-status.svg">
      </a>
    </td>
  </tr>
  <tr>
    <td><a href="https://github.com/Originate/exosphere-tweets-service">Tweets Service</a></td>
    <td>
      <a href="https://github.com/Originate/exosphere-tweets-service">
        <img src="https://circleci.com/gh/Originate/exosphere-tweets-service.svg?style=shield&circle-token=b571517a2b36b03bd440ad7056d2a072c463dc63">
      </a>
    </td>
  </tr>
</table>


## Application Architecture

Exosphere makes it very easy to create
applications consisting of lots of backend services,
and strongly encourages this pattern.
The graphic below shows the architecture of a simple Exosphere application.

![architecture diagram](architecture.gif)

* the application code is factored into individual services
  (shown in green)
* each service has one responsibility, and provides that generic and well
* service communicate via JSON (the universal data exchange language)
* ExoSphere (shown in blue)
  provides a shared communication bus called _ExoCom_
  for this application.
* the communication bus can be extended to clients via the _ExoCom gateway_
* Exosphere provides infrastructure services
  like full-stack deployment, analytics, and devops support
  for the application


