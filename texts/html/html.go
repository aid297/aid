package html

import `github.com/aid297/aid/v2/texts/buffer`

type (
	HTMLWriter interface {
		New(options ...HTMLAttr) HTMLWriter
		NewString(options ...HTMLAttr) string
		GetBuffer() buffer.Buffer
		Set(options ...HTMLAttr) HTMLWriter
		End() string
	}

	HTML struct{ buffer buffer.Buffer }
)

func (HTML) New(options ...HTMLAttr) HTMLWriter {
	return HTML{buffer: buffer.BufferImpl{}.NewString("")}.Set(options...)
}

func (HTML) NewString(options ...HTMLAttr) string {
	return HTML{}.New(options...).End()
}

func (my HTML) GetBuffer() buffer.Buffer { return my.buffer }

func (my HTML) Set(options ...HTMLAttr) HTMLWriter {
	if len(options) > 0 {
		for idx := range options {
			options[idx](&my)
		}
	}
	return my
}

func (my HTML) End() string {
	return my.buffer.String()
}
