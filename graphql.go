package graphql

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/graphql-go/graphql"
)

type CommandType int

const (
	Query            CommandType = iota
	Mutation         CommandType = iota
	Subscription     CommandType = iota
	QueryAndMutation CommandType = iota
	All              CommandType = iota
)

type Graphql struct {
	Query          RootType
	Mutation       RootType
	Subscription   RootType
	ShowPlayground bool
}

func Default() *Graphql {
	return New(QueryAndMutation)
}

func New(commandType CommandType) *Graphql {
	var ql = &Graphql{}
	switch commandType {
	case Query:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
	case Mutation:
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	case Subscription:
		ql.Subscription.new("Subscription", "訂閱相關的命令")
		go runSubscription(ql)
	case QueryAndMutation:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	case All:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
		ql.Subscription.new("Subscription", "訂閱相關的命令")
		go runSubscription(ql)
	}
	return ql
}

func runSubscription(g *Graphql) {
	for {
		time.Sleep(2 * time.Second)
		subscribers.Range(func(key, value interface{}) bool {
			msg, ok := value.(*ConnectionACKMessage)
			if !ok {
				return true
			}
			res := graphql.Do(graphql.Params{
				Schema:         g.newSchema(),
				RequestString:  msg.Payload.Query,
				VariableValues: msg.Payload.Variables,
				OperationName:  msg.Payload.OperationName,
			})
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
