package schema_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
)

func TestShouldParseTypes(t *testing.T) {
	testCases := []struct {
		name    string
		newType func(t *testing.T, name, path string) (types.TypeDescriber, error)
	}{
		{
			name: "StringType",
			newType: func(t *testing.T, name, path string) (types.TypeDescriber, error) {
				return types.NewScalar(name, path, "", false, "", nil, "banana")
			},
		},
		{
			name: "ArrayType",
			newType: func(t *testing.T, name, path string) (types.TypeDescriber, error) {
				itemsType, err := types.NewScalar(
					"IntType",
					"#/types/IntType",
					"",
					true,
					"",
					nil,
					nil,
				)
				if err != nil {
					return nil, err
				}

				return types.NewArray(
					name,
					path,
					"",
					false,
					itemsType,
				)
			},
		},
		{
			name: "ObjectType",
			newType: func(t *testing.T, name, path string) (types.TypeDescriber, error) {
				boolType, err := types.NewScalar(
					"BoolType",
					"#/types/BoolType",
					"",
					false,
					"Sim",
					nil,
					true,
				)
				if err != nil {
					return nil, err
				}

				return types.NewObject(
					name,
					path,
					"",
					false,
					[]types.TypeDescriber{boolType},
				)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resolver := schema.NewBasicResolver()
			resolver.SetProject("test project")

			expectedPath := fmt.Sprintf("#/types/%s", testCase.name)

			typeDef, err := testCase.newType(t, testCase.name, expectedPath)
			assertNoError(t, err)

			assertNoError(t, resolver.AddType(typeDef))

			types, err := resolver.GetTypes()
			assertNoError(t, err)
			assertLength(t, 1, types)

			if types[0] == nil {
				t.Fatalf("expected non nil resolved type, received nil")
			}

			assertTypeName(t, testCase.name, types[0])
			assertTypePath(t, expectedPath, types[0])
		})
	}
}

func TestShouldResolveTypeReferences(t *testing.T) {
	resolver := schema.NewBasicResolver()
	resolver.SetProject("test")

	stringType, err := types.NewScalar(
		"StringType",
		"#/types/StringType",
		"Description of string",
		false,
		"",
		[]interface{}{
			"option1",
			"option2",
		},
		"option1",
	)
	assertNoError(t, err)
	assertNoError(t, resolver.AddType(stringType))

	refToString, err := types.NewReference(
		"StringRef",
		"#/types/StringRef",
		"Description of string reference",
		false,
		stringType.Path(),
	)
	assertNoError(t, err)
	assertNoError(t, resolver.AddType(refToString))

	arrayType, err := types.NewArray(
		"StringArray",
		"#/types/StringArray",
		"Description of string array",
		true,
		refToString,
	)
	assertNoError(t, err)
	assertNoError(t, resolver.AddType(arrayType))

	idType, err := types.NewScalar(
		"id",
		"#/types/ObjectType/properties/id",
		"Description of ID type",
		false,
		"",
		nil,
		"123456",
	)
	assertNoError(t, err)

	arrayRef, err := types.NewReference(
		"strings",
		"#/types/ObjectType/properties/strings",
		"Description of string array reference",
		false,
		arrayType.Path(),
	)
	assertNoError(t, err)

	objectType, err := types.NewObject(
		"ObjectType",
		"#/types/ObjectType",
		"Description of object",
		false,
		[]types.TypeDescriber{
			idType,
			arrayRef,
		},
	)
	assertNoError(t, err)
	assertNoError(t, resolver.AddType(objectType))

	typesDefinitions, err := resolver.GetTypes()
	assertNoError(t, err)
	assertLength(t, 4, typesDefinitions)

	typesMap := typeDescriberSliceToMap(typesDefinitions)
	for i := range typesDefinitions {
		typesMap[typesDefinitions[i].Path()] = typesDefinitions[i]
	}

	assertTypeExistInMap(t, objectType.Path(), typesMap)
	assertTypeExistInMap(t, arrayType.Path(), typesMap)
	assertTypeExistInMap(t, refToString.Path(), typesMap)
	assertTypeExistInMap(t, stringType.Path(), typesMap)

	objectResolvedType := typeDefintionToRealType[*types.Object](t, typesMap[objectType.Path()])
	properties := typeDescriberSliceToMap(objectResolvedType.Properties())

	assertTypeExistInMap(t, idType.Path(), properties)
	assertTypeExistInMap(t, arrayRef.Path(), properties)

	assertTypeDescription(t, idType.Description(), properties[idType.Path()])

	idScalarType := typeDefintionToRealType[*types.Scalar](t, properties[idType.Path()])

	if idScalarType.Value() != idType.Value() {
		t.Errorf("expected '%s' id property value, received '%s'", idType.Value(), idScalarType.Value())
	}

	assertTypeDescription(t, arrayRef.Description(), properties[arrayRef.Path()])

	arrayRefType := typeDefintionToRealType[*types.Array](t, properties[arrayRef.Path()])
	arrayRefItemsType := typeDefintionToRealType[*types.Scalar](t, arrayRefType.Items())

	assertTypeName(t, refToString.Name(), arrayRefItemsType)
	assertTypePath(t, refToString.Path(), arrayRefItemsType)
	assertTypeDescription(t, refToString.Description(), arrayRefItemsType)
	assertScalarValue(t, stringType.Value(), arrayRefItemsType)
}

func TestShouldRegisterEvent(t *testing.T) {
	resolver := schema.NewBasicResolver()
	resolver.SetProject("test events")

	intType, err := types.NewScalar(
		"id",
		"#/types/ObjectType/id",
		"",
		false,
		"",
		nil,
		10,
	)
	assertNoError(t, err)

	objectType, err := types.NewObject(
		"ObjectType",
		"#/types/ObjectType",
		"",
		false,
		[]types.TypeDescriber{
			intType,
		},
	)
	assertNoError(t, err)
	assertNoError(t, resolver.AddType(objectType))

	refType, err := types.NewReference(
		"object",
		"#/events/published/SIMPLE_EVENT/properties/object",
		"",
		false,
		objectType.Path(),
	)
	assertNoError(t, err)

	stringType, err := types.NewScalar(
		"id",
		"#/events/published/SIMPLE_EVENT/entities/id",
		"",
		false,
		"",
		nil,
		"banana",
	)
	assertNoError(t, err)

	entitiesType, err := types.NewObject(
		"ObjectType",
		"#/events/published/SIMPLE_EVENT/entities",
		"",
		false,
		[]types.TypeDescriber{
			stringType,
		},
	)
	assertNoError(t, err)

	event, err := types.NewPublishdEvent(
		"SIMPLE_EVENT",
		types.EventPublic,
		"",
		"",
		refType,
		entitiesType,
	)
	assertNoError(t, err)

	assertNoError(t, resolver.AddPublishedEvent(event))
	events, err := resolver.GetPublishedEvents()
	assertNoError(t, err)

	if length := len(events); length != 1 {
		t.Fatalf("expected '1' events, received '%d'", length)
	}

	assertTypeName(t, entitiesType.Name(), events[0].Entities())
	assertTypePath(t, entitiesType.Path(), events[0].Entities())

	assertTypeName(t, refType.Name(), events[0].Attributes())
	assertTypePath(t, refType.Path(), events[0].Attributes())

	attributesType := typeDefintionToRealType[*types.Object](t, events[0].Attributes())
	properties := typeDescriberSliceToMap(attributesType.Properties())

	assertTypeExistInMap(t, "#/types/ObjectType/id", properties)
	assertScalarValue(t, 10, properties["#/types/ObjectType/id"])
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, received '%s'", err)
	}
}

func assertLength(t *testing.T, expectedLength int, types []types.TypeDescriber) {
	t.Helper()

	if total := len(types); total != expectedLength {
		t.Fatalf("expected '%d' resolved types, received '%d'", expectedLength, total)
	}
}

func assertTypeExistInMap(t *testing.T, key string, types map[string]types.TypeDescriber) {
	t.Helper()

	if _, exists := types[key]; !exists {
		t.Fatalf("not found '%s' in resolved types", key)
	}
}

func typeDescriberSliceToMap(t []types.TypeDescriber) map[string]types.TypeDescriber {
	result := make(map[string]types.TypeDescriber)
	for i := range t {
		result[t[i].Path()] = t[i]
	}

	return result
}

func assertTypeDescription(t *testing.T, expected string, typeDef types.TypeDescriber) {
	t.Helper()

	assertString(t, expected, typeDef.Description())
}

func assertTypeName(t *testing.T, expected string, typeDef types.TypeDescriber) {
	t.Helper()

	assertString(t, expected, typeDef.Name())
}

func assertTypePath(t *testing.T, expected string, typeDef types.TypeDescriber) {
	t.Helper()

	assertString(t, expected, typeDef.Path())
}

func assertString(t *testing.T, expected, value string) {
	t.Helper()

	if value != expected {
		t.Errorf("expected '%s', received '%s'", expected, value)
	}
}

func assertScalarValue(t *testing.T, expected interface{}, typeDef types.TypeDescriber) {
	t.Helper()

	scalar, is := typeDef.(*types.Scalar)
	if !is {
		t.Fatalf("definition '%s' is not a scalar type", typeDef.Path())
	}

	if value := scalar.Value(); !reflect.DeepEqual(value, expected) {
		t.Errorf("expected '%v' value, received '%v'", expected, value)
	}
}

func typeDefintionToRealType[T *types.Scalar | *types.Array | *types.Object](t *testing.T, typeDef types.TypeDescriber) T {
	t.Helper()

	result, is := typeDef.(T)
	if !is {
		t.Fatalf("definition '%s' expected type '%T', received '%T': %#v", typeDef.Path(), result, typeDef, typeDef)
	}

	return result
}
