# Update Exsiting Application

Update the Application Name 

**URL** : `/api/application/{application_token}`

**Method** : `POST`

## Success Response

**Code** : `200 OK`

**Content examples**

For an application with right token saved in DB this will be the output

```json
{
    "SuccessMessage": "App Updated Successfully"
}
```

For an application with a wrong token this will be the output

```json
{
  "ErrorMessage": "Please check your App Token"
}
```
