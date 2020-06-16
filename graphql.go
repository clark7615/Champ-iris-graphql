package graphql

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

type commandType int

type Graphql struct {
	Query          rootType
	Mutation       rootType
	Subscription   rootType
	ShowPlayground bool
}

const (
	Query            commandType = iota
	Mutation         commandType = iota
	Subscription     commandType = iota
	QueryAndMutation commandType = iota
	All              commandType = iota
)

var ql *Graphql

//Default 建立一個 Graphql 實體並包含 Query以及Mutation RootType
func Default() *Graphql {
	return New(QueryAndMutation)
}

//New 依據CommandType建立一個Graphql 實體並新增對應之RootType
func New(commandType commandType) *Graphql {
	if ql == nil {
		ql = new(Graphql)
	}
	switch commandType {
	case Query:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
	case Mutation:
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	case Subscription:
		ql.Subscription.new("Subscription", "訂閱相關的命令")
		go runSubscription(ql)
	case All:
		New(Subscription)
		fallthrough
	case QueryAndMutation:
		New(Query)
		New(Mutation)
	}
	return ql
}

func runSubscription(g *Graphql) {
	for {
		time.Sleep(2 * time.Second)
		subscribers.Range(func(key, value interface{}) bool {
			msg, ok := value.(*connectionACKMessage)
			if !ok {
				return true
			}
			res := graphql.Do(g.createParams(&msg.Payload))
			message, _ := json.Marshal(map[string]interface{}{
				"id":      msg.OperationID,
				"type":    "data",
				"payload": res,
			})
			if err := msg.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				if err == websocket.ErrCloseSent {
					subscribers.Delete(key)
					return true
				}
				log.Printf("failed to write to ws connection: %v", err)
				return true
			}
			return true
		})
	}
}
