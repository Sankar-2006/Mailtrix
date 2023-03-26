# Message:

## Message Title:
- returns json title of the message and attachment
- Request API Url: ```api/v1/message/<ID>```
- Request Method: ```GET```

## Response:
- Status: ```200```
```json
{
  "ID": "d123415-1234-478b-123c-asd12gd12",
  "Read": true,
  "From": {
    "Name": "Some One",
    "Address": "some@example.com"
  },
  "To": [
    {
      "Name": "Another person",
      "Address": "another@example.com"
    }
  ],
  "Cc": [],
  "Bcc": [],
  "Subject": "Message subject",
  "Date": "2016-09-07T16:46:00+13:00",
  "Text": "Plain text MIME part of the email",
  "HTML": "HTML MIME part (if exists)",
  "Size": 79499,
  "Inline": [
    {
      "PartID": "1.2",
      "FileName": "filename.gif",
      "ContentType": "image/gif",
      "ContentID": "919564503@07092006-1525",
      "Size": 7760
    }
  ],
  "Attachments": [
    {
      "PartID": "2",
      "FileName": "filename.doc",
      "ContentType": "application/msword",
      "ContentID": "",
      "Size": 43520
    }
  ]
}
```

## Attachments:
- Request API URL: ```api/v1/message/<ID>/part/<PartID>```
- Request Method: ```GET```
- Return the attachment using MIME Type ```ContentType```

## Headers:
- Request API URL: ```api/v1/message/<ID>/headers```
- Request Method: ```GET```
- Returns all message headers as json output.
```json
{
  "Content-Type": [
    "multipart/related; type=\"multipart/alternative\"; boundary=\"----=_NextPart_000_0013_01C6A60C.47EEAB80\""
  ],
  "Date": [
    "Wed, 30 Dec 2005 23:38:30 +1200"
  ],
  "Delivered-To": [
    "user@example.com",
    "user-alias@example.com"
  ],
  "From": [
    "\"User Name\" \\u003remote@example.com\\u003e"
  ],
  "Message-Id": [
    "\\u003c001701c6a5a7$b3205580$0201010a@HomeOfficeSM\\u003e"
  ],
}
```

## Raw email:
- URL: ```api/v1/message/<ID>/raw```
- Request Method: ```GET```
- Returns email source including headers and attachments.