Feature: Vertigo REST API Demo
  As a developer
  I want to interact with Vertigo via standard HTTP endpoints
  So that I can integrate it with web and mobile applications

  Background:
    Given the Persistence Facade is initialized
    And the database has a table "users" with rows:
      | id | name  | email             |
      | 1  | Alice | alice@example.com |
      | 2  | Bob   | bob@example.com   |

  Scenario: Fetching all users via GET /api/users
    When I send a GET request to "/api/users"
    Then the HTTP status code should be 200
    And the HTTP response body should contain "Alice"
    And the HTTP response body should contain "Bob"

  Scenario: Dispatching a raw SQL query via POST /api/dispatch
    When I send a POST request to "/api/dispatch" with body:
      """
      {
        "sql": "SELECT name FROM users WHERE id = 1",
        "channel": "test_rest"
      }
      """
    Then the HTTP status code should be 200
    And the HTTP response body should contain "Alice"
    And the HTTP response body should not contain "Bob"
