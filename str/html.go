package str

type HTML struct{ buffer Buffer }

func (HTML) New(options ...HtmlAttributer) HTML {
	return HTML{buffer: APP.Buffer.NewString("")}.Set(options...)
}

func (HTML) NewString(options ...HtmlAttributer) string {
	return APP.HTML.New(options...).End()
}

func (my HTML) Set(options ...HtmlAttributer) HTML {
	if len(options) > 0 {
		for _, option := range options {
			option.Register(&my)
		}
	}
	return my
}

func (my HTML) End() string {
	return my.buffer.String()
}
