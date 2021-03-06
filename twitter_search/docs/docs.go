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
        "/v1/tweet/search": {
            "get": {
                "description": "Search tweets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tweets"
                ],
                "summary": "Get raw tweet based on search",
                "operationId": "get-tweets",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search",
                        "name": "search",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/search.SearchTweetResponse"
                        }
                    }
                }
            }
        },
        "/v1/tweet/searchByLang": {
            "get": {
                "description": "Search tweets by lang, accepted language:\nName\tLanguage code\nEnglish (default)\ten\nArabic\tar\nBengali\tbn\nCzech\tcs\nDanish\tda\nGerman\tde\nGreek\tel\nSpanish\tes\nPersian\tfa\nFinnish\tfi\nFilipino\tfil\nFrench\tfr\nHebrew\the\nHindi\thi\nHungarian\thu\nIndonesian\tid\nItalian\tit\nJapanese\tja\nKorean\tko\nMalay\tmsa\nDutch\tnl\nNorwegian\tno\nPolish\tpl\nPortuguese\tpt\nRomanian\tro\nRussian\tru\nSwedish\tsv\nThai\tth\nTurkish\ttr\nUkrainian\tuk\nUrdu\tur\nVietnamese\tvi\nChinese (Simplified)\tzh-cn\nChinese (Traditional)\tzh-tw",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tweets"
                ],
                "summary": "Get raw tweet based on search",
                "operationId": "get-tweets-byLang",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search by lang",
                        "name": "search",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/search.SearchTweetResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "search.EsTweetObj": {
            "type": "object",
            "properties": {
                "Tweet": {
                    "$ref": "#/definitions/search.TweetStruct"
                },
                "TweetId": {
                    "type": "string"
                }
            }
        },
        "search.SearchTweetResponse": {
            "type": "object",
            "properties": {
                "tweets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/search.EsTweetObj"
                    }
                }
            }
        },
        "search.TweetStruct": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "lang": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "text": {
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
