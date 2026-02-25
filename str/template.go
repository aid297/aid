package str

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

type Template[T any] struct {
	err     error
	title   string
	content string
	s       T
	ret     []byte
	tmpl    *template.Template
}

// New 填充字符串
func NewTemplate[T any](title, content string, s T) *Template[T] {
	t := &Template[T]{title: title, content: content, s: s}

	var (
		err error
		buf bytes.Buffer
	)

	t.tmpl = template.Must(template.New(title).Parse(content))

	if err = t.tmpl.Execute(&buf, s); err != nil {
		t.err = fmt.Errorf("脚本填充失败：%w", err)
		return t
	}
	if t.ret, err = io.ReadAll(&buf); err != nil {
		t.err = fmt.Errorf("读取填充结果失败：%w", err)
		return t
	}

	return t
}

func (my *Template[T]) Error() error   { return my.err }
func (my *Template[T]) String() string { return string(my.ret) }
func (my *Template[T]) Bytes() []byte  { return my.ret }
