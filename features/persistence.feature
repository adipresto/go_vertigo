Feature: Vertigo Persistence Handler
  As a developer
  I want a simplified interface to query the database and broadcast results
  So that I don't have to manage connection pools and network resilience manually

  Scenario: Dispatching a SQL query to the database and network
    Given the Persistence Facade is initialized
    And the database has a table "users" with rows:
      | id | name  | email             |
      | 1  | Alice | alice@example.com |
      | 2  | Bob   | bob@example.com   |
    When I dispatch the SQL query "SELECT id, name, email FROM users"
    Then the result should be streamed to the network
    And the network payload should contain "Alice" and "Bob"
