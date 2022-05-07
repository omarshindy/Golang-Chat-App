# Creeate Message

**URL** : `/api/application/{application_token}/{chat_number}/message`

**Method** : `POST`

## Success Response

**Code** : `200 OK`

**Content examples**

For an Application with a valid token and valid chat token the respose should be as follows 

```json
{
    "SuccessMessage": "Message Created Successfully"
}
```

For an Application with an unvalid token and an unvalid chat token the respose should be as follows 

```json
{
      "ErrorMessage": "Please check your App Token and Chat Number"
}
```
