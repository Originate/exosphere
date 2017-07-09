@docker-cleanup
Feature: executing abbreviated Exosphere commands

  As a developer using Exosphere
  I want to be able to use abbreviations for commands
  So that I can use Exosphere without having to type so much.

  Scenario Outline: executing "exo add-template" with abbreviated commands
    When executing the abbreviated command "<ADD-TEMPLATE-COMMAND>" in the terminal
    Then the full command "exo add-template" is executed
  Examples:
    | ADD-TEMPLATE-COMMAND     |
    | exo add- name url        |
    | exo add-t name url       |


  Scenario Outline: executing "exo clean" with abbreviated commands
    When executing the abbreviated command "<CLONE-REPO-COMMAND>" in the terminal
    Then the full command "exo clean" is executed
  Examples:
    | CLONE-REPO-COMMAND |
    | exo cle            |
    | exo clea           |


  Scenario Outline: executing "exo clone" with abbreviated commands
    When executing the abbreviated command "<CLONE-REPO-COMMAND>" in the terminal
    Then the full command "exo clone" is executed
  Examples:
    | CLONE-REPO-COMMAND |
    | exo clo            |
    | exo clon           |

  Scenario Outline: executing "exo create" with abbreviated commands
    When executing the abbreviated command "<CREATE-APP-COMMAND>" in the terminal
    Then the full command "exo create" is executed
  Examples:
    | CREATE-APP-COMMAND    |
    | exo cr                |
    | exo creat             |

  Scenario Outline: executing "exo fetch-templates" with abbreviated commands
    When executing the abbreviated command "<FETCH-TEMPLATES-COMMAND>" in the terminal
    Then the full command "exo fetch-templates" is executed
  Examples:
    | FETCH-TEMPLATES-COMMAND    |
    | exo fe                     |
    | exo fetc                   |

  Scenario Outline: executing "exo remove-template" with abbreviated commands
    When executing the abbreviated command "<REMOVE-TEMPLATE-COMMAND>" in the terminal
    Then the full command "exo remove-template" is executed
  Examples:
    | REMOVE-TEMPLATE-COMMAND   |
    | exo re name               |
    | exo rem name              |

  Scenario Outline: executing "exo run" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command "<RUN-COMMAND>" in the terminal
    Then the full command "exo run" is executed
  Examples:
    | RUN-COMMAND  |
    | exo ru       |

  Scenario Outline: executing "exo setup" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command "<SETUP-COMMAND>" in the terminal
    Then the full command "exo setup" is executed
  Examples:
    | SETUP-COMMAND |
    | exo se        |
    | exo setu      |

  Scenario Outline: executing "exo test" with abbreviated commands
    Given a set-up "tests-passing" application
    When executing the abbreviated command "<TEST-COMMAND>" in the terminal
    Then the full command "exo test" is executed
  Examples:
    | TEST-COMMAND  |
    | exo t         |
    | exo tes       |
