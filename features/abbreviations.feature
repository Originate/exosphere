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

  Scenario Outline: executing "exo clone" with abbreviated commands
    When executing the abbreviated command <CLONE-REPO-COMMAND> in the terminal
    Then the full command "exo clone" is executed
  Examples:
    | CLONE-REPO-COMMAND |
    | exo cl             |
    | exo clon           |

  Scenario Outline: executing "exo create application" with abbreviated commands
    When executing the abbreviated command <CREATE-APP-COMMAND> in the terminal
    Then the full command "exo create application" is executed
  Examples:
    | CREATE-APP-COMMAND |
    | exo cr a           |
    | exo creat app      |

  Scenario Outline: executing "exo create service" with abbreviated commands
    When executing the abbreviated command <CREATE-SERVICE-COMMAND> in the terminal
    Then the full command "exo create service" is executed
  Examples:
    | CREATE-SERVICE-COMMAND |
    | exo cr s               |
    | exo creat serv         |

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
    | exo se        |
    | exo setu      |

  Scenario Outline: executing "exo test" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command <TEST-COMMAND> in the terminal
    Then the full command "exo test" is executed
  Examples:
    | TEST-COMMAND  |
    | exo t         |
    | exo tes       |
