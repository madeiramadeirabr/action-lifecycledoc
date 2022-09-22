package types

import "errors"

type Array struct {
	items TypeDescriber

	generic
}

func (*Array) Type() TypeKeyword {
	return ArrayType
}

func (a *Array) Items() TypeDescriber {
	return a.items
}

func (a *Array) SetItems(items TypeDescriber) {
	a.items = items
}

func NewArray(
	name, path, description string,
	nullable bool,
	items TypeDescriber,
) (*Array, error) {
	base, err := newGeneric(name, path, description, nullable)
	if err != nil {
		return nil, err
	}

	if items == nil {
		return nil, errors.New("the items is required")
	}

	return &Array{
		generic: *base,
		items:   items,
	}, nil
}
