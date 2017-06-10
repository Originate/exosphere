Feature: Before-All hook

  As a developer using external resources that have to be initialized before going online
  I want to be able to signal to ExoService when I am done initializing
  So that the service doesn't signal availability before everything is properly initialized.

  Rules:
  - put initializiation code into a "beforeAll" message handler


  Background:
    Given an ExoCom instance


  Scenario: code with beforeAll hook
    When starting the "example service with before-all hook" service
    Then it runs the "before-all" hook
