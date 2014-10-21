package prettyjson

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var KeyColor = color.New(color.FgBlue, color.Bold)
var StringColor = color.New(color.FgGreen, color.Bold)
var BoolColor = color.New(color.FgYellow, color.Bold)
var NumberColor = color.New(color.FgCyan, color.Bold)
var NullColor = color.New(color.FgBlack, color.Bold)
var DisabledColor = false
var Indent = 2

func MarshalPretty(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return Pretty(data)
}

func Pretty(data []byte) ([]byte, error) {
	var v interface{}
	err := json.Unmarshal(data, &v)

	if err != nil {
		return nil, err
	}

	s := pretty(v, 1)

	return []byte(s), nil
}

func sprintfColor(c *color.Color, format string, args ...interface{}) string {
	if DisabledColor || c == nil {
		return fmt.Sprintf(format, args...)
	} else {
		return c.SprintfFunc()(format, args...)
	}
}

func pretty(v interface{}, depth int) string {
	switch val := v.(type) {
	case string:
		return sprintfColor(StringColor, `"%s"`, val)
	case float64:
		return sprintfColor(NumberColor, strconv.FormatFloat(val, 'f', -1, 64))
	case bool:
		return sprintfColor(BoolColor, strconv.FormatBool(val))
	case nil:
		return sprintfColor(NullColor, "null")
	case map[string]interface{}:
		return processMap(val, depth)
	case []interface{}:
		return processArray(val, depth)
	}

	return ""
}

func processMap(m map[string]interface{}, depth int) string {
	currentIndent := generateIndent(depth - 1)
	nextIndent := generateIndent(depth)
	rows := []string{}
	keys := []string{}

	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		k := sprintfColor(KeyColor, `"%s"`, key)
		v := pretty(val, depth+1)
		row := fmt.Sprintf("%s%s: %s", nextIndent, k, v)
		rows = append(rows, row)
	}

	return fmt.Sprintf("{\n%s\n%s}", strings.Join(rows, ",\n"), currentIndent)
}

func processArray(a []interface{}, depth int) string {
	currentIndent := generateIndent(depth - 1)
	nextIndent := generateIndent(depth)
	rows := []string{}

	for _, val := range a {
		c := pretty(val, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(rows, ",\n"), currentIndent)
}

func generateIndent(depth int) string {
	return strings.Join(make([]string, Indent*depth+1), " ")
}
