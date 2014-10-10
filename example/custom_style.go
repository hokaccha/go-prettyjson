package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

func main() {
	prettyjson.Indent = 4
	prettyjson.KeyColor = color.New(color.FgMagenta)
	prettyjson.BoolColor = nil
	prettyjson.NullColor = color.New(color.Underline)

	v := map[string]interface{}{
		"str":   "foo",
		"num":   100,
		"bool":  false,
		"null":  nil,
		"array": []string{"foo", "bar", "baz"},
		"map": map[string]interface{}{
			"foo": "bar",
		},
	}
	s, _ := prettyjson.MarshalPretty(v)
	fmt.Println(string(s))
}
