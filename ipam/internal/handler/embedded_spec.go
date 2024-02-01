// Code generated by go-swagger; DO NOT EDIT.

package handler

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
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "IP address management",
    "title": "IPAM service",
    "version": "1.0.0"
  },
  "paths": {
    "/acquire": {
      "get": {
        "summary": "Get new IPv6 range",
        "parameters": [
          {
            "enum": [
              64,
              127
            ],
            "type": "integer",
            "description": "Prefix length",
            "name": "prefixlen",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "Prefix sucessfully acquired",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/release": {
      "delete": {
        "summary": "Release IPv6 range",
        "parameters": [
          {
            "maxLength": 41,
            "minLength": 14,
            "type": "string",
            "description": "Prefix",
            "name": "prefix",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              64,
              127
            ],
            "type": "integer",
            "description": "Prefix length",
            "name": "prefixlen",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Prefix sucessfully released"
          },
          "400": {
            "description": "Bad request"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "IP address management",
    "title": "IPAM service",
    "version": "1.0.0"
  },
  "paths": {
    "/acquire": {
      "get": {
        "summary": "Get new IPv6 range",
        "parameters": [
          {
            "enum": [
              64,
              127
            ],
            "type": "integer",
            "description": "Prefix length",
            "name": "prefixlen",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "Prefix sucessfully acquired",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    },
    "/release": {
      "delete": {
        "summary": "Release IPv6 range",
        "parameters": [
          {
            "maxLength": 41,
            "minLength": 14,
            "type": "string",
            "description": "Prefix",
            "name": "prefix",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              64,
              127
            ],
            "type": "integer",
            "description": "Prefix length",
            "name": "prefixlen",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Prefix sucessfully released"
          },
          "400": {
            "description": "Bad request"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    }
  }
}`))
}
