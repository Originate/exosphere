@verbose
Feature: create a reusable service

  As an application developer
  I want to be able to scaffold a reusable service
  So that I can create business logic that can be used in several apps.

  Rules:
  - call "exo create service [<name>] [<template>] [<model-name>] [<description>]" to scaffold a reusable service

  Scenario: create reusable service
    Given I am in the root directory of an empty application called "empty app"
    When executing "exo-create service --service-role=users --service-type=users-service --author=test-author --template-name=exoservice-es6-mongodb --model-name=user --description='manage users'"
    Then my application contains the file "application.yml" with the content:
      """
      name: empty app
      description: Empty test application
      version: 1.0.0

      bus:
        type: exocom
        version: 0.21.7

      services:
        public:
          users:
            location: ../users-service
      """
    And my workspace contains the file "../users-service/service.yml" with content:
      """
      type: users-service
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
          version: '3.4.0'
          docker_flags:
            volume: '-v {{EXO_DATA_PATH}}:/data/db'
            online_text: 'waiting for connections'
            port: '-p 27017:27017'
      """
