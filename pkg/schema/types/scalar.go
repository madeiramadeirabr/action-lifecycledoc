package types

const (
	ScalarIntegerType = "integer"
	ScalarNumberType  = "number"
	ScalarStringType  = "string"
	ScalarBooleanType = "boolean"
)

type ScalarDescriber interface {
	Format() string
	HasFormat() bool
	Enum() []interface{}
	HasEnum() bool
	Value() interface{}

	TypeDescriber
}

type Scalar struct {
	format      string
	typeKeyword string
	enum        []interface{}
	value       interface{}

	generic
}

func (s *Scalar) Type() string {
	return s.typeKeyword
}

// Format returns the real definition type format. Can return empty
func (s *Scalar) Format() string {
	return s.format
}

func (s *Scalar) HasFormat() bool {
	return len(s.Format()) > 0
}

// Enum returns all possible values. Can return empty
func (s *Scalar) Enum() []interface{} {
	return s.enum
}

func (s *Scalar) HasEnum() bool {
	return len(s.Enum()) > 0
}

func (s *Scalar) Value() interface{} {
	return s.value
}

func NewScalar(
	name, path, description string,
	nullable bool,
	typeKeyword, format string,
	enum []interface{},
	value interface{},
) (*Scalar, error) {
	base, err := newGeneric(name, path, description, nullable)
	if err != nil {
		return nil, err
	}

	return &Scalar{
		generic:     *base,
		typeKeyword: typeKeyword,
		format:      format,
		enum:        enum,
		value:       value,
	}, nil
}
