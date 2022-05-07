# Show All Messages for specific Chat

Getall messages for specific chat providing chat number and application token as a URL parameters

**URL** : `/api/application/{application_token}/{chat_number}/messages`

**Method** : `GET`

## Success Response

**Code** : `200 OK`

**Content examples**

For a valid application token and chat number the response should be as follows

```json
[
    {
        "MessageBody": "Hello World"
    },
    {
        "MessageBody": "Hello World 1"
    }
]
```

For an unvalid application token and chat number the response should be as follows

```json
{
    "ErrorMessage": "Please check your App Token and Chat Number"
}
```
