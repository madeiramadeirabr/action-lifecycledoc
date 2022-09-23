package types

type ScalarReference struct {
	reference *Reference

	*Scalar
}

func (s *ScalarReference) Reference() string {
	return s.Scalar.Name()
}

func (s *ScalarReference) Name() string {
	return s.reference.Name()
}

func (s *ScalarReference) Path() string {
	return s.reference.Path()
}

func (s *ScalarReference) Description() string {
	if len(s.reference.Description()) > 0 {
		return s.reference.Description()
	}

	return s.Scalar.Description()
}

func (s *ScalarReference) Nullable() bool {
	return s.reference.Nullable()
}

func NewScalarReference(reference *Reference, scalar *Scalar) *ScalarReference {
	return &ScalarReference{
		reference: reference,
		Scalar:    scalar,
	}
}
