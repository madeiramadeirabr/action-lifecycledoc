package schema

import (
	"errors"
	"fmt"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

// BasicResolver implement a simple resolver with type reference support only. This resolver does not deal with concurrency
type BasicResolver struct {
	project *types.Project

	types map[string]types.TypeDescriber
	// slice of types paths to keep declaration order
	typePaths []string

	publishedEvents map[string]*types.PublishedEvent
	// slice of published events names to keep declaration order
	publishedEventsNames []string

	consumedEvents map[string]*types.ConsumedEvent
	// slice of consumed events to keep declaration order
	consumedEventsNames []string

	// hasResolved indicates that the types have been resolved
	hasResolved bool

	// resolvedTypes stored resolved types in all levels to prevent duplicate work
	resolvedTypes map[string]types.TypeDescriber

	// resolvingTypes stores the types being resolved to identify recursive references.
	// we use maps for quick searches, the stored value doesn't matter
	resolvingTypes map[string]bool
}

func (b *BasicResolver) SetProject(name string) error {
	if b.project != nil {
		return errors.New("can't override current project")
	}

	project, err := types.NewProject(name)
	if err != nil {
		return err
	}

	b.project = project
	return nil
}

func (b *BasicResolver) AddConfluencePage(title, spaceKey, ancestorID string) error {
	if err := b.isValid(); err != nil {
		return err
	}

	if len(title) < 1 {
		title = fmt.Sprintf("Life Cycle Events: %s", b.project.Name())
	}

	page, err := types.NewConfluencePage(title, spaceKey, ancestorID)
	if err != nil {
		return err
	}

	b.project.Confluence().AddPage(page)
	return nil
}

func (b *BasicResolver) GetConfluence() (*types.Confluence, error) {
	if err := b.isValid(); err != nil {
		return nil, err
	}

	return b.project.Confluence(), nil
}

func (b *BasicResolver) AddType(t types.TypeDescriber) error {
	if err := b.isValid(); err != nil {
		return err
	}

	if _, exists := b.types[t.Path()]; exists {
		return fmt.Errorf("type '%s' has been duplicated", t.Path())
	}

	b.hasResolved = false

	b.types[string(t.Path())] = t
	b.typePaths = append(b.typePaths, t.Path())
	return nil
}

func (b *BasicResolver) GetTypes() ([]types.TypeDescriber, error) {
	if err := b.resolve(); err != nil {
		return nil, err
	}

	result := make([]types.TypeDescriber, len(b.typePaths))
	for i := range b.typePaths {
		result[i] = b.types[b.typePaths[i]]
	}

	return result, nil
}

func (b *BasicResolver) AddPublishedEvent(e *types.PublishedEvent) error {
	if err := b.isValid(); err != nil {
		return err
	}

	if _, exists := b.publishedEvents[e.Name()]; exists {
		return fmt.Errorf("published event '%s' has been duplicated", e.Name())
	}

	b.hasResolved = false

	b.publishedEvents[e.Name()] = e
	b.publishedEventsNames = append(b.publishedEventsNames, e.Name())
	return nil
}

func (b *BasicResolver) GetPublishedEvents() ([]*types.PublishedEvent, error) {
	if err := b.resolve(); err != nil {
		return nil, err
	}

	result := make([]*types.PublishedEvent, len(b.publishedEventsNames))
	for i := range b.publishedEventsNames {
		result[i] = b.publishedEvents[b.publishedEventsNames[i]]
	}

	return result, nil
}

func (b *BasicResolver) AddConsumedEvent(e *types.ConsumedEvent) error {
	if err := b.isValid(); err != nil {
		return err
	}

	if _, exists := b.consumedEvents[e.Name()]; exists {
		return fmt.Errorf("consumed event '%s' has been duplicated", e.Name())
	}

	b.consumedEvents[e.Name()] = e
	b.consumedEventsNames = append(b.consumedEventsNames, e.Name())
	return nil
}

func NewBasicResolver() *BasicResolver {
	return &BasicResolver{
		types:           make(map[string]types.TypeDescriber),
		publishedEvents: make(map[string]*types.PublishedEvent),
		consumedEvents:  make(map[string]*types.ConsumedEvent),
		resolvedTypes:   make(map[string]types.TypeDescriber),
		resolvingTypes:  make(map[string]bool),
	}
}

func (b *BasicResolver) isValid() error {
	if b.project == nil {
		return errors.New("schema not configured, please specify required fields")
	}

	return nil
}

func (b *BasicResolver) resolve() error {
	if err := b.isValid(); err != nil {
		return err
	}

	if b.hasResolved {
		return nil
	}

	for path := range b.types {
		resolvedType, err := b.getResolvedType(b.types[path])
		if err != nil {
			return err
		}

		b.types[path] = resolvedType
	}

	for name := range b.publishedEvents {
		attributesType, err := b.getResolvedType(b.publishedEvents[name].Attributes())
		if err != nil {
			return err
		}

		b.publishedEvents[name].SetAttributes(attributesType)

		entities, err := b.getResolvedType(b.publishedEvents[name].Entities())
		if err != nil {
			return err
		}

		b.publishedEvents[name].SetEntities(entities)
	}

	b.hasResolved = true
	return nil
}

func (b *BasicResolver) getResolvedType(t types.TypeDescriber) (types.TypeDescriber, error) {
	if resolved, exists := b.resolvedTypes[t.Path()]; exists {
		return resolved, nil
	}

	resolved, err := b.resolveType(t)
	if err != nil {
		return nil, err
	}

	b.resolvedTypes[resolved.Path()] = resolved
	return resolved, nil
}

func (b *BasicResolver) resolveType(t types.TypeDescriber) (types.TypeDescriber, error) {
	switch t.Type() {
	case types.ReferenceType:
		return b.resolveReferenceType(t)
	case types.ArrayType:
		return b.resolveArrayType(t)
	case types.ObjectType:
		return b.resolveObjectType(t)
	}

	return t, nil
}

func (b *BasicResolver) resolveArrayType(t types.TypeDescriber) (types.TypeDescriber, error) {
	arrayType, is := t.(*types.Array)
	if !is {
		return nil, fmt.Errorf("definition '%s' is not of array type", t.Path())
	}

	itemsType, err := b.getResolvedType(arrayType.Items())
	if err != nil {
		return nil, fmt.Errorf("can't resolve items from '%s': %w", arrayType.Path(), err)
	}

	arrayType.SetItems(itemsType)
	return arrayType, nil
}

func (b *BasicResolver) resolveObjectType(t types.TypeDescriber) (types.TypeDescriber, error) {
	objectType, is := t.(*types.Object)
	if !is {
		return nil, fmt.Errorf("definition '%s' is not of object type", t.Path())
	}

	properties := objectType.Properties()
	for i := range properties {
		propertyType, err := b.getResolvedType(properties[i])
		if err != nil {
			return nil, err
		}

		properties[i] = propertyType
	}

	objectType.SetProperties(properties)
	return objectType, nil
}

func (b *BasicResolver) resolveReferenceType(t types.TypeDescriber) (types.TypeDescriber, error) {
	referenceType, err := typeDescriberToRealType[*types.Reference](t)
	if err != nil {
		return nil, err
	}

	targetType, exists := b.types[referenceType.Reference()]
	if !exists {
		return nil, fmt.Errorf("definition '%s' referenced in '%s' not found in declared types", referenceType.Reference(), referenceType.Path())
	}

	if _, exists := b.resolvingTypes[referenceType.Path()]; exists {
		return nil, fmt.Errorf("recursive reference dected for definition '%s'", referenceType.Path())
	}
	b.resolvingTypes[referenceType.Path()] = true

	targetType, err = b.getResolvedType(targetType)
	if err != nil {
		return nil, fmt.Errorf("can't resolve '%s' reference: %w", referenceType.Path(), err)
	}

	delete(b.resolvingTypes, referenceType.Path())

	// Recreate the type definition to override generic infomation
	switch targetType.Type() {
	case types.ScalarType:
		scalarType, err := typeDescriberToRealType[*types.Scalar](targetType)
		if err != nil {
			return nil, err
		}

		return types.NewScalar(
			referenceType.Name(),
			referenceType.Path(),
			getFirstNonEmptyString(referenceType.Description(), scalarType.Description()),
			scalarType.Nullable(),
			scalarType.Format(),
			scalarType.Enum(),
			scalarType.Value(),
		)
	case types.ArrayType:
		arrayType, err := typeDescriberToRealType[*types.Array](targetType)
		if err != nil {
			return nil, err
		}

		return types.NewArray(
			referenceType.Name(),
			referenceType.Path(),
			getFirstNonEmptyString(referenceType.Description(), arrayType.Description()),
			arrayType.Nullable(),
			arrayType.Items(),
		)
	case types.ObjectType:
		objectType, err := typeDescriberToRealType[*types.Object](targetType)
		if err != nil {
			return nil, err
		}

		return types.NewObject(
			referenceType.Name(),
			referenceType.Path(),
			getFirstNonEmptyString(referenceType.Description(), objectType.Description()),
			objectType.Nullable(),
			objectType.Properties(),
		)
	default:
		return nil, fmt.Errorf("type '%s' of defintion '%s' is not supported", targetType.Type(), targetType.Path())
	}
}

func getFirstNonEmptyString(strings ...string) string {
	for i := range strings {
		if len(strings[i]) > 0 {
			return strings[i]
		}
	}

	return ""
}

type typeDescriberRealTypes interface {
	*types.Scalar | *types.Reference | *types.Array | *types.Object
}

func typeDescriberToRealType[T typeDescriberRealTypes](t types.TypeDescriber) (T, error) {
	result, is := t.(T)
	if !is {
		return nil, fmt.Errorf("definition '%s' is not of '%T' type", t.Path(), result)
	}

	return result, nil
}
