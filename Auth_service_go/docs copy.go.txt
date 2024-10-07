package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "summary": "Handle login request",
                "description": "Handle login by verifying username and nonce",
                "tags": ["auth"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "login",
                        "description": "Login data",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "username": {
                                    "type": "string"
                                },
                                "half_nonce": {
                                    "type": "string"
                                }
                            },
                            "required": ["username", "half_nonce"]
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error_code": {
                                    "type": "integer"
                                },
                                "error_message": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error_code": {
                                    "type": "integer"
                                },
                                "error_message": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/verify-challenge": {
            "post": {
                "summary": "Handle verify challenge request",
                "description": "Handle verify challenge by verifying full_nonce and challenge_response",
                "tags": ["auth"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "verify",
                        "description": "Challenge data",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "full_nonce": {
                                    "type": "string"
                                },
                                "challenge_response": {
                                    "type": "string"
                                }
                            },
                            "required": ["full_nonce", "challenge_response"]
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error_code": {
                                    "type": "integer"
                                },
                                "error_message": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error_code": {
                                    "type": "integer"
                                },
                                "error_message": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
