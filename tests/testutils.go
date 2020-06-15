package tests

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

func addSchema() {
	Graph.Subscription.AddField(&graphql.Field{
		Name: "qq",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type: graphql.Int,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return p.Args["id"], nil
		},
	})
	Graph.Query.AddField(&graphql.Field{
		Name: "qq",
		Type: graphql.Int,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return 1, errors.New("WTF")
		},
	})
	Graph.Mutation.AddField(&graphql.Field{
		Name: "ff",
		Type: graphql.Int,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"bool": &graphql.ArgumentConfig{
				Type:         graphql.Boolean,
				DefaultValue: true,
			},
			"aa": &graphql.ArgumentConfig{
				Type: graphql.Float,
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			type Member struct {
				ID   int     `json:"id"`
				Bool bool    `json:"bool"`
				AA   float64 `json:"aa"`
			}
			member := Member{}
			Graph.ToStruct(p.Args, &member)
			fmt.Println(member)
			return member.ID, nil
		},
	})
}
