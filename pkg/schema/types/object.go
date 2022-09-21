package types

import "errors"

type Object struct {
	// properties is slice instead of map to keep order of declaration
	properties []TypeDescriber

	generic
}

func (*Object) Type() TypeKeyword {
	return ObjectType
}

// Properties of object in declaration order
func (o *Object) Properties() []TypeDescriber {
	return o.properties
}

func NewObject(
	name, path, description string,
	nullable bool,
	properties []TypeDescriber,
) (*Object, error) {
	base, err := newGeneric(name, path, description, nullable)
	if err != nil {
		return nil, err
	}

	if len(properties) < 1 {
		return nil, errors.New("the properties is required")
	}

	return &Object{
		generic:    *base,
		properties: properties,
	}, nil
}
