package types

import (
	"errors"
)

type generic struct {
	name        string
	path        string
	description string
	nullable    bool
}

func (g *generic) Name() string {
	return g.name
}

func (g *generic) Path() string {
	return g.path
}

func (g *generic) Description() string {
	return g.description
}

func (g *generic) Nullable() bool {
	return g.nullable
}

func newGeneric(name, path, description string, nullable bool) (*generic, error) {
	if len(name) < 1 {
		return nil, errors.New("the name cannot be empty")
	}

	if len(path) < 1 {
		return nil, errors.New("the path cannot be empty")
	}

	return &generic{
		name:        name,
		path:        path,
		description: description,
		nullable:    nullable,
	}, nil
}
