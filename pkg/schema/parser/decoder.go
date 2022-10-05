package parser

import (
	"io"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

type SchemaStorager interface {
	SetProject(name string) error

	AddConfluencePage(title, spaceKey, ancestorID string) error

	AddType(t types.TypeDescriber) error
	AddPublishedEvent(e *types.PublishedEvent) error
	AddConsumedEvent(e *types.ConsumedEvent) error
}

type Decoder interface {
	Decode(r io.Reader, s SchemaStorager) error
}
