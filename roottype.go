package graphql

import (
	"github.com/graphql-go/graphql"
)

type rootType struct {
	Obj    *graphql.Object
	fields graphql.Fields
}

//AddField 新增Field至RootType
func (s *rootType) AddField(field ...*graphql.Field) {
	for _, f := range field {
		s.fields[f.Name] = f
	}
}

func (s *rootType) new(name, description string) {
	s.fields = make(map[string]*graphql.Field, 0)
	config := graphql.ObjectConfig{
		Name:        name,
		Fields:      s.fields,
		Description: description,
	}
	s.Obj = graphql.NewObject(config)
}
