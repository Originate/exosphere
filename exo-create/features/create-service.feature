@verbose
Feature: create a reusable service

  As an application developer
  I want to be able to scaffold a reusable service
  So that I can create business logic that can be used in several apps.

  Rules:
  - call "exo create service [<name>] [<template>] [<model-name>] [<description>]" to scaffold a reusable service

  Scenario: create reusable service
    Given I am in the root directory of an empty application called "empty app"
    When executing "exo-create service users-service test-author exoservice-es6-mongodb user manage users"
    Then my application contains the file "application.yml" with the content:
      """
      name: empty app
      description: Empty test application
      version: 1.0.0

      services:
        public:
          users-service:
            location: ../users-service
      """
    And my workspace contains the file "../users-service/service.yml" with content:
      """
      title: users-service
      description: manage users
      author: test-author

      setup: yarn install
      startup:
        command: node node_modules/exoservice/bin/exo-js
        online-text: online at port
      tests: node_modules/cucumber/bin/cucumber.js

      messages:
        receives:
          - user.create
          - user.create_many
          - user.delete
          - user.list
          - user.read
          - user.update
        sends:
          - user.created
          - user.created_many
          - user.deleted
          - user.listing
          - user.details
          - user.updated

      dependencies:
        mongo:
      """
