## Messages[list&delete]

### List:
- list messages in the mailbox, messages are returned in order of latest received to oldest.
- URL: ```api/v1/messages```
- Method: ```GET```

### query params:
| Parameter | Type    | Required | Description                |
|-----------|---------|----------|----------------------------|
| limit     | integer | false    | Limit results (default 50) |
| start     | integer | false    | Pagination offset          |

### Response:
- Status: 200
```json
{
  "total": 500,
  "unread": 500,
  "count": 50,
  "start": 0,
  "messages": [
    {
      "ID": "dbchds-7ab1-466f-8cee-2e1cf0fasdasd",
      "Read": false,
      "From": {
        "Name": "examplPerson",
        "Address": "exampleperson1@example.com"
      },
      "To": [
        {
          "Name": "Some Other",
          "Address": "example2@example.com"
        }
      ],
      "Cc": [
        {
          "Name": "Accounts",
          "Address": "accounts@example.com"
        }
      ],
      "Bcc": [],
      "Subject": "Message subject",
      "Created": "2022-10-03T21:35:32.228605299+13:00",
      "Size": 6144,
      "Attachments": 0
    },
  ]
}
```

## Delete Messages Functionalitiy:
- delete one or more messages by passing ID.
- URL: ```api/v1/messages```
- Method: ```DELETE```
- Request: 
```json
{
    "ids": ["<ID>","<ID>"]
}
```
- Response:
Status: ```200```

### Delete all messages
- delete all messages (deleting all messages):
- URL: ```api/v1/messages```
- Method: ```DELETE```

- Request:
```json
{
    "ids": []
}
```

- Response:
- Status:  ```200```

## Update individual read status:
- set the read status of one or more messages:
- URL: ```api/v1/messages`
- METHOD: ```PUT```
- Request:
```json
{
  "ids": ["<ID>", "<ID>"],
  "read": false
}
```
- Response: ```200```

## Update all messages read status:
- set the read status of all messages.
- URL: ```api/v1/messages```
- METHOD: ```PUT```
- Request:
```json
{
  "ids": [],
  "read": false
}
```
- Response:
Status: ```200```
