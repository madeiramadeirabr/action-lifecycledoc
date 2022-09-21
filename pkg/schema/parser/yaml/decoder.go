package yaml

import (
	"fmt"
	"io"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/parser"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
	"gopkg.in/yaml.v2"
)

type decoder struct{}

func NewDecoder() parser.Decoder {
	return &decoder{}
}

func (d *decoder) Decode(definition io.Reader, schema parser.SchemaStorager) error {
	var project *project
	if err := yaml.NewDecoder(definition).Decode(&project); err != nil {
		return fmt.Errorf("can't decode yaml definition: %w", err)
	}

	types, err := d.parseTypeDefinitions("#/types", project.Types)
	if err != nil {
		return fmt.Errorf("can't parse types: %w", err)
	}

	for i := range types {
		if err := schema.AddType(types[i]); err != nil {
			return fmt.Errorf("can't register type: %w", err)
		}
	}

	// @todo: parser consumed and published events

	return nil
}

func (d *decoder) parseTypeDefinitions(path string, rawTypes yaml.MapSlice) ([]types.TypeDescriber, error) {
	var typeDefinitions []types.TypeDescriber

	for i := range rawTypes {
		name := fmt.Sprint(rawTypes[i].Key)
		path := fmt.Sprintf("%s/%s", path, name)

		rawTypeDefinition, err := d.yamlMapItemValueToMap(path, rawTypes[i].Value)
		if err != nil {
			return nil, err
		}

		typeDefinition, err := d.parseTypeDefinition(name, path, rawTypeDefinition)
		if err != nil {
			return nil, err
		}

		typeDefinitions = append(typeDefinitions, typeDefinition)
	}

	return typeDefinitions, nil
}

func (d *decoder) yamlMapItemValueToMap(path string, value interface{}) (map[string]interface{}, error) {
	mapSlice, is := value.(yaml.MapSlice)
	if !is {
		return nil, fmt.Errorf("%s: unexpected structure", path)
	}

	return d.yamlMapSliceToMap(mapSlice), nil
}

func (d *decoder) extractYamlMapSliceFromMap(path, key string, m map[string]interface{}) (yaml.MapSlice, error) {
	mapSlice, is := m[key].(yaml.MapSlice)
	if !is {
		return nil, fmt.Errorf("%s: unexpected structure", path)
	}

	return mapSlice, nil
}

func (*decoder) yamlMapSliceToMap(mapSlice yaml.MapSlice) map[string]interface{} {
	result := make(map[string]interface{})
	for i := range mapSlice {
		result[fmt.Sprint(mapSlice[i].Key)] = mapSlice[i].Value
	}

	return result
}

func (d *decoder) parseTypeDefinition(
	name, path string,
	typeDefinition map[string]interface{},
) (types.TypeDescriber, error) {
	description, _ := typeDefinition["description"].(string)
	nullable, _ := typeDefinition["nullable"].(bool)

	reference, isReference := typeDefinition["$ref"].(string)
	if isReference {
		referenceType, err := types.NewReference(name, path, description, nullable, reference)
		if err != nil {
			err = addPathToErr(path, err)
		}

		return referenceType, err
	}

	typeKeyword, isTypeKeywordValid := typeDefinition["type"].(string)
	if !isTypeKeywordValid {
		return nil, fmt.Errorf("%s: invalid type declaration", path)
	}

	switch typeKeyword {
	case "integer":
		return parseScalarType[int](name, path, description, nullable, typeDefinition)
	case "number":
		return parseScalarType[float64](name, path, description, nullable, typeDefinition)
	case "string":
		return parseScalarType[string](name, path, description, nullable, typeDefinition)
	case "boolean":
		return parseScalarType[bool](name, path, description, nullable, typeDefinition)
	case "array":
		rawItems, err := d.extractYamlMapSliceFromMap(path, "items", typeDefinition)
		if err != nil {
			return nil, err
		}

		itemsType, err := d.parseTypeDefinition(
			"items",
			fmt.Sprintf("%s/items", path),
			d.yamlMapSliceToMap(rawItems),
		)
		if err != nil {
			return nil, err
		}

		arrayType, err := types.NewArray(
			name,
			path,
			description,
			nullable,
			itemsType,
		)
		if err != nil {
			return nil, addPathToErr(path, err)
		}

		return arrayType, nil
	case "object":
		rawProperties, err := d.extractYamlMapSliceFromMap(path, "properties", typeDefinition)
		if err != nil {
			return nil, err
		}

		typeDefinitions, err := d.parseTypeDefinitions(fmt.Sprintf("%s/properties", path), rawProperties)
		if err != nil {
			return nil, err
		}

		objectType, err := types.NewObject(
			name,
			path,
			description,
			nullable,
			typeDefinitions,
		)
		if err != nil {
			return nil, addPathToErr(path, err)
		}

		return objectType, nil
	default:
		return nil, fmt.Errorf("%s/type: '%s' not supported", path, typeKeyword)
	}
}

func addPathToErr(path string, err error) error {
	return fmt.Errorf("%s: %w", path, err)
}

func parseScalarType[T scalar](
	name, path, description string,
	nullable bool,
	typeDefinition map[string]interface{},
) (types.TypeDescriber, error) {
	value, err := parserScalarValue[T](path, nullable, typeDefinition)
	if err != nil {
		return nil, err
	}

	var rawValue interface{}
	if value != nil {
		// Use value if insted of pointer
		rawValue = *value
	} else {
		rawValue = nil
	}

	scalarType, err := types.NewScalar(
		name,
		path,
		description,
		nullable,
		// @todo
		"",
		nil,
		rawValue,
	)

	if err != nil {
		return nil, addPathToErr(path, err)
	}

	return scalarType, nil
}

func parserScalarValue[T scalar](path string, nullable bool, typeDefinition map[string]interface{}) (*T, error) {
	if nullable && typeDefinition["value"] == nil {
		return nil, nil
	}

	value, is := typeDefinition["value"].(T)
	if !is {
		return nil, fmt.Errorf("%s/value: is not of type '%T'", path, value)
	}

	return &value, nil
}

type scalar interface {
	int | float64 | string | bool
}
