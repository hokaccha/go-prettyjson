package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

func main() {
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
	f := prettyjson.NewFormatter()
	f.Indent = 4
	f.KeyColor = color.New(color.FgMagenta)
	f.BoolColor = nil
	f.NullColor = color.New(color.Underline)
	s, _ := f.Marshal(v)
	fmt.Println(string(s))
}
