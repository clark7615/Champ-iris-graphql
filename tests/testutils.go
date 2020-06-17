package tests

import (
	"errors"

	QL "git.championtek.com.tw/go/champiris-contrib/graphql"
	"github.com/graphql-go/graphql"
)

func addSchema() {
	type Member struct {
		Account string  `json:"account"`
		Bool    bool    `json:"bool"`
		AA      float64 `json:"aa"`
	}
	var member = graphql.NewObject(graphql.ObjectConfig{
		Name:   "Member",
		Fields: graphql.BindFields(Member{}),
	})
	Graph.Subscription.AddField(&graphql.Field{
		Name: "Member",
		Args: graphql.FieldConfigArgument{
			"account": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type: graphql.Boolean,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return true, nil
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
		Name: "Member",
		Type: member,
		Args: graphql.FieldConfigArgument{
			"account": &graphql.ArgumentConfig{
				Type: graphql.String,
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
			member := Member{}
			QL.ToStruct(p.Args, &member)
			Graph.CheckSubscription("Member", map[string]interface{}{
				"account": member.Account,
			})
			return member, nil
		},
	})
}
