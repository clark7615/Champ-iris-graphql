package tests

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

func addSchema() {

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
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			type Member struct {
				ID   int  `json:"id"`
				Bool bool `json:"bool"`
			}
			member := Member{}
			Graph.ToStruct(p.Args, &member)
			fmt.Println(member)
			return member.ID, nil
		},
	})
}
