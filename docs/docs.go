// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Idris",
            "url": "https://github.com/supuwoerc",
            "email": "zhangzhouou@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/public/user/login": {
            "post": {
                "description": "用于用户登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "注册参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "10000": {
                        "description": "操作成功",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    },
                    "10001": {
                        "description": "操作失败",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    },
                    "10002": {
                        "description": "参数错误",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    }
                }
            }
        },
        "/api/v1/public/user/signup": {
            "post": {
                "description": "用于用户注册帐号",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "注册参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "10000": {
                        "description": "操作成功",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    },
                    "10001": {
                        "description": "操作失败",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-string"
                        }
                    },
                    "10002": {
                        "description": "参数错误",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/profile": {
            "get": {
                "description": "获取用户账户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户信息",
                "responses": {
                    "10000": {
                        "description": "操作成功",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    },
                    "10001": {
                        "description": "操作失败",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    },
                    "10002": {
                        "description": "参数错误",
                        "schema": {
                            "$ref": "#/definitions/response.BasicResponse-any"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.SignUpRequest": {
            "type": "object",
            "required": [
                "code",
                "email",
                "id",
                "password"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "response.BasicResponse-any": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "response.BasicResponse-string": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Learn-Gin-Web API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
