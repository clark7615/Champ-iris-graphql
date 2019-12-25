package tests

import (
	"testing"

	"git.championtek.com.tw/go/champiris"
	"git.championtek.com.tw/go/champiris-contrib/graphql"
	"github.com/kataras/iris/v12/mvc"
)

var graph *graphql.Graphql

func init() {
	graph = graphql.Default()
	graph.ShowPlayground = true
}

func TestGraphql(t *testing.T) {
	var serivce champiris.Service
	_ = serivce.New(&champiris.NetConfig{
		Port: "8080",
	})
	serivce.App.Logger().SetLevel("debug")
	router := champiris.RouterSet{
		Party: "/service/v1",
		Router: func(m *mvc.Application) {
			m.Handle(graph)
		},
	}
	serivce.AddRoute(router)
	addSchema()
	err := serivce.Run()
	if err != nil {
		t.Error(err)
	}
}
