package str

import "github.com/spf13/cast"

type (
	MarkdownAttributer interface {
		Register(markdown *Markdown)
	}

	AttrMarkdownNormal struct{ content string }
	AttrMarkdownA      struct{ name, href string }
	AttrMarkdownBr     struct{}
	AttrMarkdownUl     struct{ contents []string }
	AttrMarkdownOl     struct{ contents []string }
)

func MarkdownNormal(content string) AttrMarkdownNormal {
	return AttrMarkdownNormal{content: content}
}

func (my AttrMarkdownNormal) Register(markdown *Markdown) {
	markdown.buffer.S(my.content, "\n")
}

func MarkdownA(name, href string) AttrMarkdownA {
	return AttrMarkdownA{name: name, href: href}
}

func (my AttrMarkdownA) Register(markdown *Markdown) {
	markdown.buffer.S("[", my.name, "]", "(", my.href, ")\n")
}

func MarkdownBr() AttrMarkdownBr { return AttrMarkdownBr{} }

func (AttrMarkdownBr) Register(markdown *Markdown) {
	markdown.buffer.S("\n\n")
}

func MarkdownUl(contents ...string) AttrMarkdownUl {
	return AttrMarkdownUl{contents: contents}
}

func (my AttrMarkdownUl) Register(markdown *Markdown) {
	if len(my.contents) > 0 {
		for idx := range my.contents {
			markdown.buffer.S("* ", my.contents[idx], "\n")
		}
	}
}

func (my AttrMarkdownUl) Append(contents ...string) AttrMarkdownUl {
	my.contents = append(my.contents, contents...)
	return my
}

func MarkdownOl(contents ...string) AttrMarkdownOl { return AttrMarkdownOl{contents: contents} }

func (my AttrMarkdownOl) Register(markdown *Markdown) {
	if len(my.contents) > 0 {
		for idx := range my.contents {
			markdown.buffer.S(cast.ToString(idx+1), ". ", my.contents[idx], "\n")
		}
	}
}

func (my AttrMarkdownOl) Append(contents ...string) AttrMarkdownOl {
	my.contents = append(my.contents, contents...)
	return my
}
