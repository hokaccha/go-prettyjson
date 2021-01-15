// Package prettyjson provides JSON pretty print.
package prettyjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	_ "unsafe" // go:linkname

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

	return f.Format(data)
}

// Format formats JSON string.
func (f *Formatter) Format(data []byte) ([]byte, error) {
	var v interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}

	buf := newBytesBuffer()
	f.pretty(buf, v, 1)
	return buf.Bytes(), nil
}

func (f *Formatter) setColor(buf *bytesBuffer, c *color.Color) {
	if f.DisabledColor || c == nil {
		return
	}
	setWriter(c, buf)
}

func (f *Formatter) unsetColor(buf *bytesBuffer, c *color.Color) {
	if f.DisabledColor || c == nil {
		return
	}
	unsetWriter(c, buf)
}

func (f *Formatter) sprintColor(c *color.Color, s string) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprint(s)
	}
	return c.SprintFunc()(s)
}

//go:linkname setWriter github.com/fatih/color.(*Color).setWriter
func setWriter(*color.Color, io.Writer) *color.Color

//go:linkname unsetWriter github.com/fatih/color.(*Color).unsetWriter
func unsetWriter(*color.Color, io.Writer)

func (f *Formatter) pretty(buf *bytesBuffer, v interface{}, depth int) {
	switch val := v.(type) {
	case string:
		f.setColor(buf, f.StringColor)
		f.writeString(buf, val)
		f.unsetColor(buf, f.StringColor)
	case float64:
		f.setColor(buf, f.NumberColor)
		buf.writeFloat64(val)
		f.unsetColor(buf, f.NumberColor)
	case json.Number:
		f.setColor(buf, f.NumberColor)
		buf.WriteString(val.String())
		f.unsetColor(buf, f.NumberColor)
	case bool:
		f.setColor(buf, f.BoolColor)
		if val {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
		f.unsetColor(buf, f.BoolColor)
	case nil:
		f.setColor(buf, f.NullColor)
		buf.WriteString("null")
		f.unsetColor(buf, f.NullColor)
	case map[string]interface{}:
		f.writeMap(buf, val, depth)
	case []interface{}:
		f.writeArray(buf, val, depth)
	}
}

func (f *Formatter) writeString(buf *bytesBuffer, s string) {
	if f.StringMaxLength != 0 {
		if r := []rune(s); len(r) >= f.StringMaxLength {
			s = string(r[0:f.StringMaxLength]) + "..."
		}
	}
	buf.writeString(s)
}

func (f *Formatter) writeMap(buf *bytesBuffer, m map[string]interface{}, depth int) {
	if len(m) == 0 {
		buf.WriteString("{}")
		return
	}

	keys := make([]string, len(m))
	var i int
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	buf.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.Newline)
		f.appendIndent(buf, depth)
		f.setColor(buf, f.KeyColor)
		buf.writeString(key)
		f.unsetColor(buf, f.KeyColor)
		buf.WriteByte(':')
		if f.Newline != "" {
			buf.WriteByte(' ')
		}
		f.pretty(buf, m[key], depth+1)
	}
	buf.WriteString(f.Newline)
	f.appendIndent(buf, depth-1)
	buf.WriteByte('}')
}

func (f *Formatter) writeArray(buf *bytesBuffer, a []interface{}, depth int) {
	if len(a) == 0 {
		buf.WriteString("[]")
		return
	}

	buf.WriteByte('[')
	for i, val := range a {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(f.Newline)
		f.appendIndent(buf, depth)
		f.pretty(buf, val, depth+1)
	}
	buf.WriteString(f.Newline)
	f.appendIndent(buf, depth-1)
	buf.WriteByte(']')
}

const spaces = "                                                                "

func (f *Formatter) appendIndent(buf *bytesBuffer, depth int) {
	if n := f.Indent * depth; n > 0 {
		for n > len(spaces) {
			buf.Write([]byte(spaces))
			n -= len(spaces)
		}
		buf.Write([]byte(spaces)[:n])
	}
}

// Marshal JSON data with default options.
func Marshal(v interface{}) ([]byte, error) {
	return NewFormatter().Marshal(v)
}

// Format JSON string with default options.
func Format(data []byte) ([]byte, error) {
	return NewFormatter().Format(data)
}

type bytesBuffer struct {
	bytes.Buffer
	enc     *json.Encoder
	scratch [64]byte
}

func newBytesBuffer() *bytesBuffer {
	var buf bytesBuffer
	buf.enc = json.NewEncoder(&buf.Buffer)
	buf.enc.SetEscapeHTML(false)
	return &buf
}

func (buf *bytesBuffer) writeString(str string) {
	if buf.enc.Encode(str) == nil {
		buf.Truncate(len(buf.Bytes()) - 1)
	}
}

func (buf *bytesBuffer) writeFloat64(f float64) {
	buf.Write(strconv.AppendFloat(buf.scratch[:0], f, 'f', -1, 64))
}
