package confluence

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/jsonc"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

type TemplateRetriver interface {
	Retrive() string
}

type TemplateRetriverFunc func() string

func (t TemplateRetriverFunc) Retrive() string {
	return t()
}

type TemplateWriter struct {
	templateRetriver TemplateRetriver
	exampleWriter    *strings.Builder
	exampleEncoder   jsonc.Encoder
}

func NewTemplateWriter(templateRetriver TemplateRetriver) *TemplateWriter {
	return &TemplateWriter{
		templateRetriver: templateRetriver,
	}
}

func (t *TemplateWriter) Write(w io.Writer, schemaResolver schema.Resolver) error {
	t.exampleWriter = &strings.Builder{}
	t.exampleEncoder = jsonc.NewEncoder(t.exampleWriter)

	out := &outputData{}

	if err := t.prepareTypes(out, schemaResolver); err != nil {
		return fmt.Errorf("can't prepare types to write: %w", err)
	}

	if err := t.preparePublishedEvents(out, schemaResolver); err != nil {
		return fmt.Errorf("can't prepare published events to write: %w", err)
	}

	if err := t.prepareConsumedEvents(out, schemaResolver); err != nil {
		return fmt.Errorf("can't prepare consumed events to write: %w", err)
	}

	template, err := template.New("events_page").Parse(t.templateRetriver.Retrive())
	if err != nil {
		return fmt.Errorf("can't parse output template: %w", err)
	}

	if err := template.Execute(w, out); err != nil {
		return fmt.Errorf("can't write output template: %w", err)
	}

	return nil
}

func (t *TemplateWriter) prepareTypes(out *outputData, schemaResolver schema.Resolver) error {
	typesDefinitions, err := schemaResolver.GetTypes()
	if err != nil {
		return fmt.Errorf("can't get types to write: %w", err)
	}

	for i := range typesDefinitions {
		typeOutput, err := t.typeDescriberToTypeOutput(typesDefinitions[i])
		if err != nil {
			return err
		}

		out.Types = append(out.Types, typeOutput)
	}

	return nil
}

func (t *TemplateWriter) preparePublishedEvents(out *outputData, schemaResolver schema.Resolver) error {
	publishedEvents, err := schemaResolver.GetPublishedEvents()
	if err != nil {
		return fmt.Errorf("can't get published events to write: %w", err)
	}

	for i := range publishedEvents {
		eventOut, err := t.publishedEventToOutput(publishedEvents[i])
		if err != nil {
			return err
		}

		out.PublishedEvents = append(out.PublishedEvents, eventOut)
	}

	return nil
}

func (t *TemplateWriter) prepareConsumedEvents(out *outputData, schemaResolver schema.Resolver) error {
	consumedEvents, err := schemaResolver.GetConsumedEvents()
	if err != nil {
		return fmt.Errorf("can't get consumed events to write: %w", err)
	}

	for i := range consumedEvents {
		out.ConsumedEvents = append(out.ConsumedEvents, &consumedEventOutput{
			Name:        consumedEvents[i].Name(),
			Description: consumedEvents[i].Description(),
		})
	}

	return nil
}

func (t *TemplateWriter) publishedEventToOutput(event *types.PublishedEvent) (*publishedEventOutput, error) {
	out := &publishedEventOutput{
		Name:        event.Name(),
		Visibility:  event.Visibility().String(),
		Module:      event.Module(),
		Description: event.Description(),
	}

	var emojiPrefix string
	switch event.Visibility() {
	case types.EventPrivate:
		emojiPrefix = "🔒 "
	case types.EventProtected:
		emojiPrefix = "🔐 "
	case types.EventPublic:
		emojiPrefix = "🔓 "
	}
	out.Name = fmt.Sprintf("%s%s", emojiPrefix, out.Name)

	var (
		eventBody jsonc.MapSlice
		err       error
	)

	eventBody, err = t.createEventExampleMapItem(event.Name(), "attributes", event.Attributes(), eventBody)
	if err != nil {
		return nil, err
	}

	eventBody, err = t.createEventExampleMapItem(event.Name(), "entities", event.Entities(), eventBody)
	if err != nil {
		return nil, err
	}

	if err := t.exampleEncoder.Encode(eventBody); err != nil {
		return nil, fmt.Errorf("can't encode event '%s' example: %w", event.Name(), err)
	}

	out.Example = t.exampleWriter.String()
	t.exampleWriter.Reset()

	return out, nil
}

func (t *TemplateWriter) createEventExampleMapItem(
	eventName, exampleKey string,
	typeDescriber types.TypeDescriber,
	eventBody jsonc.MapSlice,
) (jsonc.MapSlice, error) {
	example, err := t.typeDescriberToExample(true, typeDescriber)
	if err != nil {
		return eventBody, fmt.Errorf("can't create %s example of '%s' event: %w", exampleKey, eventName, err)
	}

	eventBody = append(eventBody, jsonc.MapItem{
		Key:   exampleKey,
		Value: example,
	})

	return eventBody, nil
}

func (t *TemplateWriter) typeDescriberToTypeOutput(typeDescriber types.TypeDescriber) (*typeOutput, error) {
	out := &typeOutput{
		Name:        typeDescriber.Name(),
		Type:        string(typeDescriber.Type()),
		Description: typeDescriber.Description(),
		Nullable:    typeDescriber.Nullable(),
	}

	scalarType, is := typeDescriber.(types.ScalarDescriber)
	if is {
		out.Example = fmt.Sprint(scalarType.Value())
		out.Format = scalarType.Format()

		enumValues := scalarType.Enum()
		lastIndex := len(enumValues) - 1

		for i := range enumValues {
			out.Enum = append(out.Enum, enumValue{
				Value:   fmt.Sprint(enumValues[i]),
				HasMore: i < lastIndex,
			})
		}
	} else {
		example, err := t.typeDescriberToExample(true, typeDescriber)
		if err != nil {
			return nil, fmt.Errorf("can't create example of type: %w", err)
		}

		if err := t.exampleEncoder.Encode(example); err != nil {
			return nil, fmt.Errorf("can't encode type '%s': %w", typeDescriber.Path(), err)
		}

		out.Example = t.exampleWriter.String()
		t.exampleWriter.Reset()
	}

	return out, nil
}

func (t *TemplateWriter) typeDescriberToExample(inRootLevel bool, typeDescriber types.TypeDescriber) (interface{}, error) {
	switch typeDescriber := typeDescriber.(type) {
	case types.ScalarDescriber:
		if !inRootLevel {
			var typeModifier string

			if typeDescriber.HasFormat() {
				typeModifier = fmt.Sprintf("%s(%s)", typeModifier, typeDescriber.Format())
			}

			if typeDescriber.HasEnum() {
				typeModifier = fmt.Sprintf("%s[%v]", typeModifier, t.formatEnum(typeDescriber.Enum()))
			}

			if comment, require := t.createComment(typeDescriber, typeModifier); require {
				return jsonc.NewCommentValue(comment, typeDescriber.Value()), nil
			}
		}

		return typeDescriber.Value(), nil
	case types.ArrayDescriber:
		items, err := t.typeDescriberToExample(false, typeDescriber.Items())
		if err != nil {
			return nil, err
		}

		result := []interface{}{items}

		if !inRootLevel {
			if comment, require := t.createComment(typeDescriber, ""); require {
				return jsonc.NewCommentValue(comment, result), nil
			}
		}

		return result, nil
	case types.ObjectDescriber:
		properties := typeDescriber.Properties()
		result := make(jsonc.MapSlice, len(properties))

		for i := range properties {
			property, err := t.typeDescriberToExample(false, properties[i])
			if err != nil {
				return nil, err
			}

			result[i] = jsonc.MapItem{
				Key:   properties[i].Name(),
				Value: property,
			}
		}

		if !inRootLevel {
			if comment, require := t.createComment(typeDescriber, ""); require {
				return jsonc.NewCommentValue(comment, result), nil
			}
		}

		return result, nil
	}

	return nil, nil
}

func (t *TemplateWriter) createComment(typeDescriber types.TypeDescriber, typeModifier string) (string, bool) {
	var identifier, nullable, description string

	refereceType, is := typeDescriber.(types.ReferenceDescriber)
	if is {
		identifier = refereceType.Reference()
		typeModifier = ""
	} else {
		identifier = typeDescriber.Type()
	}

	if typeDescriber.Nullable() {
		nullable = "|null"
	}

	if len(typeDescriber.Description()) > 0 {
		description = fmt.Sprintf(": %s", typeDescriber.Description())
	}

	comment := fmt.Sprintf("%s%s%s", typeModifier, nullable, description)
	if len(comment) < 1 {
		return "", false
	}

	return fmt.Sprintf("%s%s", identifier, comment), true
}

func (t *TemplateWriter) formatEnum(enum []interface{}) string {
	lastIndex := len(enum) - 1
	if lastIndex < 0 {
		return ""
	}

	stringBuilder := &strings.Builder{}
	for i := range enum {
		stringBuilder.WriteString(fmt.Sprint(enum[i]))

		if i < lastIndex {
			stringBuilder.WriteRune(',')
		}
	}

	return stringBuilder.String()
}
