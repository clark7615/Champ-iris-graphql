package graphql

import (
	"html/template"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/kataras/iris/v12"
)

func (g *Graphql) PostQuery(Ctx iris.Context) {
	params := g.createParams(Ctx)
	res := graphql.Do(params)
	_, _ = Ctx.JSON(res)
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
