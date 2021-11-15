// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Calendar api",
    "version": "1.0.0"
  },
  "paths": {
    "/meet/create": {
      "post": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/MeetInfo"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/MeetCreateResponse"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/meet/info": {
      "get": {
        "parameters": [
          {
            "type": "string",
            "name": "meeting_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/MeetInfo"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "Not found"
          }
        }
      }
    },
    "/ping": {
      "get": {
        "responses": {
          "200": {
            "description": "Ok"
          }
        }
      }
    },
    "/users/create": {
      "post": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserInfo"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/UsersCreateResponse"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/users/info": {
      "get": {
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query"
          },
          {
            "type": "string",
            "name": "email",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/UserInfo"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ErrorResponse": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "MeetCreateResponse": {
      "type": "object",
      "required": [
        "meet_id"
      ],
      "properties": {
        "meet_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "MeetInfo": {
      "type": "object",
      "required": [
        "name",
        "creator",
        "time_start",
        "time_end",
        "visibility"
      ],
      "properties": {
        "creator": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "meeting_link": {
          "type": "string"
        },
        "meeting_room": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "notifications": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Notification"
          }
        },
        "participants": {
          "type": "array",
          "minItems": 1,
          "items": {
            "type": "string"
          }
        },
        "repeat": {
          "type": "string",
          "enum": [
            "day",
            "week",
            "month",
            "year",
            "workday"
          ]
        },
        "time_end": {
          "type": "string",
          "format": "date-time"
        },
        "time_start": {
          "type": "string",
          "format": "date-time"
        },
        "visibility": {
          "type": "string",
          "enum": [
            "all",
            "participants"
          ]
        }
      },
      "additionalProperties": false
    },
    "Notification": {
      "type": "object",
      "required": [
        "before_start",
        "notification_type"
      ],
      "properties": {
        "before_start": {
          "description": "Minutes before meeting start",
          "type": "integer"
        },
        "notification_type": {
          "type": "string",
          "enum": [
            "email",
            "sms",
            "telegram"
          ]
        }
      },
      "additionalProperties": false
    },
    "UserInfo": {
      "type": "object",
      "required": [
        "email",
        "phone"
      ],
      "properties": {
        "day_end": {
          "type": "string",
          "example": "20:00"
        },
        "day_start": {
          "type": "string",
          "example": "10:00"
        },
        "email": {
          "type": "string"
        },
        "first_name": {
          "type": "string"
        },
        "last_name": {
          "type": "string"
        },
        "phone": {
          "type": "string"
        },
        "user_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "UsersCreateResponse": {
      "type": "object",
      "required": [
        "user_id"
      ],
      "properties": {
        "user_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Calendar api",
    "version": "1.0.0"
  },
  "paths": {
    "/meet/create": {
      "post": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/MeetInfo"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/MeetCreateResponse"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/meet/info": {
      "get": {
        "parameters": [
          {
            "type": "string",
            "name": "meeting_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/MeetInfo"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "Not found"
          }
        }
      }
    },
    "/ping": {
      "get": {
        "responses": {
          "200": {
            "description": "Ok"
          }
        }
      }
    },
    "/users/create": {
      "post": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserInfo"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/UsersCreateResponse"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/users/info": {
      "get": {
        "parameters": [
          {
            "type": "string",
            "name": "phone",
            "in": "query"
          },
          {
            "type": "string",
            "name": "email",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Ok",
            "schema": {
              "$ref": "#/definitions/UserInfo"
            }
          },
          "400": {
            "description": "Creation failed",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ErrorResponse": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "MeetCreateResponse": {
      "type": "object",
      "required": [
        "meet_id"
      ],
      "properties": {
        "meet_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "MeetInfo": {
      "type": "object",
      "required": [
        "name",
        "creator",
        "time_start",
        "time_end",
        "visibility"
      ],
      "properties": {
        "creator": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "meeting_link": {
          "type": "string"
        },
        "meeting_room": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "notifications": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Notification"
          }
        },
        "participants": {
          "type": "array",
          "minItems": 1,
          "items": {
            "type": "string"
          }
        },
        "repeat": {
          "type": "string",
          "enum": [
            "day",
            "week",
            "month",
            "year",
            "workday"
          ]
        },
        "time_end": {
          "type": "string",
          "format": "date-time"
        },
        "time_start": {
          "type": "string",
          "format": "date-time"
        },
        "visibility": {
          "type": "string",
          "enum": [
            "all",
            "participants"
          ]
        }
      },
      "additionalProperties": false
    },
    "Notification": {
      "type": "object",
      "required": [
        "before_start",
        "notification_type"
      ],
      "properties": {
        "before_start": {
          "description": "Minutes before meeting start",
          "type": "integer"
        },
        "notification_type": {
          "type": "string",
          "enum": [
            "email",
            "sms",
            "telegram"
          ]
        }
      },
      "additionalProperties": false
    },
    "UserInfo": {
      "type": "object",
      "required": [
        "email",
        "phone"
      ],
      "properties": {
        "day_end": {
          "type": "string",
          "example": "20:00"
        },
        "day_start": {
          "type": "string",
          "example": "10:00"
        },
        "email": {
          "type": "string"
        },
        "first_name": {
          "type": "string"
        },
        "last_name": {
          "type": "string"
        },
        "phone": {
          "type": "string"
        },
        "user_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "UsersCreateResponse": {
      "type": "object",
      "required": [
        "user_id"
      ],
      "properties": {
        "user_id": {
          "type": "string"
        }
      },
      "additionalProperties": false
    }
  }
}`))
}
