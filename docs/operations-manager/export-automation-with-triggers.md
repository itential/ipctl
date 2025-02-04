{
  "message": "Successfully exported automation",
  "data": {
    "_id": "672ede2975ae358dee04e8d3",
    "name": "test-with-triggers",
    "description": "This is a test automation",
    "componentName": "test",
    "componentType": "workflows",
    "gbac": {
      "write": [],
      "read": []
    },
    "created": "2024-11-09T03:59:37.091Z",
    "createdBy": "admin@pronghorn",
    "lastUpdated": "2024-11-09T03:59:37.091Z",
    "lastUpdatedBy": "admin@pronghorn",
    "triggers": [
      {
        "_id": "672ede2975ae358dee04e8d6",
        "created": "2024-11-09T03:59:37.097Z",
        "createdBy": "admin@pronghorn",
        "lastUpdated": "2024-11-09T03:59:54.856Z",
        "lastUpdatedBy": "admin@pronghorn",
        "name": "test-api-trigger",
        "type": "endpoint",
        "enabled": true,
        "actionType": "automations",
        "actionId": "672ede2975ae358dee04e8d3",
        "description": "This is a test api trigger",
        "verb": "POST",
        "routeName": "test-with-triggers",
        "schema": {
          "type": "object",
          "properties": {},
          "additionalProperties": true
        },
        "jst": {
          "_id": "6716dafdd1a3e0115acac3a8",
          "name": "@6716dafd113f9679380359e0: Process Instance Response",
          "description": "",
          "incoming": [
            {
              "$id": "instanceResponse",
              "properties": {
                "response": {
                  "properties": {
                    "DescribeInstancesResponse": {
                      "properties": {
                        "reservationSet": {
                          "properties": {
                            "item": {
                              "items": {
                                "properties": {
                                  "instancesSet": {
                                    "properties": {
                                      "item": {
                                        "items": {
                                          "properties": {
                                            "instanceId": {
                                              "examples": [
                                                "i-091629e3a896331c2"
                                              ],
                                              "type": "string"
                                            }
                                          },
                                          "required": [],
                                          "type": "object"
                                        },
                                        "type": "array"
                                      }
                                    },
                                    "required": [],
                                    "type": "object"
                                  }
                                },
                                "required": [],
                                "type": "object"
                              },
                              "type": "array"
                            }
                          },
                          "required": [],
                          "type": "object"
                        }
                      },
                      "required": [],
                      "type": "object"
                    }
                  },
                  "required": [],
                  "type": "object"
                }
              },
              "required": [],
              "type": "object"
            }
          ],
          "outgoing": [
            {
              "$id": "instanceIdArray",
              "type": "array"
            }
          ],
          "steps": [
            {
              "context": "#",
              "from": {
                "location": "incoming",
                "name": "instanceResponse",
                "ptr": "/response/DescribeInstancesResponse/reservationSet/item"
              },
              "id": 5,
              "to": {
                "location": "method",
                "name": 4,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "id": 4,
              "type": "method",
              "library": "Array",
              "method": "flatMap",
              "args": [
                null,
                "extractAndFlattenInstances"
              ],
              "view": {
                "col": 1,
                "row": 1
              },
              "context": "#"
            },
            {
              "context": "#",
              "from": {
                "location": "method",
                "name": 4,
                "ptr": "/return"
              },
              "id": 7,
              "to": {
                "location": "method",
                "name": 6,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "id": 6,
              "type": "method",
              "library": "Array",
              "method": "map",
              "args": [
                null,
                "extractInstanceIds"
              ],
              "view": {
                "col": 2,
                "row": 1
              },
              "context": "#"
            },
            {
              "context": "#",
              "from": {
                "location": "method",
                "name": 6,
                "ptr": "/return"
              },
              "id": 8,
              "to": {
                "location": "outgoing",
                "name": "instanceIdArray",
                "ptr": ""
              },
              "type": "assign"
            }
          ],
          "functions": [
            {
              "incoming": [
                {
                  "$id": "currentValue",
                  "properties": {
                    "instancesSet": {
                      "properties": {
                        "item": {
                          "items": {
                            "properties": {
                              "instanceId": {
                                "examples": [
                                  "i-091629e3a896331c2"
                                ],
                                "type": "string"
                              }
                            },
                            "required": [],
                            "type": "object"
                          },
                          "type": "array"
                        }
                      },
                      "required": [],
                      "type": "object"
                    }
                  },
                  "required": [],
                  "type": "object"
                },
                {
                  "$id": "index",
                  "optional": true,
                  "title": "index",
                  "type": "number"
                },
                {
                  "$id": "array",
                  "items": {
                    "properties": {
                      "instancesSet": {
                        "properties": {
                          "item": {
                            "items": {
                              "properties": {
                                "instanceId": {
                                  "examples": [
                                    "i-091629e3a896331c2"
                                  ],
                                  "type": "string"
                                }
                              },
                              "required": [],
                              "type": "object"
                            },
                            "type": "array"
                          }
                        },
                        "required": [],
                        "type": "object"
                      }
                    },
                    "required": [],
                    "type": "object"
                  },
                  "optional": true,
                  "type": "array"
                }
              ],
              "outgoing": [
                {
                  "$id": "newValue",
                  "title": "newValue",
                  "type": [
                    "array",
                    "boolean",
                    "number",
                    "integer",
                    "string",
                    "object",
                    "null"
                  ]
                }
              ],
              "steps": [
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": "/instancesSet/item"
                  },
                  "id": 1,
                  "to": {
                    "location": "outgoing",
                    "name": "newValue",
                    "ptr": ""
                  },
                  "type": "assign"
                }
              ],
              "functions": [],
              "name": "extractAndFlattenInstances",
              "view": {
                "col": 2,
                "row": 4
              },
              "id": "extractAndFlattenInstances",
              "comments": []
            },
            {
              "incoming": [
                {
                  "$id": "currentValue",
                  "type": [
                    "array",
                    "boolean",
                    "number",
                    "integer",
                    "string",
                    "object",
                    "null"
                  ]
                },
                {
                  "$id": "index",
                  "optional": true,
                  "title": "index",
                  "type": "number"
                },
                {
                  "$id": "array",
                  "optional": true,
                  "type": "array"
                }
              ],
              "outgoing": [
                {
                  "$id": "newValue",
                  "editable": true,
                  "title": "newValue",
                  "type": [
                    "array",
                    "boolean",
                    "number",
                    "integer",
                    "string",
                    "object",
                    "null"
                  ]
                }
              ],
              "steps": [
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": ""
                  },
                  "id": 2,
                  "to": {
                    "location": "method",
                    "name": 1,
                    "ptr": "/args/0/value"
                  },
                  "type": "assign"
                },
                {
                  "id": 1,
                  "type": "method",
                  "library": "Object",
                  "method": "getProperty",
                  "args": [
                    null,
                    "instanceId"
                  ],
                  "view": {
                    "col": 1,
                    "row": 1
                  },
                  "context": "#"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "method",
                    "name": 1,
                    "ptr": "/return"
                  },
                  "id": 3,
                  "to": {
                    "location": "outgoing",
                    "name": "newValue",
                    "ptr": ""
                  },
                  "type": "assign"
                }
              ],
              "functions": [],
              "name": "extractInstanceIds",
              "view": {
                "col": 1,
                "row": 4
              },
              "id": "extractInstanceIds",
              "comments": []
            }
          ],
          "comments": [],
          "view": {
            "col": 2,
            "row": 5
          },
          "created": "2024-11-08T14:36:17.114Z",
          "createdBy": {
            "_id": "668c58df4f234baee4996cfb",
            "provenance": "local_aaa",
            "username": "admin@pronghorn"
          },
          "lastUpdated": "2024-11-08T14:36:17.180Z",
          "lastUpdatedBy": {
            "_id": "668c58df4f234baee4996cfb",
            "provenance": "local_aaa",
            "username": "admin@pronghorn"
          },
          "version": "4.3.6-2023.2.2",
          "tags": [],
          "namespace": {
            "type": "project",
            "_id": "6716dafd113f9679380359e0",
            "name": "AWS EC2",
            "accessControl": {
              "manage": [
                "account:668c58df4f234baee4996cfb"
              ],
              "write": [
                "account:668c58df4f234baee4996cfb"
              ],
              "execute": [
                "account:668c58df4f234baee4996cfb"
              ],
              "read": [
                "account:668c58df4f234baee4996cfb"
              ]
            }
          }
        },
        "migrationVersion": 3
      },
      {
        "_id": "672ede2975ae358dee04e8d7",
        "created": "2024-11-09T03:59:37.097Z",
        "createdBy": "admin@pronghorn",
        "lastUpdated": "2024-11-09T03:59:37.097Z",
        "lastUpdatedBy": "admin@pronghorn",
        "name": "test-event-trigger",
        "type": "eventSystem",
        "enabled": true,
        "actionType": "automations",
        "actionId": "672ede2975ae358dee04e8d3",
        "description": "this is a test event trigger",
        "source": "@itential/app-lifecycle_manager",
        "topic": "actionComplete",
        "schema": {
          "type": "object",
          "properties": {},
          "additionalProperties": true
        },
        "jst": {
          "_id": "6716dafdd1a3e0115acac3a1",
          "name": "@6716dafd113f9679380359e0: Build EC2 Tag Data",
          "description": "",
          "incoming": [
            {
              "$id": "vpcName",
              "type": "string"
            },
            {
              "$id": "launchResponse",
              "properties": {
                "response": {
                  "properties": {
                    "RunInstancesResponse": {
                      "properties": {
                        "instancesSet": {
                          "properties": {
                            "item": {
                              "items": {
                                "properties": {
                                  "instanceId": {
                                    "examples": [
                                      "i-0e5a940bc1fc03f1a"
                                    ],
                                    "type": "string"
                                  }
                                },
                                "required": [],
                                "type": "object"
                              },
                              "type": "array"
                            }
                          },
                          "required": [],
                          "type": "object"
                        }
                      },
                      "required": [],
                      "type": "object"
                    }
                  },
                  "required": [],
                  "type": "object"
                }
              },
              "required": [],
              "type": "object"
            }
          ],
          "outgoing": [
            {
              "$id": "ec2Tags",
              "type": "array"
            },
            {
              "$id": "instanceIdArray",
              "type": "array"
            }
          ],
          "steps": [
            {
              "context": "#",
              "from": {
                "location": "incoming",
                "name": "vpcName",
                "ptr": ""
              },
              "id": 26,
              "to": {
                "location": "method",
                "name": 25,
                "ptr": "/args/1/value"
              },
              "type": "assign"
            },
            {
              "context": "#",
              "from": {
                "location": "incoming",
                "name": "launchResponse",
                "ptr": "/response/RunInstancesResponse/instancesSet/item"
              },
              "id": 32,
              "to": {
                "location": "method",
                "name": 31,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "id": 25,
              "type": "method",
              "library": "String",
              "method": "concat",
              "args": [
                "Apache Web Server for ",
                null
              ],
              "view": {
                "col": 1,
                "row": 1
              },
              "context": "#"
            },
            {
              "id": 31,
              "type": "method",
              "library": "Array",
              "method": "getIndex",
              "args": [
                null,
                0
              ],
              "view": {
                "col": 1,
                "row": 2
              },
              "context": "#"
            },
            {
              "context": "#",
              "from": {
                "location": "method",
                "name": 25,
                "ptr": "/return"
              },
              "id": 27,
              "to": {
                "location": "function",
                "name": 21,
                "ptr": "/args/1/value"
              },
              "type": "assign"
            },
            {
              "context": "#",
              "from": {
                "location": "method",
                "name": 31,
                "ptr": "/return"
              },
              "id": 34,
              "to": {
                "location": "method",
                "name": 33,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "id": 21,
              "type": "function",
              "function": "buildEc2Tags",
              "args": [
                "Name",
                ""
              ],
              "view": {
                "col": 2,
                "row": 1
              }
            },
            {
              "id": 33,
              "type": "method",
              "library": "Object",
              "method": "getProperty",
              "args": [
                null,
                "instanceId"
              ],
              "view": {
                "col": 2,
                "row": 2
              },
              "context": "#"
            },
            {
              "context": "#",
              "from": {
                "location": "function",
                "name": 21,
                "ptr": "/return/ec2Tags"
              },
              "id": 23,
              "to": {
                "location": "declaration",
                "name": 22,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "context": "#",
              "from": {
                "location": "method",
                "name": 33,
                "ptr": "/return"
              },
              "id": 36,
              "to": {
                "location": "declaration",
                "name": 35,
                "ptr": "/args/0/value"
              },
              "type": "assign"
            },
            {
              "id": 22,
              "type": "declaration",
              "library": "Array",
              "method": "new Array",
              "args": [
                null
              ],
              "view": {
                "col": 3,
                "row": 1
              },
              "context": "#",
              "polymorphIndex": 0
            },
            {
              "id": 35,
              "type": "declaration",
              "library": "Array",
              "method": "new Array",
              "args": [
                null
              ],
              "view": {
                "col": 3,
                "row": 2
              },
              "context": "#",
              "polymorphIndex": 0
            },
            {
              "context": "#",
              "from": {
                "location": "declaration",
                "name": 22,
                "ptr": "/return"
              },
              "id": 24,
              "to": {
                "location": "outgoing",
                "name": "ec2Tags",
                "ptr": ""
              },
              "type": "assign"
            },
            {
              "context": "#",
              "from": {
                "location": "declaration",
                "name": 35,
                "ptr": "/return"
              },
              "id": 37,
              "to": {
                "location": "outgoing",
                "name": "instanceIdArray",
                "ptr": ""
              },
              "type": "assign"
            }
          ],
          "functions": [
            {
              "incoming": [
                {
                  "$id": "currentValue",
                  "properties": {
                    "port": {
                      "examples": [
                        443
                      ],
                      "type": "integer"
                    },
                    "protocol": {
                      "examples": [
                        "TCP"
                      ],
                      "type": "string"
                    },
                    "sourceIp": {
                      "examples": [
                        "1%2E2%2E3%2E4"
                      ],
                      "format": "ipv4",
                      "type": "string"
                    },
                    "sourceSubnetCidrMask": {
                      "examples": [
                        32
                      ],
                      "type": "integer"
                    }
                  },
                  "required": [],
                  "type": "object"
                },
                {
                  "$id": "index",
                  "optional": true,
                  "title": "index",
                  "type": "number"
                },
                {
                  "$id": "array",
                  "items": {
                    "properties": {
                      "port": {
                        "examples": [
                          443
                        ],
                        "type": "integer"
                      },
                      "protocol": {
                        "examples": [
                          "TCP"
                        ],
                        "type": "string"
                      },
                      "sourceIp": {
                        "examples": [
                          "1%2E2%2E3%2E4"
                        ],
                        "format": "ipv4",
                        "type": "string"
                      },
                      "sourceSubnetCidrMask": {
                        "examples": [
                          32
                        ],
                        "type": "integer"
                      }
                    },
                    "required": [],
                    "type": "object"
                  },
                  "optional": true,
                  "type": "array"
                }
              ],
              "outgoing": [
                {
                  "$id": "newValue",
                  "editable": true,
                  "title": "newValue",
                  "type": [
                    "array",
                    "boolean",
                    "number",
                    "integer",
                    "string",
                    "object",
                    "null"
                  ]
                }
              ],
              "steps": [
                {
                  "id": 1,
                  "type": "function",
                  "function": "buildRule",
                  "args": [
                    null,
                    null,
                    null,
                    null
                  ],
                  "view": {
                    "col": 1,
                    "row": 1
                  },
                  "context": "#"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": "/sourceSubnetCidrMask"
                  },
                  "id": 2,
                  "to": {
                    "location": "function",
                    "name": 1,
                    "ptr": "/args/0/value"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": "/port"
                  },
                  "id": 3,
                  "to": {
                    "location": "function",
                    "name": 1,
                    "ptr": "/args/1/value"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": "/sourceIp"
                  },
                  "id": 4,
                  "to": {
                    "location": "function",
                    "name": 1,
                    "ptr": "/args/2/value"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "currentValue",
                    "ptr": "/protocol"
                  },
                  "id": 5,
                  "to": {
                    "location": "function",
                    "name": 1,
                    "ptr": "/args/3/value"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "function",
                    "name": 1,
                    "ptr": "/return/rule"
                  },
                  "id": 6,
                  "to": {
                    "location": "outgoing",
                    "name": "newValue",
                    "ptr": ""
                  },
                  "type": "assign"
                }
              ],
              "functions": [],
              "name": "buildRuleList",
              "view": {
                "col": 2,
                "row": 4
              },
              "id": "buildRuleList",
              "comments": []
            },
            {
              "incoming": [
                {
                  "$id": "sourceSubnetCidrMask",
                  "type": "number"
                },
                {
                  "$id": "port",
                  "type": "number"
                },
                {
                  "$id": "sourceIp",
                  "type": "string"
                },
                {
                  "$id": "protocol",
                  "type": "string"
                }
              ],
              "outgoing": [
                {
                  "$id": "rule",
                  "properties": {
                    "port": {
                      "examples": [
                        443
                      ],
                      "type": "number"
                    },
                    "protocol": {
                      "examples": [
                        "TCP"
                      ],
                      "type": "string"
                    },
                    "sourceIp": {
                      "examples": [
                        "192%2E168%2E30%2E10"
                      ],
                      "format": "ipv4",
                      "type": "string"
                    },
                    "sourceSubnetCidrMask": {
                      "examples": [
                        32
                      ],
                      "type": "number"
                    }
                  },
                  "required": [],
                  "type": "object"
                }
              ],
              "steps": [
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "sourceSubnetCidrMask",
                    "ptr": ""
                  },
                  "id": 1,
                  "to": {
                    "location": "outgoing",
                    "name": "rule",
                    "ptr": "/sourceSubnetCidrMask"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "port",
                    "ptr": ""
                  },
                  "id": 2,
                  "to": {
                    "location": "outgoing",
                    "name": "rule",
                    "ptr": "/port"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "sourceIp",
                    "ptr": ""
                  },
                  "id": 3,
                  "to": {
                    "location": "outgoing",
                    "name": "rule",
                    "ptr": "/sourceIp"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "protocol",
                    "ptr": ""
                  },
                  "id": 4,
                  "to": {
                    "location": "outgoing",
                    "name": "rule",
                    "ptr": "/protocol"
                  },
                  "type": "assign"
                }
              ],
              "functions": [],
              "name": "buildRule",
              "view": {
                "col": 2,
                "row": 4
              },
              "id": "buildRule",
              "comments": []
            },
            {
              "incoming": [
                {
                  "$id": "Key",
                  "type": "string"
                },
                {
                  "$id": "Value",
                  "type": "string"
                }
              ],
              "outgoing": [
                {
                  "$id": "ec2Tags",
                  "properties": {
                    "Key": {
                      "examples": [
                        "Name"
                      ],
                      "type": "string"
                    },
                    "Value": {
                      "examples": [
                        "ec2"
                      ],
                      "type": "string"
                    }
                  },
                  "required": [],
                  "type": "object"
                }
              ],
              "steps": [
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "Key",
                    "ptr": ""
                  },
                  "id": 1,
                  "to": {
                    "location": "outgoing",
                    "name": "ec2Tags",
                    "ptr": "/Key"
                  },
                  "type": "assign"
                },
                {
                  "context": "#",
                  "from": {
                    "location": "incoming",
                    "name": "Value",
                    "ptr": ""
                  },
                  "id": 2,
                  "to": {
                    "location": "outgoing",
                    "name": "ec2Tags",
                    "ptr": "/Value"
                  },
                  "type": "assign"
                }
              ],
              "functions": [],
              "name": "buildEc2Tags",
              "view": {
                "col": 1,
                "row": 4
              },
              "id": "buildEc2Tags",
              "comments": []
            }
          ],
          "comments": [],
          "view": {
            "col": 3,
            "row": 8
          },
          "created": "2024-11-08T14:36:17.083Z",
          "createdBy": {
            "_id": "668c58df4f234baee4996cfb",
            "provenance": "local_aaa",
            "username": "admin@pronghorn"
          },
          "lastUpdated": "2024-11-08T14:36:17.180Z",
          "lastUpdatedBy": {
            "_id": "668c58df4f234baee4996cfb",
            "provenance": "local_aaa",
            "username": "admin@pronghorn"
          },
          "version": "4.3.6-2023.2.2",
          "tags": [],
          "namespace": {
            "type": "project",
            "_id": "6716dafd113f9679380359e0",
            "name": "AWS EC2",
            "accessControl": {
              "manage": [
                "account:668c58df4f234baee4996cfb"
              ],
              "write": [
                "account:668c58df4f234baee4996cfb"
              ],
              "execute": [
                "account:668c58df4f234baee4996cfb"
              ],
              "read": [
                "account:668c58df4f234baee4996cfb"
              ]
            }
          }
        },
        "legacyWrapper": false,
        "migrationVersion": 3
      },
      {
        "_id": "672ede2975ae358dee04e8d5",
        "created": "2024-11-09T03:59:37.097Z",
        "createdBy": "admin@pronghorn",
        "lastUpdated": "2024-12-25T11:59:23.253Z",
        "lastUpdatedBy": "admin@pronghorn",
        "name": "test-manual-trigger",
        "type": "manual",
        "enabled": true,
        "actionType": "automations",
        "actionId": "672ede2975ae358dee04e8d3",
        "description": "this is a test manual trigger",
        "formId": "cli-devel-test",
        "formData": {},
        "formSchemaHash": "a41642a90e3fbc500b5096e6dd82579aa5df6b0d4ca070b491a91f0c9c63c09c",
        "legacyWrapper": false,
        "migrationVersion": 3
      },
      {
        "_id": "672ede2975ae358dee04e8d4",
        "created": "2024-11-09T03:59:37.096Z",
        "createdBy": "admin@pronghorn",
        "lastUpdated": "2024-12-25T11:59:13.359Z",
        "lastUpdatedBy": "admin@pronghorn",
        "name": "test-schedule-trigger",
        "type": "schedule",
        "enabled": true,
        "actionType": "automations",
        "actionId": "672ede2975ae358dee04e8d3",
        "description": "this is a test schedule trigger",
        "formId": "cli-devel-test",
        "formData": {},
        "formSchemaHash": "a41642a90e3fbc500b5096e6dd82579aa5df6b0d4ca070b491a91f0c9c63c09c",
        "legacyWrapper": false,
        "firstRunAt": 1730122369158,
        "processMissedRuns": "none",
        "repeatUnit": "day",
        "repeatFrequency": 1,
        "repeatInterval": 86400000,
        "migrationVersion": 3
      }
    ]
  },
  "metadata": {}
}
