package markdown

import `github.com/aid297/aid/v2/texts/buffer`

type (
	MarkdownWriter interface {
		New(options ...MarkdownAttr) MarkdownWriter
		GetBuffer() buffer.Buffer
		NewString(options ...MarkdownAttr) string
		Set(options ...MarkdownAttr) MarkdownWriter
		End() string
	}

	Markdown struct{ buffer buffer.Buffer }
)

func (Markdown) New(options ...MarkdownAttr) MarkdownWriter {
	return Markdown{buffer: buffer.BufferImpl{}.NewString("")}.Set(options...)
}

func (my Markdown) GetBuffer() buffer.Buffer { return my.buffer }

func (Markdown) NewString(options ...MarkdownAttr) string {
	return Markdown{}.New(options...).End()
}

func (my Markdown) Set(options ...MarkdownAttr) MarkdownWriter {
	if len(options) > 0 {
		for idx := range options {
			options[idx](&my)
		}
	}
	return my
}

func (my Markdown) End() string {
	return my.buffer.String()
}
