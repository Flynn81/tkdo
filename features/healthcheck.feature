Feature: health check
  In order to be happy
  As an api consumer
  I need to know the api is healthy

  Scenario: Health check is up
    Given I make a request to the health check
    When I get a response
    Then a 200 response code is returned And there is no body
