package types

type ArrayReference struct {
	reference *Reference

	*Array
}

func (a *ArrayReference) Reference() string {
	return a.Array.Name()
}

func (s *ArrayReference) Name() string {
	return s.reference.Name()
}

func (s *ArrayReference) Path() string {
	return s.reference.Path()
}

func (s *ArrayReference) Description() string {
	if len(s.reference.Description()) > 0 {
		return s.reference.Description()
	}

	return s.Array.Description()
}

func (s *ArrayReference) Nullable() bool {
	return s.reference.Nullable()
}

func NewArrayReference(reference *Reference, array *Array) *ArrayReference {
	return &ArrayReference{
		reference: reference,
		Array:     array,
	}
}
