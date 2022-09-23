package types

type ObjectReference struct {
	reference *Reference

	*Object
}

func (o *ObjectReference) Reference() string {
	return o.Object.Name()
}

func (s *ObjectReference) Name() string {
	return s.reference.Name()
}

func (s *ObjectReference) Path() string {
	return s.reference.Path()
}

func (s *ObjectReference) Description() string {
	if len(s.reference.Description()) > 0 {
		return s.reference.Description()
	}

	return s.Object.Description()
}

func (s *ObjectReference) Nullable() bool {
	return s.reference.Nullable()
}

func NewObjectReference(reference *Reference, object *Object) *ObjectReference {
	return &ObjectReference{
		reference: reference,
		Object:    object,
	}
}
