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
			input: &commenterRetriverStub{
				comment: "This is a int",
				val:     67,
			},
			expectedWithComment: "67 // This is a int",
		},
		{
			name: "should add comment to string",
			input: &commenterRetriverStub{
				comment: "This is a string",
				val:     "yes!, comments",
			},
			expectedWithComment: `"yes!, comments" // This is a string`,
		},
		{
			name: "should add comment to map",
			input: map[string]interface{}{
				"key1": &commenterRetriverStub{
					comment: "Comment in key1!",
					val:     867.123,
				},
				"key2": &commenterRetriverStub{
					comment: "Its' hard to overstat my satisfaction",
					val: map[string]interface{}{
						"hasComment": false,
						"triumph": &commenterRetriverStub{
							comment: "Huge Success",
							val:     "I'm making a note here",
						},
					},
				},
				"key3": "Simple string",
				"key4": &commenterRetriverStub{
					comment: "We do what we must, because we can",
					val:     true,
				},
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

type commenterRetriverStub struct {
	val     interface{}
	comment string
}

func (t *commenterRetriverStub) GetComment() string {
	return t.comment
}

func (t *commenterRetriverStub) GetValue() interface{} {
	return t.val
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
