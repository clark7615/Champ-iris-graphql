package graphql

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

var subscribers sync.Map

type connectionACKMessage struct {
	*websocket.Conn
	OperationID string         `json:"id,omitempty"`
	Type        string         `json:"type"`
	Payload     requestOptions `json:"payload,omitempty"`
}

//PostQuery Graphql API
//進入點 POST http://Host:Port/Party/query
//基於iris mvc 架構
func (g *Graphql) PostQuery(Ctx iris.Context) {
	params := g.createParams(g.getRequestOptions(Ctx))
	res := graphql.Do(params)
	_, _ = Ctx.JSON(res)
}

//GetSubscriptions Graphql Apollo Server Subscriptions websocket
//進入點 GET ws://Host:Port/Party/subscription
//基於iris mvc 架構
func (g *Graphql) GetSubscription(Ctx iris.Context) {
	conn, err := upgrader.Upgrade(Ctx.ResponseWriter(), Ctx.Request(), nil)
	if err != nil {
		fmt.Printf("failed to do websocket upgrade: %v", err)
		return
	}
	connectionACK, err := json.Marshal(map[string]string{
		"type": "connection_ack",
	})
	if err != nil {
		fmt.Printf("failed to marshal ws connection ack: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
		fmt.Printf("failed to write to ws connection: %v", err)
		return
	}
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				return
			}
			if err != nil {
				fmt.Println("failed to read websocket message: %v", err)
				return
			}
			var msg connectionACKMessage
			msg.Conn = conn
			if err := json.Unmarshal(p, &msg); err != nil {
				fmt.Printf("failed to unmarshal: %v", err)
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

//GetPg Graphql Playground 靜態網頁
//進入點 GET http://Host:Port/Party/pg
func (g *Graphql) GetPg(Ctx iris.Context) {
	if !g.ShowPlayground {
		Ctx.StatusCode(404)
		return
	}
	t := template.New("Playground")
	te, _ := t.Parse(html)
	endpoint := strings.Replace(Ctx.RequestPath(true), "pg", "query", 1)
	subEndpoint := strings.Replace(endpoint, "query", "subscription", 1)
	_ = te.ExecuteTemplate(Ctx.ResponseWriter(), "index", struct {
		Endpoint             string
		SubscriptionEndpoint string
	}{
		Endpoint:             endpoint,
		SubscriptionEndpoint: subEndpoint,
	})
}
