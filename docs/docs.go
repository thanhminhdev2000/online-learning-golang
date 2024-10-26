// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/forgot-password": {
            "post": {
                "description": "Send a password reset link to the user's email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Request password reset",
                "parameters": [
                    {
                        "description": "User email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Log in using email or username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Log in",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Log out by clearing the refresh token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Log out",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    }
                }
            }
        },
        "/auth/refresh-token": {
            "post": {
                "description": "Refresh the access token using the refresh token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Refresh access token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AccessTokenResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "description": "Reset the user's password using a valid token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Reset user password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Reset token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "New password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ResetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/contacts/": {
            "post": {
                "description": "Send email contact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contact"
                ],
                "summary": "Send email contact",
                "parameters": [
                    {
                        "description": "Send email",
                        "name": "contact",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Contact"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/documents/": {
            "get": {
                "description": "Trả về danh sách các tài liệu, có thể lọc theo ` + "`" + `subjectId` + "`" + ` và ` + "`" + `title` + "`" + `. Giới hạn số lượng tài liệu trả về bằng tham số ` + "`" + `limit` + "`" + `.",
                "tags": [
                    "Document"
                ],
                "summary": "Lấy danh sách tài liệu",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 40,
                        "description": "Giới hạn số lượng tài liệu trả về",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "ID của môn học",
                        "name": "subjectId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Tiêu đề của tài liệu (tìm kiếm bằng LIKE)",
                        "name": "title",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.DocumentsResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/documents/subjects": {
            "get": {
                "description": "List of classes with their subjects and document counts",
                "tags": [
                    "Document"
                ],
                "summary": "List of classes with their subjects and document counts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ClassWithSubjects"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/documents/upload": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Upload document file",
                "tags": [
                    "Document"
                ],
                "summary": "Upload document file",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Subject ID",
                        "name": "subjectId",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Document title",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/users/": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieve a list of all users, with optional filters for email, username, full name, date of birth, role, and pagination.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "string",
                        "name": "dateOfBirth",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "fullName",
                        "in": "query"
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "role",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "username",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Register a new user with email, username, and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/users/admin": {
            "post": {
                "description": "Register a new user with email, username, and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/users/{userId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieve user information by user ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update user information by user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user information",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserDetail"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UpdateUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Delete a user by user ID",
                "tags": [
                    "User"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/users/{userId}/avatar": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update the avatar for a specific user",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user avatar",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "User Avatar",
                        "name": "avatar",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UpdateUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/users/{userId}/password": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Change the user's password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Change user password",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Password data",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PasswordUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AccessTokenResponse": {
            "type": "object",
            "required": [
                "accessToken",
                "expiresIn"
            ],
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expiresIn": {
                    "type": "integer"
                }
            }
        },
        "models.ClassWithSubjects": {
            "type": "object",
            "required": [
                "classId",
                "className",
                "count",
                "subjects"
            ],
            "properties": {
                "classId": {
                    "type": "integer"
                },
                "className": {
                    "type": "string"
                },
                "count": {
                    "type": "integer"
                },
                "subjects": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SubjectId"
                    }
                }
            }
        },
        "models.Contact": {
            "type": "object",
            "required": [
                "content",
                "email",
                "fullName",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.CreateUserRequest": {
            "type": "object",
            "required": [
                "avatar",
                "dateOfBirth",
                "email",
                "fullName",
                "gender",
                "password",
                "username"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "dateOfBirth": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "gender": {
                    "$ref": "#/definitions/models.UserGender"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "role": {
                    "$ref": "#/definitions/models.UserRole"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.DocumentsResponse": {
            "type": "object",
            "required": [
                "author",
                "category",
                "documentType",
                "downloads",
                "fileUrl",
                "id",
                "title",
                "views"
            ],
            "properties": {
                "author": {
                    "type": "string"
                },
                "category": {
                    "type": "string"
                },
                "documentType": {
                    "type": "string"
                },
                "downloads": {
                    "type": "integer"
                },
                "fileUrl": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "views": {
                    "type": "integer"
                }
            }
        },
        "models.Error": {
            "type": "object",
            "required": [
                "error"
            ],
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.ForgotPasswordRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "models.LoginRequest": {
            "type": "object",
            "required": [
                "identifier",
                "password"
            ],
            "properties": {
                "identifier": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.LoginResponse": {
            "type": "object",
            "required": [
                "accessToken",
                "expiresIn",
                "message",
                "user"
            ],
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expiresIn": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/models.UserDetail"
                }
            }
        },
        "models.Message": {
            "type": "object",
            "required": [
                "message"
            ],
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.PagingInfo": {
            "type": "object",
            "required": [
                "limit",
                "page",
                "totalCount"
            ],
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "models.PasswordUpdateRequest": {
            "type": "object",
            "required": [
                "currentPassword",
                "newPassword"
            ],
            "properties": {
                "currentPassword": {
                    "type": "string",
                    "minLength": 6
                },
                "newPassword": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "models.ResetPasswordRequest": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "models.SubjectId": {
            "type": "object",
            "required": [
                "count",
                "subjectId",
                "subjectName"
            ],
            "properties": {
                "count": {
                    "type": "integer"
                },
                "subjectId": {
                    "type": "integer"
                },
                "subjectName": {
                    "type": "string"
                }
            }
        },
        "models.UpdateUserResponse": {
            "type": "object",
            "required": [
                "message",
                "user"
            ],
            "properties": {
                "message": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/models.UserDetail"
                }
            }
        },
        "models.UserDetail": {
            "type": "object",
            "required": [
                "avatar",
                "dateOfBirth",
                "email",
                "fullName",
                "gender",
                "id",
                "role",
                "username"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "dateOfBirth": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "gender": {
                    "$ref": "#/definitions/models.UserGender"
                },
                "id": {
                    "type": "integer"
                },
                "role": {
                    "$ref": "#/definitions/models.UserRole"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.UserGender": {
            "type": "string",
            "enum": [
                "female",
                "male"
            ],
            "x-enum-varnames": [
                "GenderFemale",
                "GenderMale"
            ]
        },
        "models.UserResponse": {
            "type": "object",
            "required": [
                "data",
                "paging"
            ],
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UserDetail"
                    }
                },
                "paging": {
                    "$ref": "#/definitions/models.PagingInfo"
                }
            }
        },
        "models.UserRole": {
            "type": "string",
            "enum": [
                "user",
                "admin"
            ],
            "x-enum-varnames": [
                "RoleUser",
                "RoleAdmin"
            ]
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Online Learning API",
	Description:      "This is an online learning API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
