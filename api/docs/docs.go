// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/auth/login": {
            "post": {
                "description": "Checks the user data and returns a jwt token on correct Login",
                "tags": [
                    "Auth"
                ],
                "summary": "Logs in the user",
                "parameters": [
                    {
                        "description": "Login details",
                        "name": "details",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginDetails"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Checks the user data and adds it to the repo",
                "tags": [
                    "Auth"
                ],
                "summary": "Registers the user",
                "parameters": [
                    {
                        "description": "Register details",
                        "name": "details",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterDetails"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/card": {
            "get": {
                "description": "Fetches all cards that match the query",
                "tags": [
                    "Card"
                ],
                "summary": "Fetch card by query",
                "parameters": [
                    {
                        "type": "string",
                        "name": "lang",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "name": "maxPrice",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "name": "minPrice",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCard"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new card",
                "tags": [
                    "Card"
                ],
                "summary": "Create new card",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "new card data",
                        "name": "card",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PostCard"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCard"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/card/{id}": {
            "get": {
                "description": "Fetches a card by it's id",
                "tags": [
                    "Card"
                ],
                "summary": "Fetch card by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Card ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCard"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Updates an existing card",
                "tags": [
                    "Card"
                ],
                "summary": "Update card",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "Card ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "new card data",
                        "name": "card",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PostCard"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCard"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/collection": {
            "post": {
                "description": "Creates a new card collection",
                "tags": [
                    "Collection"
                ],
                "summary": "Create new collection",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "new collection data",
                        "name": "collection",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PostCollection"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCollection"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/collection/all": {
            "get": {
                "description": "Fetches all the user's collections",
                "tags": [
                    "Collection"
                ],
                "summary": "Fetch all collections",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCollection"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/collection/{collectionId}": {
            "post": {
                "description": "Adds, removes or alters a card slot in an existing collection",
                "tags": [
                    "Collection"
                ],
                "summary": "Add, remove or alter card slot",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "Collection ID",
                        "name": "collectionId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "new card slot data",
                        "name": "collectionSlot",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PostCollectionSlot"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCollection"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        },
        "/collection/{id}": {
            "get": {
                "description": "Fetches a collection by it's id",
                "tags": [
                    "Collection"
                ],
                "summary": "Fetch collection by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Collection ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetCollection"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a collection by it's id",
                "tags": [
                    "Collection"
                ],
                "summary": "Delete collection",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Collection ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Deletes a collection's info by it's id",
                "tags": [
                    "Collection"
                ],
                "summary": "Update collection info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Collection ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authenticator",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "new collection data",
                        "name": "collection",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.PostCollection"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.ErrResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "dto.GetCard": {
            "type": "object",
            "properties": {
                "cardType": {
                    "$ref": "#/definitions/model.CardType"
                },
                "id": {
                    "type": "integer"
                },
                "imageUrl": {
                    "type": "string"
                },
                "language": {
                    "$ref": "#/definitions/model.Language"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dto.GetCollection": {
            "type": "object",
            "properties": {
                "cards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.GetCollectionSlot"
                    }
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "dto.GetCollectionSlot": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "card": {
                    "$ref": "#/definitions/dto.GetCard"
                }
            }
        },
        "dto.LoginDetails": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.PostCard": {
            "type": "object",
            "required": [
                "language",
                "name",
                "price",
                "text",
                "type"
            ],
            "properties": {
                "imageUrl": {
                    "type": "string"
                },
                "language": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "text": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "dto.PostCollection": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "minLength": 3
                }
            }
        },
        "dto.PostCollectionSlot": {
            "type": "object",
            "required": [
                "amount",
                "cardId"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "cardId": {
                    "type": "integer"
                }
            }
        },
        "dto.RegisterDetails": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 4
                }
            }
        },
        "model.CardType": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "longName": {
                    "type": "string"
                }
            }
        },
        "model.Language": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "longName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Card store api",
	Description:      "Service for buying/selling collectable cards",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
