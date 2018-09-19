package main

import (
	"fmt"

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
	f.Indent = 0
	f.Newline = ""
	s, _ := f.Marshal(v)
	fmt.Println(string(s))
}
