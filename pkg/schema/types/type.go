package types

const (
	ScalarType    TypeKeyword = "scalar"
	ReferenceType TypeKeyword = "reference"
	ObjectType    TypeKeyword = "object"
	ArrayType     TypeKeyword = "array"
)

type TypeKeyword string

type TypeDescriber interface {
	// Name of definition
	Name() string

	// Path of definition in schema
	Path() string

	// Type of definition
	Type() TypeKeyword

	// Description of definition
	Description() string

	// Nullable indicates if type accepts null values
	Nullable() bool
}
