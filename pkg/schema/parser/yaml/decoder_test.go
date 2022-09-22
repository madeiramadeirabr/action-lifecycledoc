package yaml_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/parser/yaml"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

func TestShouldParserValidYamlDefinition(t *testing.T) {
	input := strings.NewReader(`
version: "1.0"
name: super-cool-service

confluence:
  pages:
    - spaceKey: "SPACEKEY"
      ancestorId: "123456789"
      title: Titulo

events:
  published:
    SOME_COOL_EVENT:
      visibility: public
      module: test
      description: Evento de teste

      attributes:
        type: object
        properties:
          id:
            type: integer
            value: 9932
            description: ID do evento
          description:
            type: string
            nullable: true
            value: null
      
      entities:
        type: object
        properties:
          id:
            $ref: '#/types/IntType'

  consumed:
    EXTERNAL_EVENT:
      description: Usado para demonstração

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
    format: /\d+.\d{14}/
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
			typeAssertions:      newScalarAssertion(t, "", nil, nil),
		},
		"#/types/IntType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Um int",
			expectedNullable:    false,
			typeAssertions:      newScalarAssertion(t, "", nil, 10),
		},
		"#/types/FloatType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Um float",
			expectedNullable:    false,
			typeAssertions:      newScalarAssertion(t, `/\d+.\d{14}/`, nil, 22.12345678901234),
		},
		"#/types/StringType": {
			expectedTypeKeyword: types.ScalarType,
			expectedDescription: "Uma string",
			expectedNullable:    false,
			typeAssertions: newScalarAssertion(
				t,
				"",
				[]interface{}{
					"option1",
					"option2",
				},
				"option1",
			),
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

				assertScalarType(t, "", nil, "Sou uma string", arrayType.Items())
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
							typeAssertions:      newScalarAssertion(t, "", nil, 10),
						},
						"name": {
							expectedTypeKeyword: types.ScalarType,
							expectedDescription: "O nome do object",
							expectedNullable:    true,
							typeAssertions:      newScalarAssertion(t, "", nil, nil),
						},
					},
				)
			},
		},
	}

	assertTypes(t, schemaSpy.types, testCases)

	t.Run("should register published event", func(t *testing.T) {
		event, exists := schemaSpy.publishedEvents["SOME_COOL_EVENT"]
		if !exists {
			t.Fatal("expected event SOME_COOL_EVENT in published events")
			return
		}

		if visibility := event.Visibility(); visibility != types.EventPublic {
			t.Errorf("expected 'public' visibility, received '%s'", visibility)
		}

		attributes, err := assertTypeCasting[*types.Object](t, event.Attributes())
		if err != nil {
			t.Fatal(err)
		}

		if len(attributes.Properties()) != 2 {
			t.Errorf("expected '2' properties, received '%d'", len(attributes.Properties()))
		}

		entities, err := assertTypeCasting[*types.Object](t, event.Entities())
		if err != nil {
			t.Fatal(err)
		}

		refDefinition, err := assertTypeCasting[*types.Reference](t, entities.Properties()[0])
		if err != nil {
			t.Fatal(err)
		}

		if refence := refDefinition.Reference(); refence != "#/types/IntType" {
			t.Errorf("expected '#/types/IntType' reference, received '%s'", refence)
		}
	})

	t.Run("should register consumed event", func(t *testing.T) {
		event, exists := schemaSpy.consumedEvents["EXTERNAL_EVENT"]
		if !exists {
			t.Fatal("expected event EXTERNAL_EVENT in consumed events")
		}

		expected := "Usado para demonstração"
		if event.Description() != expected {
			t.Errorf("expected '%s' description, received '%s'", expected, event.Description())
		}
	})

	t.Run("should add name and confluence pages", func(t *testing.T) {
		expectedName := "super-cool-service"
		if schemaSpy.name != expectedName {
			t.Errorf("expected '%s' name, received '%s'", expectedName, schemaSpy.name)
		}

		if length := len(schemaSpy.confluencePages); length != 1 {
			t.Fatalf("expected '1' confluence page, received '%d'", length)
		}

		var (
			expectedSpaceKey   = "SPACEKEY"
			expectedAncestorId = "123456789"
			expectedTitle      = "Titulo"
		)

		if spaceKey := schemaSpy.confluencePages[0].SpaceKey(); spaceKey != expectedSpaceKey {
			t.Errorf("expected '%s' space key, received '%s'", expectedSpaceKey, spaceKey)
		}

		if ancestorID := schemaSpy.confluencePages[0].AncestorID(); ancestorID != expectedAncestorId {
			t.Errorf("expected '%s' ancestor id, received '%s'", expectedAncestorId, ancestorID)
		}

		if title := schemaSpy.confluencePages[0].Title(); title != expectedTitle {
			t.Errorf("expected '%s' title, received '%s'", expectedTitle, title)
		}
	})
}

type schemaStoragerSpy struct {
	name            string
	confluencePages []*types.ConfluencePage

	types           map[string]types.TypeDescriber
	publishedEvents map[string]*types.PublishedEvent
	consumedEvents  map[string]*types.ConsumedEvent
}

func (s *schemaStoragerSpy) SetProject(name string) error {
	s.name = name
	return nil
}

func (s *schemaStoragerSpy) AddConfluencePage(title, spaceKey, ancestorID string) error {
	page, err := types.NewConfluencePage(title, spaceKey, ancestorID)
	if err != nil {
		return err
	}

	s.confluencePages = append(s.confluencePages, page)
	return nil
}

func (s *schemaStoragerSpy) AddType(t types.TypeDescriber) error {
	s.types[t.Path()] = t
	return nil
}

func (s *schemaStoragerSpy) AddPublishedEvent(e *types.PublishedEvent) error {
	s.publishedEvents[e.Name()] = e
	return nil
}

func (s *schemaStoragerSpy) AddConsumedEvent(e *types.ConsumedEvent) error {
	s.consumedEvents[e.Name()] = e
	return nil
}

func newSchameStorageSpy() *schemaStoragerSpy {
	return &schemaStoragerSpy{
		types:           make(map[string]types.TypeDescriber),
		publishedEvents: make(map[string]*types.PublishedEvent),
		consumedEvents:  make(map[string]*types.ConsumedEvent),
	}
}

func assertTypeCasting[T *types.Scalar | *types.Array | *types.Object | *types.Reference](t *testing.T, typeDef types.TypeDescriber) (T, error) {
	t.Helper()

	result, is := typeDef.(T)
	if !is {
		return nil, fmt.Errorf("expected '%T' type, received '%T'", result, typeDef)
	}

	return result, nil
}

func newScalarAssertion(
	t *testing.T,
	expectedFormat string,
	expectedEnum []interface{},
	expectedValue interface{},
) func(*testing.T, types.TypeDescriber) {
	t.Helper()

	return func(t *testing.T, typeDef types.TypeDescriber) {
		t.Helper()
		assertScalarType(t, expectedFormat, expectedEnum, expectedValue, typeDef)
	}
}

func assertScalarType(
	t *testing.T,
	expectedFormat string,
	expectedEnum []interface{},
	expectedValue interface{},
	typeDef types.TypeDescriber,
) {
	scalar, err := assertTypeCasting[*types.Scalar](t, typeDef)
	if err != nil {
		t.Error(err)
		return
	}

	if format := scalar.Format(); format != expectedFormat {
		t.Errorf("expected '%s' format, received '%s'", expectedFormat, format)
	}

	if enum := scalar.Enum(); !reflect.DeepEqual(enum, expectedEnum) {
		t.Errorf("expected '%v' enum, received '%v'", expectedEnum, enum)
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
				return
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
