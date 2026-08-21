package texts_test

import (
	"testing"

	"github.com/aid297/aid/v3/texts"
	"github.com/aid297/aid/v3/texts/template"
)

func Test1(t *testing.T) {
	type A struct {
		Name string `template:"name"`
	}

	a := &A{Name: "A"}

	t.Logf("%v", texts.Template.New("11{{name}}22", template.Struct(a)).String())

}
