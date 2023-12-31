// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
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
        "/api/v1/advert": {
            "get": {
                "description": "广告列表接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "广告列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Data"
                        }
                    }
                }
            }
        },
        "/api/v1/advert/{id}": {
            "get": {
                "description": "广告详情接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "广告详情",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Data"
                        }
                    }
                }
            }
        },
        "/api/v1/captcha": {
            "get": {
                "description": "获取登陆验证码接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Data"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.Data": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "自定义返回码  0:表示正常",
                    "type": "integer"
                },
                "data": {
                    "description": "返回的具体数据"
                },
                "message": {
                    "description": "关于这次响应的说明信息",
                    "type": "string"
                },
                "meta": {
                    "description": "数据meta"
                },
                "namespace": {
                    "description": "异常的范围",
                    "type": "string"
                },
                "reason": {
                    "description": "异常原因",
                    "type": "string"
                },
                "recommend": {
                    "description": "推荐链接",
                    "type": "string"
                },
                "request_id": {
                    "description": "请求Id",
                    "type": "string"
                },
                "type": {
                    "description": "数据类型, 可以缺省",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}
