# Create Chat

**URL** : `/api/application/{application_token}/chat`

**Method** : `POST`

## Success Response

**Code** : `200 OK`

**Content examples**

For an Application with correct token should recieve response like follows

```json
{
  	"ChatNumber": 1
}
```

For an Application with wrong token should recieve response like follows

```json
{
	"ErrorMessage": "Please check your App Token"
}
```
