Feature: create user
  In order to use the api
  As a user
  I need to create a user

  Scenario: User creation
    Given the api is up
    When I create a user
    Then get back the user I sent with an id
