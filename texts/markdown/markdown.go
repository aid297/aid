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

	MarkdownWriterImpl struct{ buffer buffer.Buffer }
)

func (MarkdownWriterImpl) New(options ...MarkdownAttr) MarkdownWriter {
	return MarkdownWriterImpl{buffer: buffer.BufferImpl{}.NewString("")}.Set(options...)
}

func (my MarkdownWriterImpl) GetBuffer() buffer.Buffer { return my.buffer }

func (MarkdownWriterImpl) NewString(options ...MarkdownAttr) string {
	return MarkdownWriterImpl{}.New(options...).End()
}

func (my MarkdownWriterImpl) Set(options ...MarkdownAttr) MarkdownWriter {
	if len(options) > 0 {
		for idx := range options {
			options[idx](&my)
		}
	}
	return my
}

func (my MarkdownWriterImpl) End() string {
	return my.buffer.String()
}
