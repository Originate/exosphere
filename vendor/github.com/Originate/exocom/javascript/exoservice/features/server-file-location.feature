Feature: server file location

  As an Exoservice developer
  I want to have flexibility on where I put by server file
  So that I can keep simple services simple and at the same time organize complex services.

  Rules:
  - server.ls can be in the home directory or in the "src" folder


  Background:
    Given an ExoCom instance


  Scenario: server.ls is in the home directory
    Then it can run the "server-in-root" service


  Scenario: server.ls is in the "src" folder
    Then it can run the "server-in-src" service

  Scenario: dependencies also contain "server.js"
    Then it can run the "dependencies-contain-server-js" service
