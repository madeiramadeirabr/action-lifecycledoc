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

	if project.Version != "1.0" {
		return fmt.Errorf("unsupported '%s' version", project.Version)
	}

	if err := schema.SetProject(project.Name); err != nil {
		return err
	}

	if err := d.parseConfluence(project, schema); err != nil {
		return err
	}

	if err := d.parseTypes(project, schema); err != nil {
		return err
	}

	if err := d.parsePublishedEvents(project, schema); err != nil {
		return err
	}

	if err := d.parseConsumedEvents(project, schema); err != nil {
		return err
	}

	return nil
}

func (d *decoder) parseConfluence(project *project, schema parser.SchemaStorager) error {
	for i := range project.Confluence.Pages {
		err := schema.AddConfluencePage(
			project.Confluence.Pages[i].Title,
			project.Confluence.Pages[i].SpaceKey,
			project.Confluence.Pages[i].AncestorID,
		)

		if err != nil {
			return addPathToError("#/confluence/pages", fmt.Errorf("can't add page at '%d' index: %w", i, err))
		}
	}

	return nil
}

func (d *decoder) parseTypes(project *project, schema parser.SchemaStorager) error {
	types, err := d.parseTypeDefinitions("#/types", project.Types)
	if err != nil {
		return fmt.Errorf("can't parse types: %w", err)
	}

	for i := range types {
		if err := schema.AddType(types[i]); err != nil {
			return fmt.Errorf("can't register type: %w", err)
		}
	}

	return nil
}

func (d *decoder) parsePublishedEvents(project *project, schema parser.SchemaStorager) error {
	for i := range project.Events.Published {
		name, path := d.yamlMapItemToNameAndPath("#/events/published", project.Events.Published[i])

		rawEventDefinition, err := d.yamlMapItemValueToMap(path, project.Events.Published[i].Value)
		if err != nil {
			return err
		}

		visibility, err := d.parseEventVisibility(path, rawEventDefinition)
		if err != nil {
			return err
		}

		module, _ := rawEventDefinition["module"].(string)
		description, _ := rawEventDefinition["description"].(string)

		attributesType, err := d.parserEventTypeDefinition(path, "attributes", rawEventDefinition)
		if err != nil {
			return err
		}

		entitiesType, err := d.parserEventTypeDefinition(path, "entities", rawEventDefinition)
		if err != nil {
			return err
		}

		event, err := types.NewPublishdEvent(
			name,
			visibility,
			module,
			description,
			attributesType,
			entitiesType,
		)
		if err != nil {
			return addPathToError(path, err)
		}

		if err := schema.AddPublishedEvent(event); err != nil {
			return fmt.Errorf("can't register published event: %w", err)
		}
	}

	return nil
}

func (d *decoder) parseConsumedEvents(project *project, schema parser.SchemaStorager) error {
	for i := range project.Events.Consumed {
		name, path := d.yamlMapItemToNameAndPath("#/events/consumed", project.Events.Consumed[i])

		rawEventDefinition, err := d.yamlMapItemValueToMap(path, project.Events.Consumed[i].Value)
		if err != nil {
			return err
		}

		description, _ := rawEventDefinition["description"].(string)

		event, err := types.NewConsumedEvent(name, description)
		if err != nil {
			return addPathToError(path, err)
		}

		if err := schema.AddConsumedEvent(event); err != nil {
			return fmt.Errorf("can't register consumed event: %w", err)
		}
	}

	return nil
}

func (d *decoder) parseEventVisibility(path string, eventDefinition map[string]interface{}) (types.EventVisibility, error) {
	visibility, _ := eventDefinition["visibility"].(string)

	v, err := types.NewEventVisibility(visibility)
	if err != nil {
		err = addPathToError(path, err)
	}

	return v, err
}

func (d *decoder) parserEventTypeDefinition(path, key string, rawEventDefinition map[string]interface{}) (types.TypeDescriber, error) {
	rawType, err := d.extractYamlMapSliceFromMap(path, key, rawEventDefinition)
	if err != nil {
		return nil, err
	}

	typeDefinition, err := d.parseTypeDefinition(
		key,
		fmt.Sprintf("%s/%s", path, key),
		d.yamlMapSliceToMap(rawType),
	)
	if err != nil {
		return nil, err
	}

	return typeDefinition, nil
}

func (d *decoder) parseTypeDefinitions(path string, rawTypes yaml.MapSlice) ([]types.TypeDescriber, error) {
	var typeDefinitions []types.TypeDescriber

	for i := range rawTypes {
		name, path := d.yamlMapItemToNameAndPath(path, rawTypes[i])

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

func (d *decoder) yamlMapItemToNameAndPath(path string, item yaml.MapItem) (string, string) {
	name := fmt.Sprint(item.Key)
	return name, fmt.Sprintf("%s/%s", path, name)
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
		return nil, fmt.Errorf("%s/%s: unexpected structure", path, key)
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
			err = addPathToError(path, err)
		}

		return referenceType, err
	}

	typeKeyword, isTypeKeywordValid := typeDefinition["type"].(string)
	if !isTypeKeywordValid {
		return nil, fmt.Errorf("%s/type: invalid type identifier declaration", path)
	}

	switch typeKeyword {
	case "integer":
		return parseScalarType[int](name, path, description, typeKeyword, nullable, typeDefinition)
	case "number":
		return parseScalarType[float64](name, path, description, typeKeyword, nullable, typeDefinition)
	case "string":
		return parseScalarType[string](name, path, description, typeKeyword, nullable, typeDefinition)
	case "boolean":
		return parseScalarType[bool](name, path, description, typeKeyword, nullable, typeDefinition)
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
			return nil, addPathToError(path, err)
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
			return nil, addPathToError(path, err)
		}

		return objectType, nil
	default:
		return nil, fmt.Errorf("%s/type: '%s' not supported", path, typeKeyword)
	}
}

func addPathToError(path string, err error) error {
	return fmt.Errorf("%s: %w", path, err)
}

func parseScalarType[T scalar](
	name, path, description, typeKeyword string,
	nullable bool,
	typeDefinition map[string]interface{},
) (types.TypeDescriber, error) {
	value, err := parserScalarValue[T](path, nullable, typeDefinition)
	if err != nil {
		return nil, err
	}

	var rawValue interface{}
	if value != nil {
		// Use value insted of pointer
		rawValue = *value
	} else {
		rawValue = nil
	}

	enumValues, err := parserScalarEnum[T](path, nullable, typeDefinition)
	if err != nil {
		return nil, err
	}

	format, _ := typeDefinition["format"].(string)

	scalarType, err := types.NewScalar(
		name,
		path,
		description,
		nullable,
		typeKeyword,
		format,
		enumValues,
		rawValue,
	)

	if err != nil {
		return nil, addPathToError(path, err)
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

func parserScalarEnum[T scalar](path string, nullable bool, typeDefinition map[string]interface{}) ([]interface{}, error) {
	enumValues, _ := typeDefinition["enum"].([]interface{})

	for i := range enumValues {
		if nullable && enumValues[i] == nil {
			continue
		}

		_, is := enumValues[i].(T)
		if !is {
			return nil, fmt.Errorf("%s/enum: invalid enum type at %d position", path, i)
		}
	}

	return enumValues, nil
}

type scalar interface {
	int | float64 | string | bool
}
