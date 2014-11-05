package main

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func main() {
	f := prettyjson.NewFormatter()
	s, _ := f.Format([]byte("{\"foo\":\"bar\"}"))
	fmt.Println(string(s))
}
