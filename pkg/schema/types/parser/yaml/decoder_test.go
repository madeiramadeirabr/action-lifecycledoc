package yaml_test

import (
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
    description: Um float
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

	typePathsToAssert := []string{
		"#/types/BoolType",
		"#/types/IntType",
		"#/types/FloatType",
		"#/types/StringType",
		"#/types/ArrayType",
		"#/types/ObjectType",
	}

	for _, path := range typePathsToAssert {
		if _, exists := schemaSpy.types[path]; !exists {
			t.Errorf("expected type added with '%s' path, no types have been added", path)
		}
	}
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
