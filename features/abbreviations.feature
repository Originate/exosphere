Feature: executing abbreviated Exosphere commands

  As a developer using Exosphere
  I want to be able to use abbreviations for commands
  So that I can use Exosphere without having to type so much.

  Scenario Outline: executing "exo add" with abbreviated commands
    When executing the abbreviated command <ADD-COMMAND> in the terminal
    Then the full command "exo add" is executed
  Examples:
    | ADD-COMMAND |
    | exo a       |
    | exo ad      |

  Scenario Outline: executing "exo create" with abbreviated commands
    When executing the abbreviated command <CREATE-COMMAND> in the terminal
    Then the full command "exo create" is executed
  Examples:
    | CREATE-COMMAND |
    | exo cr         |
    | exo creat      |

  Scenario Outline: executing "exo run" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command <RUN-COMMAND> in the terminal
    Then the full command "exo run" is executed
  Examples:
    | RUN-COMMAND |
    | exo r       |
    | exo ru      |

  Scenario Outline: executing "exo setup" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command <SETUP-COMMAND> in the terminal
    Then the full command "exo setup" is executed
  Examples:
    | SETUP-COMMAND |
    | exo s         |
    | exo setu      |

  Scenario Outline: executing "exo test" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command <TEST-COMMAND> in the terminal
    Then the full command "exo test" is executed
  Examples:
    | TEST-COMMAND  |
    | exo t         |
    | exo tes       |
    
