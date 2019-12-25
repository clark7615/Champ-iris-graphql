package graphql

import (
	"github.com/graphql-go/graphql"
)

type RootType struct {
	Obj    *graphql.Object
	fields graphql.Fields
}

func (s *RootType) new(name, description string) {
	s.fields = make(map[string]*graphql.Field, 0)
	config := graphql.ObjectConfig{
		Name:        name,
		Fields:      s.fields,
		Description: description,
	}
	s.Obj = graphql.NewObject(config)
}

func (s *RootType) AddField(field ...*graphql.Field) {
	for _, f := range field {
		s.fields[f.Name] = f
	}
}
