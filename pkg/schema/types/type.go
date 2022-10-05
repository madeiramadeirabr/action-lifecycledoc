package types

const (
	ReferenceType string = "reference"
	ObjectType    string = "object"
	ArrayType     string = "array"
)

type TypeDescriber interface {
	// Name of definition
	Name() string

	// Path of definition in schema
	Path() string

	// Type of definition
	Type() string

	// Description of definition
	Description() string

	// Nullable indicates if type accepts null values
	Nullable() bool
}
