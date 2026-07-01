package markdown

import "github.com/spf13/cast"

type MarkdownAttr func(m MarkdownWriter)

func Normal(content string) MarkdownAttr {
	return func(m MarkdownWriter) { m.GetBuffer().S(content) }
}

func A(name, href string) MarkdownAttr {
	return func(m MarkdownWriter) { m.GetBuffer().S("[", name, "]", "(", href, ")\n") }
}

func Br() MarkdownAttr { return func(m MarkdownWriter) { m.GetBuffer().S("\n\n") } }

func Ul(contents ...string) MarkdownAttr {
	return func(m MarkdownWriter) {
		if len(contents) > 0 {
			for idx := range contents {
				m.GetBuffer().S("* ", contents[idx], "\n")
			}
		}
	}
}

func Ol(contents ...string) MarkdownAttr {
	return func(m MarkdownWriter) {
		if len(contents) > 0 {
			for idx := range contents {
				m.GetBuffer().S(cast.ToString(idx+1), ". ", contents[idx], "\n")
			}
		}
	}
}
