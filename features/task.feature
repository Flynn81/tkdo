Feature: tasks
  In order to keep track of things
  As a person with things to do
  I need to put my things in one place

  Scenario: No tasks no user
    Given the api is up
    When tasks are requested
    Then I get back an error

  Scenario: No tasks
    Given the api is up
    When tasks are requested by a user
    Then I get back an empty array

  Scenario: Create a task
    Given the api is up
    When a task is created
    Then I get back created status

  Scenario: There is a task
    Given the api is up
    When a user is created
    And the user creates a task
    And tasks are requested by the user
    Then the one task is returned
