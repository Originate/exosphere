Feature: help command

  As a developer not familiar with how to call "exo add" correctly
  I want to be given helpful information how to use this application
  So that I can do the things I intended without having to look this up somewhere else.


  Scenario: the user enters 'exo add help'
    When running "exo-add help" in the terminal
    Then I see:
    """
    Usage: exo add [<entity-name>]

    Adds a new service to the current application.
    This command must be called in the root directory of the application.

    options: [<service-role>] [<service-type>] [<template>] [<model>] [<description>]
    """
