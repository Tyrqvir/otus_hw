## file: list_event.feature
Feature: show list of events
  As an application user
  I want to see list of events for me
  so that i can analyze events

  Scenario: show result schema
    When I send "GET" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-01-01T23:59:00Z/owner/1"
    Then the response code should be 200
    And the response should match json:
    """
    {
      "events": [
          {
              "id": "1",
              "common_event": {
                  "title": "title 1",
                  "start_date": "2022-01-01T00:00:00Z",
                  "end_date": "2022-02-01T00:00:00Z",
                  "description": "description 1",
                  "owner_id": "1",
                  "notification_date": "2022-01-02T00:00:00Z"
              }
          }
      ]
  }
    """

  Scenario: events as list of day
    When I send "GET" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-01-02T00:00:00Z/owner/1"
    Then the response code should be 200
    And the response body has 2 items


  Scenario: events as list of week
    When I send "GET" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-01-07T00:00:00Z/owner/1"
    Then the response code should be 200
    And the response body has 7 items

  Scenario: events as list of month
    When I send "GET" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-02-01T00:00:00Z/owner/1"
    Then the response code should be 200
    And the response body has 9 items

  Scenario: events has flag as already sent
    When I wait 4s
    Then I send "GET" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-01-01T23:59:00Z/owner/1"
    And the response should match json:
    """
    {
        "events": [
            {
                "id": "1",
                "common_event": {
                    "title": "title 1",
                    "start_date": "2022-01-01T00:00:00Z",
                    "end_date": "2022-02-01T00:00:00Z",
                    "description": "description 1",
                    "owner_id": "1",
                    "is_notified": "1",
                    "notification_date": "2022-01-02T00:00:00Z"
                }
            }
        ]
    }
    """

  Scenario: wrong http method
    When I send "POST" request to "/v1/events/period/2022-01-01T00:00:00Z/2022-01-04T00:00:00Z/owner/1"
    Then the response code should be 501
