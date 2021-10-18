// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/Updatelistinfo": {
            "post": {
                "description": "格式：POST /api/v1/Updatelistinfo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "修改名单列表(会同步刷新缓存和推送waf)",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "WafLists",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.WafLists"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/addlistinfo": {
            "post": {
                "description": "格式：POST /api/v1/addlistinfo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "增加名单元素(会同步刷新缓存和推送waf)",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "WafList",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.WafList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/cachesync": {
            "get": {
                "description": "格式：get /api/v1/cachesync?pushwaf=yes|no\n参数pushwaf可选，默认为yes，即缓存同步redis后推送至waf。no为不推送",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "waf缓存同步",
                "parameters": [
                    {
                        "type": "string",
                        "description": "可选参数",
                        "name": "pushwaf",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/dellistinfo": {
            "post": {
                "description": "格式：POST /api/v1/dellistinfo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "删除名单元素(会同步刷新缓存和推送waf)",
                "parameters": [
                    {
                        "description": "参数",
                        "name": "WafList",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.WafList"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/disablelock": {
            "get": {
                "description": "格式：get /api/v1/disablelock?key=\n不能为空字段：key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "关锁",
                "parameters": [
                    {
                        "type": "string",
                        "description": "锁名称",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/disablelocklist": {
            "get": {
                "description": "格式：get /api/v1/disablelocklist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "关闭状态锁列表",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/enablelock": {
            "get": {
                "description": "格式：get /api/v1/enablelock?key=\n不能为空字段：key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "开锁",
                "parameters": [
                    {
                        "type": "string",
                        "description": "锁名称",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/enablelocklist": {
            "get": {
                "description": "格式：GET /api/v1/enablelocklist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "开启状态锁列表",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/getalllistinfo": {
            "get": {
                "description": "格式：get /api/v1/getalllistinfo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取名单列表",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/getlistinfo": {
            "get": {
                "description": "格式：get /api/v1/getlistinfo?key=\n不能为空字段：key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取单个名单列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "锁名称",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/getlistnamelists": {
            "get": {
                "description": "格式：get /api/v1/getlistnamelists",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取名单名称列表",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/wafclientlist": {
            "get": {
                "description": "格式：get /api/v1/wafclientlist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取在线waf client列表",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/wafclientstart": {
            "get": {
                "description": "格式：get /api/v1/wafclientstart?host=x.x.x.x:port\nhost参数不能为空，用于标识本机",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "客户端启动通知接口",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/waflog": {
            "post": {
                "description": "格式：POST /api/v1/waflog\n不能为空字段：key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "设置waf地址",
                "parameters": [
                    {
                        "description": "相关信息",
                        "name": "WafLog",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.WafLog"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        },
        "/api/v1/wafsync": {
            "get": {
                "description": "格式：get /api/v1/wafsync?host=xxx\nhost非必须参数，有参数时同步指定节点，无参数时同步所有节点\nhost列表获取接口：get /api/v1/wafclientlist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "waf客户端同步",
                "parameters": [
                    {
                        "type": "string",
                        "description": "主机",
                        "name": "host",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/main.Resp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Resp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "main.WafList": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "listinfo": {
                    "type": "string"
                },
                "listname": {
                    "type": "string"
                }
            }
        },
        "main.WafLists": {
            "type": "object",
            "properties": {
                "listinfo": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "listname": {
                    "type": "string"
                }
            }
        },
        "main.WafLog": {
            "type": "object",
            "properties": {
                "client": {
                    "type": "string"
                },
                "dtime": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
