# returns the OpenAPI document

GET /integration-models/Swagger%20Petstore%20-%20OpenAPI%203.0%3A1.0.20-SNAPSHOT/export


# returns the integraiton model

GET /integration-models/Swagger%20Petstore%20-%20OpenAPI%203.0%3A1.0.20-SNAPSHOT

Request

None

Response

```
{
  "model": "@itential/adapter_Swagger Petstore - OpenAPI 3.0:1.0.20-SNAPSHOT",
  "versionId": "Swagger Petstore - OpenAPI 3.0:1.0.20-SNAPSHOT",
  "description": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [http://swagger.io](http://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)",
  "properties": {
    "authentication": {
      "petstore_auth": {
        "token": {
          "access_token": ""
        }
      },
      "api_key": {
        "value": ""
      }
    },
    "server": {
      "protocol": "http",
      "host": "/v3",
      "base_path": ""
    },
    "tls": {
      "enabled": false,
      "rejectUnauthorized": true
    },
    "version": "1.0.20-SNAPSHOT"
  }
}
```
