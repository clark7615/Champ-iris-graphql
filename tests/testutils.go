package tests

import (
	"errors"

	"github.com/graphql-go/graphql"
)

func addSchema() {

	graph.Query.AddField(&graphql.Field{
		Name: "qq",
		Type: graphql.Int,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return 1, errors.New("WTF")
		},
	})
	graph.Mutation.AddField(&graphql.Field{
		Name: "ff",
		Type: graphql.Int,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			type Member struct {
				ID int `json:"id"`
			}
			member := Member{}
			return member.ID, nil
		},
	})
}
