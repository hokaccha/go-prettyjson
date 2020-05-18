// Package prettyjson provides JSON pretty print.
package prettyjson

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/singlemusic/go-ordered-json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Formatter is a struct to format JSON data. `color` is github.com/fatih/color: https://github.com/fatih/color
type Formatter struct {
	// JSON key color. Default is `color.New(color.FgBlue, color.Bold)`.
	KeyColor *color.Color

	// JSON string value color. Default is `color.New(color.FgGreen, color.Bold)`.
	StringColor *color.Color

	// JSON boolean value color. Default is `color.New(color.FgYellow, color.Bold)`.
	BoolColor *color.Color

	// JSON number value color. Default is `color.New(color.FgCyan, color.Bold)`.
	NumberColor *color.Color

	// JSON null value color. Default is `color.New(color.FgBlack, color.Bold)`.
	NullColor *color.Color

	// Max length of JSON string value. When the value is 1 and over, string is truncated to length of the value.
	// Default is 0 (not truncated).
	StringMaxLength int

	// Boolean to disable color. Default is false.
	DisabledColor bool

	// Indent space number. Default is 2.
	Indent int

	// Newline string. To print without new lines set it to empty string. Default is \n.
	Newline string
}

// NewFormatter returns a new formatter with following default values.
func NewFormatter() *Formatter {
	return &Formatter{
		KeyColor:        color.New(color.FgBlue, color.Bold),
		StringColor:     color.New(color.FgGreen, color.Bold),
		BoolColor:       color.New(color.FgYellow, color.Bold),
		NumberColor:     color.New(color.FgCyan, color.Bold),
		NullColor:       color.New(color.FgBlack, color.Bold),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          2,
		Newline:         "\n",
	}
}

// Marshal marshals and formats JSON data.
func (f *Formatter) Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	rt := reflect.TypeOf(v)
	kind := rt.Kind()
	if reflect.Ptr == kind {
		kind = rt.Elem().Kind()
	}
	switch kind {
	case reflect.Slice:
		return f.FormatArray(data)
	case reflect.Array:
		return f.FormatArray(data)
	case reflect.Struct:
		return f.Format(data)
	case reflect.Interface:
		return f.Format(data)
	case reflect.Map:
		return f.Format(data)
	default:
		return f.FormatLiteral(data)
	}
}

// Format formats JSON string.
func (f *Formatter) Format(data []byte) ([]byte, error) {
	var om = ordered.NewOrderedMap()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(om); err != nil {
		return nil, err
	}
	return []byte(f.pretty(om, 1)), nil
}

func (f *Formatter) FormatArray(data []byte) ([]byte, error) {
	var array = make([]ordered.OrderedMap, 0)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&array); err != nil {
		return nil, err
	}
	return []byte(f.pretty(array, 1)), nil
}

func (f *Formatter) FormatLiteral(data []byte) ([]byte, error) {
	var v interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	return []byte(f.pretty(v, 1)), nil
}

func (f *Formatter) sprintfColor(c *color.Color, format string, args ...interface{}) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprintf(format, args...)
	}
	return c.SprintfFunc()(format, args...)
}

func (f *Formatter) sprintColor(c *color.Color, s string) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprint(s)
	}
	return c.SprintFunc()(s)
}

func (f *Formatter) pretty(v interface{}, depth int) string {
	switch val := v.(type) {
	case string:
		return f.processString(val)
	case float64:
		return f.sprintColor(f.NumberColor, strconv.FormatFloat(val, 'f', -1, 64))
	case json.Number:
		return f.sprintColor(f.NumberColor, string(val))
	case bool:
		return f.sprintColor(f.BoolColor, strconv.FormatBool(val))
	case nil:
		return f.sprintColor(f.NullColor, "null")
	case *ordered.OrderedMap:
		return f.processOrderedMapPtr(val, depth)
	case ordered.OrderedMap:
		return f.processOrderedMap(val, depth)
	case *list.List:
		return f.processListPtr(val, depth)
	case list.List:
		return f.processList(val, depth)
	case []ordered.OrderedMap:
		return f.processOrderedMapArray(val, depth)
	case map[string]interface{}:
		return f.processMap(val, depth)
	case []interface{}:
		return f.processArray(val, depth)
	}
	return ""
}

func (f *Formatter) processString(s string) string {
	r := []rune(s)
	if f.StringMaxLength != 0 && len(r) >= f.StringMaxLength {
		s = string(r[0:f.StringMaxLength]) + "..."
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(s)
	s = string(buf.Bytes())
	s = strings.TrimSuffix(s, "\n")

	return f.sprintColor(f.StringColor, s)
}

func (f *Formatter) processMap(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}
	keys := []string{}

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		k := f.sprintfColor(f.KeyColor, `"%s"`, key)
		v := f.pretty(val, depth+1)

		valueIndent := " "
		if f.Newline == "" {
			valueIndent = ""
		}
		row := fmt.Sprintf("%s%s:%s%s", nextIndent, k, valueIndent, v)
		rows = append(rows, row)
	}

	return fmt.Sprintf("{%s%s%s%s}", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processOrderedMapPtr(m *ordered.OrderedMap, depth int) string {
	_, notEmpty := m.EntriesIter()()
	if !notEmpty {
		return "{}"
	}
	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}
	keys := []string{}

	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		keys = append(keys, pair.Key)
		val := pair.Value
		k := f.sprintfColor(f.KeyColor, `"%s"`, pair.Key)
		v := f.pretty(val, depth+1)

		valueIndent := " "
		if f.Newline == "" {
			valueIndent = ""
		}
		row := fmt.Sprintf("%s%s:%s%s", nextIndent, k, valueIndent, v)
		rows = append(rows, row)
	}
	return fmt.Sprintf("{%s%s%s%s}", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processOrderedMap(m ordered.OrderedMap, depth int) string {
	_, notEmpty := m.EntriesIter()()
	if !notEmpty {
		return "{}"
	}
	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}
	keys := []string{}

	iter := m.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		keys = append(keys, pair.Key)
		val := pair.Value
		k := f.sprintfColor(f.KeyColor, `"%s"`, pair.Key)
		v := f.pretty(val, depth+1)

		valueIndent := " "
		if f.Newline == "" {
			valueIndent = ""
		}
		row := fmt.Sprintf("%s%s:%s%s", nextIndent, k, valueIndent, v)
		rows = append(rows, row)
	}
	return fmt.Sprintf("{%s%s%s%s}", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processArray(a []interface{}, depth int) string {
	if len(a) == 0 {
		return "[]"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	for _, val := range a {
		c := f.pretty(val, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}
	return fmt.Sprintf("[%s%s%s%s]", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processListPtr(list *list.List, depth int) string {
	if list.Len() == 0 {
		return "[]"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	for e := list.Front(); e != nil; e = e.Next() {
		c := f.pretty(e.Value, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}
	return fmt.Sprintf("[%s%s%s%s]", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processList(list list.List, depth int) string {
	if list.Len() == 0 {
		return "[]"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	for e := list.Front(); e != nil; e = e.Next() {
		c := f.pretty(e.Value, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}
	return fmt.Sprintf("[%s%s%s%s]", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processOrderedMapArray(a []ordered.OrderedMap, depth int) string {
	if len(a) == 0 {
		return "[]"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	for _, val := range a {
		c := f.pretty(val, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}
	return fmt.Sprintf("[%s%s%s%s]", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) generateIndent(depth int) string {
	return strings.Repeat(" ", f.Indent*depth)
}

// Marshal JSON data with default options.
func Marshal(v interface{}) ([]byte, error) {
	return NewFormatter().Marshal(v)
}

// Format JSON string with default options.
func Format(data []byte) ([]byte, error) {
	return NewFormatter().Format(data)
}
