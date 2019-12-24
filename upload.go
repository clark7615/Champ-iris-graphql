package graphql

import (
	"github.com/graphql-go/graphql"
	"mime/multipart"
)

type MultipartFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}

var Upload = graphql.NewScalar(graphql.ScalarConfig{
	Name: "Upload",
	ParseValue: func(value interface{}) interface{} {
		if v, ok := value.(*MultipartFile); ok {
			return v
		}
		return nil
	},
})
