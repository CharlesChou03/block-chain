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
        "/blocks": {
            "get": {
                "description": "get latest n block information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Block Chain Information"
                ],
                "summary": "get latest n block information",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/models.GetLatestNBlockRes"
                        }
                    },
                    "204": {
                        "description": "no content"
                    },
                    "500": {
                        "description": "internal error"
                    }
                }
            }
        },
        "/blocks/{id}": {
            "get": {
                "description": "get block information by block number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Block Chain Information"
                ],
                "summary": "get block information by block number",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/models.GetBlockByNumRes"
                        }
                    },
                    "204": {
                        "description": "no content"
                    },
                    "500": {
                        "description": "internal error"
                    }
                }
            }
        },
        "/health": {
            "get": {
                "summary": "health checker API",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/transaction/{txHash}": {
            "get": {
                "description": "get transaction information by hash",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Block Chain Information"
                ],
                "summary": "get transaction information by hash",
                "parameters": [
                    {
                        "type": "string",
                        "description": "txHash",
                        "name": "txHash",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/models.GetTransactionByHashRes"
                        }
                    },
                    "204": {
                        "description": "no content"
                    },
                    "500": {
                        "description": "internal error"
                    }
                }
            }
        },
        "/version": {
            "get": {
                "summary": "service version API",
                "responses": {
                    "200": {
                        "description": "0.0.1",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.BlockInfo": {
            "type": "object",
            "properties": {
                "block_hash": {
                    "type": "string"
                },
                "block_num": {
                    "type": "integer"
                },
                "block_time": {
                    "type": "integer"
                },
                "parent_hash": {
                    "type": "string"
                }
            }
        },
        "models.GetBlockByNumRes": {
            "type": "object",
            "properties": {
                "block_hash": {
                    "type": "string"
                },
                "block_num": {
                    "type": "integer"
                },
                "block_time": {
                    "type": "integer"
                },
                "parent_hash": {
                    "type": "string"
                },
                "transactions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.GetLatestNBlockRes": {
            "type": "object",
            "properties": {
                "blocks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.BlockInfo"
                    }
                }
            }
        },
        "models.GetTransactionByHashRes": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "from": {
                    "type": "string"
                },
                "logs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Log"
                    }
                },
                "nonce": {
                    "type": "integer"
                },
                "to": {
                    "type": "string"
                },
                "tx_hash": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "models.Log": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "index": {
                    "type": "integer"
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
	Version:     "0.0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Swagger",
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