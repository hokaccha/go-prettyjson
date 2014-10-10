package main

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func main() {
	prettyjson.DisabledColor = true

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
