# Search Messages Model

This an integration endpoint with elasticsearch to search over messages model

**URL** : `/api/search/{message}`

**Method** : `POST`

## Success Response

**Code** : `200 OK`

**Content examples**

For a valid search keyword as "he" the response should be as follows 

```json
[
  {
    "chat_id":6,
    "id":6,
    "message_body":"Hello World"
  }
]
```

