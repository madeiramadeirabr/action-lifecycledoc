// jsonc package provides a simple JSONC (JSON with comments) encoder
package jsonc

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type CommentRetriver interface {
	GetComment() string
	GetValue() interface{}
}

type FormattedEncoder struct {
	writer io.Writer
}

func NewEncoder(w io.Writer) *FormattedEncoder {
	return &FormattedEncoder{
		writer: w,
	}
}

/*
Encode the "v" value formmated for humans.

Supported types:

  - nil
  - numeric (int, int8, int16, int32, uint, uint8, uint32, uint64, float32, float64)
  - bool
  - string
  - map[string]interface{}
  - []interface{}
  - CommentRetriver
*/
func (f *FormattedEncoder) Encode(v interface{}) error {
	return f.writeValue(v, 1, false)
}

func (f *FormattedEncoder) writeValue(v interface{}, level int, addComma bool) error {
	var err error

	if v == nil {
		err = f.write("null")
		err = f.writeCommaIfNoErrorOccured(err, addComma)
	} else {
		switch t := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			err = f.writefWithComma(addComma, "%d", v)
		case float32, float64:
			err = f.writefWithComma(addComma, "%v", v)
		case bool:
			err = f.writefWithComma(addComma, "%t", v)
		case string:
			err = f.writefWithComma(addComma, `"%s"`, v)
		case CommentRetriver:
			err = f.writeValue(t.GetValue(), level, addComma)

			if err == nil && len(t.GetComment()) > 0 {
				err = f.writef(" // %s", t.GetComment())
			}
		case map[string]interface{}:
			err = f.write("{")
			if err != nil {
				break
			}

			// Sort keys in asc
			var keys []string
			for k := range t {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			lastIndex := len(keys) - 1

			for index, key := range keys {
				err = f.writeNewLine(level)
				if err != nil {
					break
				}

				err = f.writef(`"%s": `, key)
				if err != nil {
					break
				}

				err = f.writeValue(t[key], level+1, index < lastIndex)
				if err != nil {
					break
				}
			}

			if err == nil {
				err = f.writeCloseBlock("}", level, addComma)
			}
		case []interface{}:
			err = f.write("[")
			if err != nil {
				break
			}

			lastIndex := len(t) - 1

			for index := range t {
				err = f.writeNewLine(level)
				if err != nil {
					break
				}

				err = f.writeValue(t[index], level+1, index < lastIndex)
				if err != nil {
					break
				}
			}

			if err == nil {
				err = f.writeCloseBlock("]", level, addComma)
			}
		default:
			err = fmt.Errorf("the type '%T' is not supported", v)
		}
	}

	return err
}

func (f *FormattedEncoder) write(val string) (err error) {
	_, err = f.writer.Write([]byte(val))
	return
}

func (f *FormattedEncoder) writef(format string, vals ...interface{}) (err error) {
	_, err = fmt.Fprintf(f.writer, format, vals...)
	return
}

func (f *FormattedEncoder) writeCommaIfNoErrorOccured(err error, addComma bool) error {
	if err != nil {
		return err
	}

	if !addComma {
		return nil
	}

	return f.write(",")
}

func (f *FormattedEncoder) writefWithComma(addComma bool, format string, vals ...interface{}) error {
	err := f.writef(format, vals...)
	err = f.writeCommaIfNoErrorOccured(err, addComma)

	return err
}

func (f *FormattedEncoder) writeNewLine(level int) error {
	return f.writef("\n%s", strings.Repeat("\t", level))
}

func (f *FormattedEncoder) writeCloseBlock(closingSequence string, currentLevel int, addComma bool) error {
	err := f.writef("\n%s%s", strings.Repeat("\t", currentLevel-1), closingSequence)
	err = f.writeCommaIfNoErrorOccured(err, addComma)

	return err
}
