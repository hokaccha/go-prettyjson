package main

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func main() {
	s, _ := prettyjson.Format([]byte(`{"foo":"bar"}`))
	fmt.Println(string(s))
}
