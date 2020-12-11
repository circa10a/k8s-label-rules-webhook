// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-12-19 10:22:13.231925 -0600 CST m=+0.071816604

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
        "contact": {
            "name": "GitHub",
            "url": "https://github.com/circa10a/k8s-label-rules-webhook/"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/circa10a/k8s-label-rules-webhook/blob/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "post": {
                "description": "Respond to k8s with approval or rejection of labels compared against the ruleset",
                "produces": [
                    "application/json"
                ],
                "summary": "Respond to k8s with approval or rejection of labels compared against the ruleset",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.webhookResponse"
                        }
                    }
                }
            }
        },
        "/reload": {
            "post": {
                "description": "Hot reload ruleset file into memory without downtime",
                "produces": [
                    "application/json"
                ],
                "summary": "Reload ruleset file",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.reloadRulesResponse"
                            }
                        }
                    }
                }
            }
        },
        "/rules": {
            "get": {
                "description": "Show current ruleset being used to validate labels against",
                "produces": [
                    "application/json"
                ],
                "summary": "Display loaded rules",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.rule"
                            }
                        }
                    }
                }
            }
        },
        "/validate": {
            "get": {
                "description": "Validate regex of every rule in ruleset, respond with rule errors",
                "produces": [
                    "application/json"
                ],
                "summary": "Validate regex of all loaded rules",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.validRulesResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.reloadRulesResponse": {
            "type": "object",
            "properties": {
                "newRules": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.rule"
                    }
                },
                "reloaded": {
                    "type": "boolean"
                },
                "ruleErr": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.ruleError"
                    }
                },
                "yamlErr": {
                    "type": "string"
                }
            }
        },
        "main.response": {
            "type": "object",
            "properties": {
                "allowed": {
                    "type": "boolean"
                },
                "status": {
                    "type": "object",
                    "$ref": "#/definitions/main.status"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "main.rule": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "object",
                    "$ref": "#/definitions/main.value"
                }
            }
        },
        "main.ruleError": {
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "rulename": {
                    "type": "string"
                }
            }
        },
        "main.status": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "main.validRulesResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.ruleError"
                    }
                },
                "rulesValid": {
                    "type": "boolean"
                }
            }
        },
        "main.value": {
            "type": "object",
            "properties": {
                "regex": {
                    "type": "string"
                }
            }
        },
        "main.webhookResponse": {
            "type": "object",
            "properties": {
                "apiVersion": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                },
                "response": {
                    "type": "object",
                    "$ref": "#/definitions/main.response"
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
	Version:     "0.2.9",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "k8s-label-rules-webhook",
	Description: "A kubernetes webhook to standardize labels on resources",
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
