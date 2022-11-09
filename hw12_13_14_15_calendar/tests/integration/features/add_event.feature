# file: add_event.feature
Feature: add event
    As an application user
    I want to add event
    so that i can make notification for send

  Scenario: successfully create event
    When I send "POST" request to "/v1/events" with json body:
      """
       {
        "common_event" : {
           "title": "created title",
           "description": "created description",
           "startDate": "2022-01-15T00:00:00Z",
           "endDate": "2022-04-01T00:00:00Z",
           "ownerId":1,
           "notificationDate" : "2022-03-05T00:00:00Z"
        }
      }
      """
    Then the response code should be 200
    And the response body should be contains key: "inserted_id"

  Scenario: error event time busy
    When I send "POST" request to "/v1/events" with json body:
      """
        {
            "common_event" : {
              "title": "created title",
              "description": "created description",
              "startDate": "2022-01-01T00:00:00Z",
              "endDate": "2022-04-01T00:00:00Z",
              "ownerId":1,
              "notificationDate" : "2022-03-05T00:00:00Z"
          }
        }
      """
    Then the response code should be 400
    And the response should match json:
      """
        {
          "code": 3, "message": "date is already busy : 2022-01-01 00:00:00 +0000 UTC"
        }
      """

  Scenario: add event with other start date
    When I send "POST" request to "/v1/events" with json body:
      """
        {
          "common_event" : {
              "title": "created title",
              "description": "created description",
              "startDate": "2022-01-16T00:00:00Z",
              "endDate": "2022-04-01T00:00:00Z",
              "ownerId":1,
              "notificationDate" : "2022-03-05T00:00:00Z"
          }
        }
      """
    Then the response code should be 200
    And the response body should be contains key: "inserted_id"

  Scenario: wrong http method
    When I send "GET" request to "/v1/events" with json body:
      """
         {
           "event" : {
              "title": "created title",
              "description": "created description",
              "startDate": "2022-03-01T00:00:00Z",
              "endDate": "2022-04-01T00:00:00Z",
              "ownerId":1,
              "notificationDate" : "2022-03-05T00:00:00Z"
            }
        }
      """
    Then the response code should be 501
