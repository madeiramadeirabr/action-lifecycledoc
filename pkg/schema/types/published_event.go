package types

import "errors"

type PublishedEvent struct {
	name        string
	visibility  EventVisibility
	module      string
	description string

	attributes TypeDescriber
	entities   TypeDescriber
}

func (p *PublishedEvent) Name() string {
	return p.name
}

func (p *PublishedEvent) Visibility() EventVisibility {
	return p.visibility
}

func (p *PublishedEvent) Module() string {
	return p.module
}

func (p *PublishedEvent) Description() string {
	return p.description
}

func (p *PublishedEvent) Attributes() TypeDescriber {
	return p.attributes
}

func (p *PublishedEvent) Entities() TypeDescriber {
	return p.entities
}

func NewPublishdEvent(
	name string,
	visibility EventVisibility,
	module string,
	description string,
	attributes TypeDescriber,
	entities TypeDescriber,
) (*PublishedEvent, error) {
	if len(name) < 1 {
		return nil, errors.New("the name cannot be empty")
	}

	if attributes == nil {
		return nil, errors.New("the attributes is required")
	}

	if entities == nil {
		return nil, errors.New("the entities is required")
	}

	return &PublishedEvent{
		name:        name,
		visibility:  visibility,
		module:      module,
		description: description,
		attributes:  attributes,
		entities:    entities,
	}, nil
}
