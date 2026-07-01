package texts_test

import (
	`testing`

	`github.com/aid297/aid/v2/texts`
	`github.com/aid297/aid/v2/texts/template`
)

func Test1(t *testing.T) {
	type A struct {
		Name string `template:"name"`
	}

	a := &A{Name: "A"}

	t.Logf("%v", texts.Templater.New("{{name}}", template.Struct(a)).String())

}
