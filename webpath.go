package graphql

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

var subscribers sync.Map

type ConnectionACKMessage struct {
	*websocket.Conn
	OperationID string         `json:"id,omitempty"`
	Type        string         `json:"type"`
	Payload     requestOptions `json:"payload,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"graphql-ws"},
}

func (g *Graphql) PostQuery(Ctx iris.Context) {
	params := g.createParams(Ctx)
	res := graphql.Do(params)
	_, _ = Ctx.JSON(res)
}

func (g *Graphql) GetSubscriptions(Ctx iris.Context) {
	conn, err := upgrader.Upgrade(Ctx.ResponseWriter(), Ctx.Request(), nil)
	if err != nil {
		log.Printf("failed to do websocket upgrade: %v", err)
		return
	}
	connectionACK, err := json.Marshal(map[string]string{
		"type": "connection_ack",
	})
	if err != nil {
		log.Printf("failed to marshal ws connection ack: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
		log.Printf("failed to write to ws connection: %v", err)
		return
	}
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				return
			}
			if err != nil {
				log.Println("failed to read websocket message: %v", err)
				return
			}
			var msg ConnectionACKMessage
			msg.Conn = conn
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Printf("failed to unmarshal: %v", err)
				return
			}
			if msg.Type == "start" {
				length := 0
				subscribers.Range(func(key, value interface{}) bool {
					length++
					return true
				})
				subscribers.Store(msg.OperationID, &msg)
			}
			if msg.Type == "stop" {
				subscribers.Delete(msg.OperationID)
			}

		}
	}()
}

func (g *Graphql) GetPg(Ctx iris.Context) {
	if !g.ShowPlayground {
		Ctx.StatusCode(404)
		return
	}
	t := template.New("Playground")
	te, _ := t.Parse(html)
	path := Ctx.RequestPath(true)
	path = strings.Replace(path, "pg", "query", 1)
	_ = te.ExecuteTemplate(Ctx.ResponseWriter(), "index", struct {
		Endpoint             string
		SubscriptionEndpoint string
	}{
		Endpoint:             path,
		SubscriptionEndpoint: path,
	})
}

func (g *Graphql) createParams(Ctx iris.Context) graphql.Params {
	opt := g.getRequestOptions(Ctx)
	return graphql.Params{
		Schema:         g.newSchema(),
		RequestString:  opt.Query,
		VariableValues: opt.Variables,
		OperationName:  opt.OperationName,
	}
}
