package graphql

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
	case QueryAndMutation:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
	case All:
		ql.Query.new("Query", "收尋&取得資料的相關命令")
		ql.Mutation.new("Mutation", "主要用在建立、修改、刪除的相關命令")
		ql.Subscription.new("Subscription", "訂閱相關的命令")
	}
	return ql
}
