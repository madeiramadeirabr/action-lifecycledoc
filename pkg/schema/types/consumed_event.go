package types

import "errors"

type ConsumedEvent struct {
	name        string
	description string
}

func (c *ConsumedEvent) Name() string {
	return c.name
}

func (c *ConsumedEvent) Description() string {
	return c.description
}

func NewConsumedEvent(name, description string) (*ConsumedEvent, error) {
	if len(name) < 1 {
		return nil, errors.New("the name cannot be empty")
	}

	if len(description) < 1 {
		return nil, errors.New("the description cannot be empty")
	}

	return &ConsumedEvent{
		name:        name,
		description: description,
	}, nil
}
