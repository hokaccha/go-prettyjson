package prettyjson_test

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

func TestMarshalPretty(t *testing.T) {
	errFormat := "\nexpected:\n%s\n\nactual:\n%s"
	v := map[string]interface{}{
		"key": []interface{}{
			"str",
			100, 
			nil,
			true,
			false,
			map[string]string{"key":"str"},
		},
	}
	b, err := prettyjson.MarshalPretty(v)

	if err != nil {
		t.Error(err)
	}

	s := string(b)

	key := color.New(color.FgBlue, color.Bold).SprintFunc()(`"key"`)
	str := color.New(color.FgGreen, color.Bold).SprintFunc()(`"str"`)
	num := color.New(color.FgCyan, color.Bold).SprintFunc()("100")
	null := color.New(color.FgBlack, color.Bold).SprintFunc()("null")
	tru := color.New(color.FgYellow, color.Bold).SprintFunc()("true")
	fal := color.New(color.FgYellow, color.Bold).SprintFunc()("false")

	expected := fmt.Sprintf(`{
  %s: [
    %s,
    %s,
    %s,
    %s,
    %s,
    {
      %s: %s
    }
  ]
}`, key, str, num, null, tru, fal, key, str)

	if s != expected {
		t.Errorf(errFormat, expected, s)
	}
}
