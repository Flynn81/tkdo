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

  Scenario: Updating a task
    Given the api is up
    When a user is created
    And the user creates a task
    And the user updates the task
    Then the api returns an ok status code

  Scenario: Deleting a task
    Given the api is up
    When a user is created
    And the user creates a task
    And the user deletes the task
    Then the api returns an ok status code

  Scenario: Getting mutliple tasks
    Given the api is up
    When a user is created
    And the user creates multiple tasks
    And then requests their tasks
    Then multiple tasks are returned

  Scenario: Requesting a page of tasks
    Given the api is up
    When a user is created
    And the user creates twenty tasks
    And then requests their tasks with page size of five
    Then only five tasks are returned

  Scenario: Requesting a second page of tasks
    Given the api is up
    When a user is created
    And the user creates twenty tasks
    And then requests their tasks with page size of five
    And saves the returned requests
    And then requests a second page of tasks with page size of five
    Then only five tasks are returned
    And the five tasks are different from the first five

  Scenario: Requesting a page of tasks with page size greater than available tasks
    Given the api is up
    When a user is created
    And the user creates twenty tasks
    And then requests their tasks with page size of fifty
    Then all twenty tasks are returned

  Scenario: Requesting a page of tasks beyond what is available
    Given the api is up
    When a user is created
    And the user creates twenty tasks
    And then requests their tasks with page size of five and page five
    Then zero tasks are returned
