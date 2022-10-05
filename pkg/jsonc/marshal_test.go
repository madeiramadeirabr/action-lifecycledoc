package jsonc_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/jsonc"
)

func TestEncodePrimiteValues(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "should encode nil value",
			input:    nil,
			expected: "null",
		},
		{
			name:     "should encode int value",
			input:    99962,
			expected: "99962",
		},
		{
			name:     "should encode float value",
			input:    22.33,
			expected: "22.33",
		},
		{
			name:     "should encode bool value",
			input:    false,
			expected: "false",
		},
		{
			name:     "should encode string value",
			input:    "yes!",
			expected: `"yes!"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assertEncode(t, false, testCase.input, testCase.expected)
		})
	}
}

func TestEncodeMapValues(t *testing.T) {
	input := map[string]interface{}{
		"key1": "a",
		"key2": 10,
		"key3": map[string]interface{}{
			"subkey1": true,
			"subkey2": 33.22,
		},
	}

	expected := `{
	"key1": "a",
	"key2": 10,
	"key3": {
		"subkey1": true,
		"subkey2": 33.22
	}
}`
	assertEncode(t, false, input, expected)
}

func TestEncodeMapSlice(t *testing.T) {
	input := jsonc.MapSlice{
		jsonc.MapItem{
			Key:   "Z",
			Value: 10,
		},
		jsonc.MapItem{
			Key:   "A",
			Value: "Yes!",
		},
		jsonc.MapItem{
			Key:   "0",
			Value: false,
		},
	}

	expected := `{
	"Z": 10,
	"A": "Yes!",
	"0": false
}`

	assertEncode(t, false, input, expected)
}

func TestEncodeArrayValues(t *testing.T) {
	input := map[string]interface{}{
		"array": []interface{}{
			56,
			85,
			2,
		},
	}

	expect := `{
	"array": [
		56,
		85,
		2
	]
}`

	assertEncode(t, false, input, expect)
}

func TestEncodeCommentableValues(t *testing.T) {
	testCases := []struct {
		name                string
		input               interface{}
		expectedWithComment string
	}{
		{
			name: "should add comment to int",
			input: jsonc.NewCommentValue(
				"This is a int",
				67,
			),
			expectedWithComment: "67 // This is a int",
		},
		{
			name: "should add comment to string",
			input: jsonc.NewCommentValue(
				"This is a string",
				"yes!, comments",
			),
			expectedWithComment: `"yes!, comments" // This is a string`,
		},
		{
			name: "should add comment to map",
			input: map[string]interface{}{
				"key1": jsonc.NewCommentValue(
					"Comment in key1!",
					867.123,
				),
				"key2": jsonc.NewCommentValue(
					"Its' hard to overstat my satisfaction",
					map[string]interface{}{
						"hasComment": false,
						"triumph": jsonc.NewCommentValue(
							"Huge Success",
							"I'm making a note here",
						),
					},
				),
				"key3": "Simple string",
				"key4": jsonc.NewCommentValue(
					"We do what we must, because we can",
					true,
				),
			},
			expectedWithComment: `{
	"key1": 867.123, // Comment in key1!
	"key2": {
		"hasComment": false,
		"triumph": "I'm making a note here" // Huge Success
	}, // Its' hard to overstat my satisfaction
	"key3": "Simple string",
	"key4": true // We do what we must, because we can
}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s with comment", testCase.name), func(t *testing.T) {
			assertEncode(t, true, testCase.input, testCase.expectedWithComment)
		})
	}
}

func assertCanEncode(t *testing.T, writer io.Writer, withComments bool, input interface{}) {
	t.Helper()

	if err := jsonc.NewEncoder(writer).Encode(input); err != nil {
		t.Fatalf("can't encode value '%#v': %s", input, err)
	}
}

func assertString(t *testing.T, expected, result string) {
	t.Helper()

	if result != expected {
		t.Errorf("expected '%s', but receive '%s'", expected, result)
	}
}

func assertEncode(t *testing.T, withComments bool, input interface{}, expected string) {
	t.Helper()

	result := &strings.Builder{}

	assertCanEncode(t, result, withComments, input)
	assertString(t, expected, result.String())
}
