package prettyjson

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Formatter struct {
	KeyColor        *color.Color
	StringColor     *color.Color
	BoolColor       *color.Color
	NumberColor     *color.Color
	NullColor       *color.Color
	StringMaxLength int
	DisabledColor   bool
	Indent          int
}

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
	}
}

func (f *Formatter) Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return f.Format(data)
}

func (f *Formatter) Format(data []byte) ([]byte, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)

	if err != nil {
		return nil, err
	}

	s := f.pretty(v, 1)

	return []byte(s), nil
}

func (f *Formatter) sprintfColor(c *color.Color, format string, args ...interface{}) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprintf(format, args...)
	} else {
		return c.SprintfFunc()(format, args...)
	}
}

func (f *Formatter) pretty(v interface{}, depth int) string {
	switch val := v.(type) {
	case string:
		return f.processString(val)
	case float64:
		return f.sprintfColor(f.NumberColor, strconv.FormatFloat(val, 'f', -1, 64))
	case bool:
		return f.sprintfColor(f.BoolColor, strconv.FormatBool(val))
	case nil:
		return f.sprintfColor(f.NullColor, "null")
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

	return f.sprintfColor(f.StringColor, `"%s"`, s)
}

func (f *Formatter) processMap(m map[string]interface{}, depth int) string {
	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}
	keys := []string{}

	if len(m) == 0 {
		return "{}"
	}

	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		k := f.sprintfColor(f.KeyColor, `"%s"`, key)
		v := f.pretty(val, depth+1)
		row := fmt.Sprintf("%s%s: %s", nextIndent, k, v)
		rows = append(rows, row)
	}

	return fmt.Sprintf("{\n%s\n%s}", strings.Join(rows, ",\n"), currentIndent)
}

func (f *Formatter) processArray(a []interface{}, depth int) string {
	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	if len(a) == 0 {
		return "[]"
	}

	for _, val := range a {
		c := f.pretty(val, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(rows, ",\n"), currentIndent)
}

func (f *Formatter) generateIndent(depth int) string {
	return strings.Join(make([]string, f.Indent*depth+1), " ")
}

func Marshal(v interface{}) ([]byte, error) {
	return NewFormatter().Marshal(v)
}

func Format(data []byte) ([]byte, error) {
	return NewFormatter().Format(data)
}
