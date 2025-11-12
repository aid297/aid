package str

type Markdown struct{ buffer Buffer }

func (Markdown) New(options ...MarkdownAttributer) Markdown {
	return Markdown{buffer: APP.Buffer.NewString("")}.Set(options...)
}

func (Markdown) NewString(options ...MarkdownAttributer) string {
	return APP.Markdown.New(options...).End()
}

func (my Markdown) Set(options ...MarkdownAttributer) Markdown {
	if len(options) > 0 {
		for idx := range options {
			options[idx].Register(&my)
		}
	}
	return my
}

func (my Markdown) End() string {
	return my.buffer.String()
}
