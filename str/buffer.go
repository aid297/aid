package str

import (
	"bytes"
	"net/url"

	"github.com/aid297/aid/digest"
	"github.com/spf13/cast"
)

type Buffer struct {
	original *bytes.Buffer
}

var BufferApp Buffer

func (Buffer) JoinString(strings ...string) string {
	switch len(strings) {
	case 0:
		return ""
	case 1:
		return strings[0]
	default:
		return BufferApp.NewString(strings[0]).S(strings[1:]...).String()
	}
}

func (Buffer) JoinStringLimit(limit string, strings ...string) string {
	switch len(strings) {
	case 0:
		return ""
	case 1:
		return strings[0]
	default:
		buf := APP.Buffer.NewString(strings[0])
		for idx := range strings[1:] {
			buf.S(limit).S(strings[idx+1])
		}
		return buf.String()
	}
}

func (Buffer) JoinAny(values ...any) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return cast.ToString(values[0])
	default:
		buf := APP.Buffer.NewAny(values[0])
		for idx := range values[1:] {
			buf.Any(values[idx+1])
		}
		return buf.String()
	}
}

func (Buffer) JoinAnyLimit(limit string, values ...any) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return cast.ToString(values[0])
	default:
		buf := APP.Buffer.NewAny(values[0])
		for idx := range values[1:] {
			buf.S(limit).Any(values[idx+1])
		}
		return buf.String()
	}
}

// NewString 实例化：通过字符串
func (Buffer) NewString(data ...string) Buffer {
	switch len(data) {
	case 0:
		return Buffer{bytes.NewBufferString("")}
	case 1:
		return Buffer{bytes.NewBufferString(data[0])}
	default:
		buf := Buffer{bytes.NewBufferString(data[0])}
		for idx := range data[1:] {
			buf.original.WriteString(data[idx+1])
		}
		return buf
	}
}

// NewBytes 实例化：通过字节码
func (Buffer) NewBytes(original []byte) Buffer { return Buffer{bytes.NewBuffer(original)} }

// NewRune 实例化：通过字符
func (Buffer) NewRune(data ...rune) Buffer {
	switch len(data) {
	case 0:
		return Buffer{bytes.NewBufferString("")}
	case 1:
		return Buffer{bytes.NewBufferString(string(data[0]))}
	default:
		buf := Buffer{bytes.NewBufferString(string(data))}
		for idx := range data[1:] {
			buf.original.WriteRune(data[idx+1])
		}
		return buf
	}
}

func (Buffer) NewAny(data ...any) Buffer {
	switch len(data) {
	case 0:
		return Buffer{bytes.NewBufferString("")}
	case 1:
		return Buffer{bytes.NewBufferString(cast.ToString(data[0]))}
	default:
		buf := Buffer{bytes.NewBufferString(cast.ToString(data[0]))}
		for idx := range data[1:] {
			buf.original.WriteString(cast.ToString(data[idx+1]))
		}
		return buf
	}
}

// Any 追加任意内容到字符串
func (my Buffer) Any(values ...any) Buffer {
	for _, value := range values {
		my.original.WriteString(cast.ToString(value))
	}

	return my
}

// SUrlPath 字符串 -> 追加写入 URL 路径
func (my Buffer) SUrlPath(values ...string) Buffer {
	for _, value := range values {
		my.original.WriteString(url.PathEscape(value))
	}

	return my
}

// BUrlPath 字节 -> 追加写入 URL 路径
func (my Buffer) BUrlPath(values ...byte) Buffer {
	for _, value := range values {
		my.original.WriteString(url.PathEscape(string(value)))
	}

	return my
}

// RUrlPath 字符 -> 追加写入 URL 路径
func (my Buffer) RUrlPath(values ...rune) Buffer {
	for _, value := range values {
		my.original.WriteString(url.PathEscape(string(value)))
	}

	return my
}

// SUrlQuery 字符串 -> 追加写入 URL 查询
func (my Buffer) SUrlQuery(values ...string) Buffer {
	for _, value := range values {
		my.original.WriteString(url.QueryEscape(value))
	}

	return my
}

// BUrlQuery 字节 -> 追加写入 URL 查询
func (my Buffer) BUrlQuery(values ...byte) Buffer {
	for _, value := range values {
		my.original.WriteString(url.QueryEscape(string(value)))
	}

	return my
}

// RUrlQuery 字符 -> 追加写入 URL 查询
func (my Buffer) RUrlQuery(values ...rune) Buffer {
	for _, value := range values {
		my.original.WriteString(url.QueryEscape(string(value)))
	}

	return my
}

// S 追加写入字符串
func (my Buffer) S(values ...string) Buffer {
	for _, value := range values {
		my.original.WriteString(value)
	}

	return my
}

// B 追加写入字节
func (my Buffer) B(values ...byte) Buffer {
	for _, b := range values {
		my.original.WriteByte(b)
	}

	return my
}

// R 追加写入字符
func (my Buffer) R(values ...rune) Buffer {
	for _, v := range values {
		my.original.WriteRune(v)
	}

	return my
}

// String 获取字符串
func (my Buffer) String() string {
	defer my.original.Reset()
	return my.original.String()
}

// Bytes 获取字节码
func (my Buffer) Bytes() []byte {
	defer my.original.Reset()
	return my.original.Bytes()
}

// Ptr 获取字符串指针
func (my Buffer) Ptr() *string {
	defer my.original.Reset()
	ret := my.original.String()
	return &ret
}

func (my Buffer) Sha256Sum256() string {
	defer my.original.Reset()
	return digest.Sha256Sum256(my.original.Bytes())
}

// Copy 复制当前对象
func (my Buffer) Copy() Buffer { return Buffer{bytes.NewBuffer(my.original.Bytes())} }
