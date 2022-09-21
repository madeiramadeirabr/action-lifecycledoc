package types

type Scalar struct {
	format string
	enum   []interface{}
	value  interface{}

	generic
}

func (*Scalar) Type() TypeKeyword {
	return ScalarType
}

// Format returns the real definition type format. Can return empty
func (s *Scalar) Format() string {
	return s.format
}

// Enum returns all possible values. Can return empty
func (s *Scalar) Enum() []interface{} {
	return s.enum
}

func (s *Scalar) Value() interface{} {
	return s.value
}

func NewScalar(
	name, path, description string,
	nullable bool,
	format string,
	enum []interface{},
	value interface{},
) (*Scalar, error) {
	base, err := newGeneric(name, path, description, nullable)
	if err != nil {
		return nil, err
	}

	return &Scalar{
		generic: *base,
		format:  format,
		enum:    enum,
		value:   value,
	}, nil
}
