# Search

**URL** : `api/v1/search?query=<string>`

**Method** : `GET`

The search returns the most recent matches (default 50).
Matching messages are returned in the order of latest received to oldest.


## Query parameters

| Parameter | Type    | Required | Description                |
|-----------|---------|----------|----------------------------|
| query     | string  | true     | Search query               |
| limit     | integer | false    | Limit results (default 50) |
| start     | integer | false    | Pagination offset          |


## Response

**Status** : `200`

```json
{
  "total": 500,
  "unread": 500,
  "count": 25,
  "start": 0,
  "messages": [
    {
      "ID": "ascasd-70ba-466f-8cee-sda",
      "Read": false,
      "From": {
        "Name": "TestPerson1",
        "Address": "testperson1@example.com"
      },
      "To": [
        {
          "Name": "Test",
          "Address": "test@example.com"
        }
      ],
      "Cc": [
        {
          "Name": "Accounts",
          "Address": "accounts@example.com"
        }
      ],
      "Bcc": [],
      "Subject": "Test email",
      "Created": "2022-10-03T21:35:32.228605299+13:00",
      "Size": 6144,
      "Attachments": 0
    },
  ]
}
```