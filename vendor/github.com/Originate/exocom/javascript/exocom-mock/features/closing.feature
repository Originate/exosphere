Feature: Closing the instance

  As a developer usin ExoComMock in my tests
  I want to be able to shut it down clearly and completely
  So that I can set up a fresh instance for each test.

  Rules:
  - call "close" on your ExoComMock instance to make it remove all its side effects


  Scenario: closing an active instance
    Given an ExoComMock instance
    When closing it
    Then it is no longer listening


  Scenario: closing an inactive instance
    Given an ExoComMock instance
    Then I can close it without errors
