package graphql

import (
	"html/template"

	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

type Graphql struct {
	Ctx            iris.Context
	Query          RootType
	Mutation       RootType
	Subscription   RootType
	ShowPlayground bool
}

func (g *Graphql) PostQuery() {
	params := g.createParams()
	res := graphql.Do(params)
	_, _ = g.Ctx.JSON(res)
}

func (g *Graphql) GetPg() {
	if !g.ShowPlayground {
		return
	}
	t := template.New("Playground")
	te, _ := t.Parse(html)
	_ = te.ExecuteTemplate(g.Ctx.ResponseWriter(), "index", nil)
}

func (g *Graphql) createParams() graphql.Params {
	opt := g.getRequestOptions()
	return graphql.Params{
		Schema:         g.newSchema(),
		RequestString:  opt.Query,
		VariableValues: opt.Variables,
		OperationName:  opt.OperationName,
	}
}
