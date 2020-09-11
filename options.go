package graphql

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

type requestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
	UUID          string                 `json:"UUID,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"graphql-ws"},
}

func (g *Graphql) newSchema() graphql.Schema {
	s, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:        g.Query.Obj,
		Mutation:     g.Mutation.Obj,
		Subscription: g.Subscription.Obj,
	})
	return s
}

func (g *Graphql) createParams(opt *requestOptions) graphql.Params {
	return graphql.Params{
		Schema:         g.newSchema(),
		RequestString:  opt.Query,
		VariableValues: opt.Variables,
		OperationName:  opt.OperationName,
	}
}

func (g *Graphql) getRequestOptions(Ctx iris.Context) *requestOptions {
	r := Ctx.Request()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var opts requestOptions
		_ = Ctx.ReadJSON(&opts)
		return &opts
	}

	operationsParam := r.FormValue("operations")
	var opts requestOptions
	if err := json.Unmarshal([]byte(operationsParam), &opts); err != nil {
		return &requestOptions{}
	}

	mapParam := r.FormValue("map")
	mapValues := make(map[string][]string)
	if len(mapParam) != 0 {
		if err := json.Unmarshal([]byte(mapParam), &mapValues); err != nil {
			return &requestOptions{}
		}
	}

	variables := opts

	for key, value := range mapValues {
		for _, v := range value {
			if file, header, err := r.FormFile(key); err == nil {

				var node interface{} = variables

				parts := strings.Split(v, ".")
				last := parts[len(parts)-1]

				for _, vv := range parts[:len(parts)-1] {
					switch node.(type) {
					case requestOptions:
						if vv == "variables" {
							node = opts.Variables
						} else {
							return &requestOptions{}
						}
					case map[string]interface{}:
						node = node.(map[string]interface{})[vv]
					case []interface{}:
						if idx, err := strconv.ParseInt(vv, 10, 64); err == nil {
							node = node.([]interface{})[idx]
						} else {
							return &requestOptions{}
						}
					default:
						return &requestOptions{}
					}
				}

				data := &MultipartFile{File: file, Header: header}

				switch node.(type) {
				case map[string]interface{}:
					node.(map[string]interface{})[last] = data
				case []interface{}:
					if idx, err := strconv.ParseInt(last, 10, 64); err == nil {
						node.([]interface{})[idx] = data
					} else {
						return &requestOptions{}
					}
				default:
					return &requestOptions{}
				}
			}
		}
	}
	return &opts
}
