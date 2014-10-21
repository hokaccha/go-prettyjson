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
		"key": map[string]interface{}{
			"a": "str",
			"b": 100,
			"c": nil,
			"d": true,
			"e": false,
			"f": map[string]string{"key": "str"},
		},
	}
	b, err := prettyjson.MarshalPretty(v)

	if err != nil {
		t.Error(err)
	}

	blueBold := color.New(color.FgBlue, color.Bold).SprintFunc()
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	blackBold := color.New(color.FgBlack, color.Bold).SprintFunc()
	yelloBold := color.New(color.FgYellow, color.Bold).SprintFunc()

	format := `{
  %s: {
    %s: %s,
    %s: %s,
    %s: %s,
    %s: %s,
    %s: %s,
    %s: {
      %s: %s
    }
  }
}`

	expected := fmt.Sprintf(format,
		blueBold(`"key"`),
		blueBold(`"a"`), greenBold(`"str"`),
		blueBold(`"b"`), cyanBold("100"),
		blueBold(`"c"`), blackBold("null"),
		blueBold(`"d"`), yelloBold("true"),
		blueBold(`"e"`), yelloBold("false"),
		blueBold(`"f"`), blueBold(`"key"`), greenBold(`"str"`),
	)

	s := string(b)

	if s != expected {
		t.Errorf(errFormat, expected, s)
	}
}
