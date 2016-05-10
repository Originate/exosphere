Feature: Unknown command

  As an Exosphere user accidentally entering an unknown Exosphere command
  I want to be given a useful error message and a list of available commands
  So that I can realize my mistake and remember the correct name for the command I wanted to run.


  Scenario: the user enters an unknown command
    When trying to run "exo zonk" in the terminal
    Then I see
      """
      Error: unknown command 'zonk'

      Usage: exo <command> [options]

      Available commands are:
      * create
      * install
      * run
      """
