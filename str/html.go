package str

type (
	Html struct{ buffer Buffer }
)

func (Html) New(options ...HtmlAttributer) Html {
	return Html{buffer: APP.Buffer.NewString("")}.Set(options...)
}

func (Html) NewString(options ...HtmlAttributer) string {
	return APP.Html.New(options...).End()
}

func (my Html) Set(options ...HtmlAttributer) Html {
	if len(options) > 0 {
		for _, option := range options {
			option.Register(&my)
		}
	}
	return my
}

func (my Html) End() string {
	return my.buffer.String()
}
