POST /integrations

REQUEST

```
{
  "properties": {
    "name": "v3",
    "type": "Adapter",
    "properties": {
      "id": "v3",
      "type": "Swagger Petstore:1.0.0"
    },
    "virtual": true
  }
}
```

RESPONSE

```
{
  "status": "Created",
  "message": "Successfully created v3",
  "data": {
    "_id": "676c001652dda825325cc68a",
    "name": "v3",
    "model": "@itential/adapter_Swagger Petstore:1.0.0",
    "type": "Adapter",
    "properties": {
      "id": "v3",
      "type": "Swagger Petstore:1.0.0",
      "brokers": [],
      "groups": [],
      "properties": {
        "authentication": {},
        "server": {
          "protocol": "http",
          "host": "petstore.swagger.io",
          "base_path": "/v1"
        },
        "tls": {
          "enabled": false,
          "rejectUnauthorized": true
        },
        "version": "1.0.0"
      }
    },
    "isEncrypted": true,
    "loggerProps": {
      "description": "Logging",
      "log_max_files": 10,
      "log_max_file_size": 10485760,
      "log_level": "info",
      "log_directory": "/opt/itential/logs",
      "log_filename": "itential.log",
      "console_level": "info",
      "syslog": {
        "level": "warning",
        "host": "127.0.0.1",
        "port": 514,
        "protocol": "udp4",
        "facility": "local0",
        "type": "BSD",
        "path": "",
        "pid": "process.pid",
        "localhost": "",
        "app_name": "",
        "eol": ""
      }
    },
    "virtual": true
  }
}
```

