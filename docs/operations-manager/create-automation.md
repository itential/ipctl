POST /operations-manager/automations

Request

```
{
  "name": "UI-create-automation",
  "description": "this will create an automation",
  "componentType": "workflows"
}
```


Response

```
{
  "message": "Successfully created UI-create-automation",
  "data": {
    "_id": "676ab26f5f359ee5cbf41172",
    "name": "UI-create-automation",
    "description": "this will create an automation",
    "componentType": "workflows",
    "componentId": null,
    "gbac": {
      "write": [],
      "read": []
    },
    "created": "2024-12-24T13:09:03.356Z",
    "createdBy": "668c58df4f234baee4996cfb",
    "lastUpdated": "2024-12-24T13:09:03.356Z",
    "lastUpdatedBy": "668c58df4f234baee4996cfb"
  },
  "metadata": {}
}
```
