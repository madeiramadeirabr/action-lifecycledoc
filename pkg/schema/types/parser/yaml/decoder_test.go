package yaml_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types/parser/yaml"
)

func TestShouldParserValidYamlDefinition(t *testing.T) {
	input := strings.NewReader(`
version: "1.0"
name: super-cool-service

types:
  BoolType:
    description: Uma boolean
    type: boolean
    nullable: true
    value: null

  IntType:
    description: Um int
    type: integer
    value: 10

  FloatType:
    description: Um float
    type: number
    format: float com 14 casas decimais
    value: 22.12345678901234

  StringType:
    description: Uma string
    type: string
    enum:
      - option1
      - option2
    value: option1

  ArrayType:
    type: array
    description: Um array de strings!
    items:
      type: string
      value: Sou uma string

  ObjectType:
    type: object
    nullable: false # redundant, the default is false
    properties:
      id:
        type: integer
        value: 10
      name:
        type: string
        nullable: true
        description: O nome do object
        value: null`)

	schemaSpy := newSchameStorageSpy()

	if err := yaml.NewDecoder().Decode(input, schemaSpy); err != nil {
		t.Fatal(err)
	}

	testCases := map[string]testCase{
		"#/types/BoolType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Uma boolean",
			expectedNullable:    true,
			typeAssertions:      newScalarValueAssertion(t, nil),
		},
		"#/types/IntType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Um int",
			expectedNullable:    false,
			typeAssertions:      newScalarValueAssertion(t, 10),
		},
		"#/types/FloatType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Um float",
			expectedNullable:    false,
			typeAssertions:      newScalarValueAssertion(t, 22.12345678901234),
		},
		"#/types/StringType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Uma string",
			expectedNullable:    false,
			typeAssertions:      newScalarValueAssertion(t, "option10"),
		},
		"#/types/ArrayType": {
			expectedTypeKeyword: types.ArrayType,
			expectedDescription: "Um array de strings!",
			expectedNullable:    false,
			typeAssertions: func(t *testing.T, typeDef types.TypeDescriber) {
				arrayType, err := assertTypeCasting[*types.Array](t, typeDef)
				if err != nil {
					t.Error(err)
					return
				}

				assertScalarType(t, "Sou uma string", arrayType.Items())
			},
		},
		"#/types/ObjectType": {
			expectedTypeKeyword: types.ObjectType,
			expectedDescription: "",
			expectedNullable:    false,
			typeAssertions: func(t *testing.T, typeDef types.TypeDescriber) {
				objectType, err := assertTypeCasting[*types.Object](t, typeDef)
				if err != nil {
					t.Error(err)
					return
				}

				sliceProperties := objectType.Properties()
				properties := make(map[string]types.TypeDescriber)
				for i := range sliceProperties {
					properties[sliceProperties[i].Name()] = sliceProperties[i]
				}

				assertTypes(
					t,
					properties,
					map[string]testCase{
						"id": {
							expectedTypeKeyword: types.ScalarType,
							expectedDescription: "",
							expectedNullable:    false,
							typeAssertions:      newScalarValueAssertion(t, 10),
						},
						"name": {
							expectedTypeKeyword: types.ScalarType,
							expectedDescription: "O nome do object",
							expectedNullable:    true,
							typeAssertions:      newScalarValueAssertion(t, nil),
						},
					},
				)
			},
		},
	}

	assertTypes(t, schemaSpy.types, testCases)
}

type schemaStoragerSpy struct {
	types map[string]types.TypeDescriber
}

func (s *schemaStoragerSpy) AddType(t types.TypeDescriber) error {
	s.types[t.Path()] = t
	return nil
}

func newSchameStorageSpy() *schemaStoragerSpy {
	return &schemaStoragerSpy{
		types: make(map[string]types.TypeDescriber),
	}
}

func assertTypeCasting[T *types.Scalar | *types.Array | *types.Object](t *testing.T, typeDef types.TypeDescriber) (T, error) {
	t.Helper()

	result, is := typeDef.(T)
	if !is {
		return nil, fmt.Errorf("expected '%T' type, received '%T'", result, typeDef)
	}

	return result, nil
}

func newScalarValueAssertion(t *testing.T, expectedValue interface{}) func(*testing.T, types.TypeDescriber) {
	t.Helper()

	return func(t *testing.T, typeDef types.TypeDescriber) {
		t.Helper()
		assertScalarType(t, expectedValue, typeDef)
	}
}

func assertScalarType(
	t *testing.T,
	expectedValue interface{},
	typeDef types.TypeDescriber,
) {
	scalar, err := assertTypeCasting[*types.Scalar](t, typeDef)
	if err != nil {
		t.Error(err)
		return
	}

	if value := scalar.Value(); !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("expected '%v' value, received '%v'", expectedValue, value)
	}
}

type testCase struct {
	expectedTypeKeyword types.TypeKeyword
	expectedDescription string
	expectedNullable    bool
	typeAssertions      func(t *testing.T, typeDef types.TypeDescriber)
}

func assertTypes(t *testing.T, types map[string]types.TypeDescriber, testCases map[string]testCase) {
	for path, testCase := range testCases {
		t.Run(fmt.Sprintf("asserting %s", path), func(t *testing.T) {
			typeDef, exists := types[path]
			if !exists {
				t.Errorf("expected type added with '%s' path, no types have been added", path)
			}

			assertType(
				t,
				typeDef,
				testCase.expectedTypeKeyword,
				testCase.expectedDescription,
				testCase.expectedNullable,
				testCase.typeAssertions,
			)
		})
	}
}

func assertType(
	t *testing.T,
	typeDef types.TypeDescriber,
	expectedTypeKeyword types.TypeKeyword,
	expectedDescription string,
	expectedNullable bool,
	typeAssertions func(t *testing.T, typeDef types.TypeDescriber),
) {
	if typeKeyword := typeDef.Type(); typeKeyword != expectedTypeKeyword {
		t.Errorf("expected type '%s', received '%s'", expectedTypeKeyword, typeKeyword)
	}

	if description := typeDef.Description(); description != expectedDescription {
		t.Errorf("expected description '%s', received '%s'", expectedDescription, description)
	}

	if nullable := typeDef.Nullable(); nullable != expectedNullable {
		t.Errorf("expected nullable '%v', received '%v'", expectedNullable, nullable)
	}

	if typeAssertions != nil {
		t.Run("type assertions", func(t *testing.T) {
			typeAssertions(t, typeDef)
		})
	}
}
