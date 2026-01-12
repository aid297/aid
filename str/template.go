package str

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

type Template[T any] struct {
}

// New 填充字符串
func NewTemplate[T any](title, content string, s T) (string, error) {
	var (
		err error
		buf bytes.Buffer
		ret []byte
	)

	tmpl := template.Must(template.New(title).Parse(content))

	if err = tmpl.Execute(&buf, s); err != nil {
		return "", fmt.Errorf("脚本填充失败：%w", err)
	}
	if ret, err = io.ReadAll(&buf); err != nil {
		return "", fmt.Errorf("生成脚本失败：%w", err)
	}

	return string(ret), nil
}
