package main

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/mattn/go-colorable"
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
	s, _ := prettyjson.Marshal(v)
	colorable.NewColorableStdout().Write(s)
	fmt.Print("\n")
}
