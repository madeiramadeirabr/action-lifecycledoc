package schema

import "github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"

type Resolver interface {
	GetConfluence() (*types.Confluence, error)
	GetPublishedEvents() ([]*types.PublishedEvent, error)
	GetConsumedEvents() ([]*types.ConsumedEvent, error)
	GetTypes() ([]types.TypeDescriber, error)
}
