package graphql

import (
	"encoding/json"
	"log"

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

//Default 建立一個 Graphql 實體並包含 Query以及Mutation RootType
func Default() *Graphql {
	return New(QueryAndMutation)
}

//New 依據CommandType建立一個Graphql 實體並新增對應之RootType
func New(commandType commandType) *Graphql {
	ql := new(Graphql)
	switch commandType {
	case Query:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
	case Mutation:
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	case Subscription:
		ql.Subscription.new("Subscription", "訂閱相關的命令")
	case All:
		ql.Subscription.new("Subscription", "訂閱相關的命令")
		fallthrough
	case QueryAndMutation:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	}
	return ql
}

//CheckSubscription subscriptionName:監聽的Schema名稱 check:需要驗證的Graphql Variables
func (g *Graphql) CheckSubscription(subscriptionName string, check map[string]interface{}) bool {
	var b bool
	subscribers.Range(func(key, value interface{}) bool {
		msg, ok := value.(*connectionACKMessage)
		if !ok {
			return true
		}
		if msg.Payload.OperationName == subscriptionName {
			opt := g.createParams(&msg.Payload)
			for s, i := range opt.VariableValues {
				if check[s] != nil {
					if check[s] != i {
						return true
					}
				}
			}
			res := graphql.Do(opt)

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
			b = true
			return true
		}
		return true
	})
	return b
}
