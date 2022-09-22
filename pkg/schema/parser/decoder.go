package parser

import (
	"io"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

type SchemaStorager interface {
	AddType(types.TypeDescriber) error
	AddPublishedEvent(*types.PublishedEvent) error
}

type Decoder interface {
	Decode(io.Reader, SchemaStorager) error
}
