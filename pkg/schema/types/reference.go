package types

import "errors"

type Reference struct {
	reference string

	generic
}

func (*Reference) Type() string {
	return ReferenceType
}

func (r *Reference) Reference() string {
	return r.reference
}

func NewReference(
	name, path, description string,
	nullable bool,
	reference string,
) (*Reference, error) {
	base, err := newGeneric(name, path, description, nullable)
	if err != nil {
		return nil, err
	}

	if len(reference) < 1 {
		return nil, errors.New("the reference cannot be empty")
	}

	return &Reference{
		generic:   *base,
		reference: reference,
	}, nil
}
