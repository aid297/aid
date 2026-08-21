package html

import "github.com/aid297/aid/v3/texts/buffer"

type (
	HTMLWriter interface {
		New(options ...HTMLAttr) HTMLWriter
		NewString(options ...HTMLAttr) string
		GetBuffer() buffer.Buffer
		Set(options ...HTMLAttr) HTMLWriter
		End() string
	}

	HTMLWriterImpl struct{ buffer buffer.Buffer }
)

func (HTMLWriterImpl) New(options ...HTMLAttr) HTMLWriter {
	return HTMLWriterImpl{buffer: buffer.BufferImpl{}.NewString("")}.Set(options...)
}

func (HTMLWriterImpl) NewString(options ...HTMLAttr) string {
	return HTMLWriterImpl{}.New(options...).End()
}

func (my HTMLWriterImpl) GetBuffer() buffer.Buffer { return my.buffer }

func (my HTMLWriterImpl) Set(options ...HTMLAttr) HTMLWriter {
	if len(options) > 0 {
		for idx := range options {
			options[idx](&my)
		}
	}
	return my
}

func (my HTMLWriterImpl) End() string {
	return my.buffer.String()
}
