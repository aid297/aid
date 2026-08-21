package buffer

import (
	"bytes"
	"net/url"
	"sync"
	"unicode/utf8"

	"github.com/spf13/cast"

	"github.com/aid297/aid/v3/digest"
)

type (
	Buffer interface {
		GetOriginal() *bytes.Buffer
		JoinString(strings ...string) string
		JoinStringLimit(limit string, strings ...string) string
		JoinAny(values ...any) string
		JoinAnyLimit(limit string, values ...any) string
		NewString(data ...string) Buffer
		NewBytes(original []byte) Buffer
		NewRune(data ...rune) Buffer
		NewAny(data ...any) Buffer
		Any(values ...any) Buffer
		URLPath(values ...string) Buffer
		URLQuery(values ...string) Buffer
		s(values ...string) Buffer
		S(values ...string) Buffer
		WhenS(condition bool, values ...string) Buffer
		b(values ...byte) Buffer
		B(values ...byte) Buffer
		WhenB(condition bool, values ...byte) Buffer
		r(values ...rune) Buffer
		R(values ...rune) Buffer
		WhenR(condition bool, values ...rune) Buffer
		String() string
		Bytes() []byte
		Ptr() *string
		Sha256Sum256() string
		Copy() Buffer
	}

	BufferImpl struct {
		original *bytes.Buffer
		lock     *sync.RWMutex
	}
)

func totalLenStrings(ss []string) int {
	total := 0
	for _, s := range ss {
		total += len(s)
	}
	return total
}

func totalLenAny(values []any) int {
	total := 0
	for _, v := range values {
		total += len(cast.ToString(v))
	}
	return total
}

func totalLenRunes(rs []rune) int {
	total := 0
	for _, r := range rs {
		total += utf8.RuneLen(r)
	}
	return total
}

func (my BufferImpl) GetOriginal() *bytes.Buffer { return my.original }

func (BufferImpl) JoinString(strings ...string) string {
	switch len(strings) {
	case 0:
		return ""
	case 1:
		return strings[0]
	default:
		return BufferImpl{}.NewString(strings[0]).S(strings[1:]...).String()
	}
}

func (BufferImpl) JoinStringLimit(limit string, strings ...string) string {
	switch len(strings) {
	case 0:
		return ""
	case 1:
		return strings[0]
	default:
		buf := BufferImpl{}.NewString(strings[0])
		remainingTotal := 0
		for _, s := range strings[1:] {
			remainingTotal += len(limit) + len(s)
		}
		buf.GetOriginal().Grow(remainingTotal)
		for idx := range strings[1:] {
			buf.S(limit).S(strings[idx+1])
		}
		return buf.String()
	}
}

func (BufferImpl) JoinAny(values ...any) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return cast.ToString(values[0])
	default:
		buf := BufferImpl{}.NewAny(values[0])
		buf.GetOriginal().Grow(totalLenAny(values))
		for idx := range values[1:] {
			buf.Any(values[idx+1])
		}
		return buf.String()
	}
}

func (BufferImpl) JoinAnyLimit(limit string, values ...any) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return cast.ToString(values[0])
	default:
		buf := BufferImpl{}.NewAny(values[0])
		extra := len(limit) * (len(values) - 1)
		buf.GetOriginal().Grow(totalLenAny(values) + extra)
		for idx := range values[1:] {
			buf.S(limit).Any(values[idx+1])
		}
		return buf.String()
	}
}

// NewString 实例化：通过字符串
func (BufferImpl) NewString(data ...string) Buffer {
	switch len(data) {
	case 0:
		return BufferImpl{original: bytes.NewBufferString(""), lock: &sync.RWMutex{}}
	case 1:
		return BufferImpl{original: bytes.NewBufferString(data[0]), lock: &sync.RWMutex{}}
	default:
		buf := BufferImpl{original: bytes.NewBufferString(data[0]), lock: &sync.RWMutex{}}
		buf.GetOriginal().Grow(totalLenStrings(data[1:]))
		for idx := range data[1:] {
			buf.GetOriginal().WriteString(data[idx+1])
		}
		return buf
	}
}

// NewBytes 实例化：通过字节码
func (BufferImpl) NewBytes(original []byte) Buffer {
	return BufferImpl{original: bytes.NewBuffer(original), lock: &sync.RWMutex{}}
}

// NewRune 实例化：通过字符
func (BufferImpl) NewRune(data ...rune) Buffer {
	switch len(data) {
	case 0:
		return BufferImpl{original: bytes.NewBufferString(""), lock: &sync.RWMutex{}}
	case 1:
		return BufferImpl{original: bytes.NewBufferString(string(data[0])), lock: &sync.RWMutex{}}
	default:
		buf := BufferImpl{original: bytes.NewBufferString(string(data)), lock: &sync.RWMutex{}}
		buf.GetOriginal().Grow(totalLenRunes(data[1:]))
		for idx := range data[1:] {
			buf.GetOriginal().WriteRune(data[idx+1])
		}
		return buf
	}
}

func (BufferImpl) NewAny(data ...any) Buffer {
	switch len(data) {
	case 0:
		return BufferImpl{original: bytes.NewBufferString(""), lock: &sync.RWMutex{}}
	case 1:
		return BufferImpl{original: bytes.NewBufferString(cast.ToString(data[0])), lock: &sync.RWMutex{}}
	default:
		buf := BufferImpl{original: bytes.NewBufferString(cast.ToString(data[0])), lock: &sync.RWMutex{}}
		buf.GetOriginal().Grow(totalLenAny(data[1:]))
		for idx := range data[1:] {
			buf.GetOriginal().WriteString(cast.ToString(data[idx+1]))
		}
		return buf
	}
}

// Any 追加任意内容到字符串
func (my BufferImpl) Any(values ...any) Buffer {
	for _, value := range values {
		my.original.Grow(totalLenAny(values))
		my.original.WriteString(cast.ToString(value))
	}

	return my
}

// URLPath 字符串 -> 追加写入 URL 路径
func (my BufferImpl) URLPath(values ...string) Buffer {
	total := 0
	escs := make([]string, 0, len(values))
	for _, v := range values {
		e := url.PathEscape(v)
		escs = append(escs, e)
		total += len(e)
	}
	my.original.Grow(total)
	// 然后使用 escs 而不是重复调用 PathEscape
	for _, e := range escs {
		my.original.WriteString(e)
	}

	return my
}

// URLQuery 字符串 -> 追加写入 URL 查询
func (my BufferImpl) URLQuery(values ...string) Buffer {
	total := 0
	escs := make([]string, 0, len(values))
	for _, v := range values {
		e := url.PathEscape(v)
		escs = append(escs, e)
		total += len(e)
	}
	my.original.Grow(total)
	// 然后使用 escs 而不是重复调用 PathEscape
	for _, e := range escs {
		my.original.WriteString(e)
	}

	return my
}

func (my BufferImpl) s(values ...string) Buffer {
	for _, value := range values {
		my.original.Grow(totalLenStrings(values))
		my.original.WriteString(value)
	}

	return my
}

// S 追加写入字符串
func (my BufferImpl) S(values ...string) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.s(values...)
}

func (my BufferImpl) WhenS(condition bool, values ...string) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	if condition {
		return my.s(values...)
	}
	return my
}

func (my BufferImpl) b(values ...byte) Buffer {
	for _, b := range values {
		my.original.Grow(len(values) - 1)
		my.original.WriteByte(b)
	}

	return my
}

// B 追加写入字节
func (my BufferImpl) B(values ...byte) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.b(values...)
}

func (my BufferImpl) WhenB(condition bool, values ...byte) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	if condition {
		return my.b(values...)
	}
	return my
}

func (my BufferImpl) r(values ...rune) Buffer {
	for _, v := range values {
		my.original.Grow(totalLenRunes(values))
		my.original.WriteRune(v)
	}

	return my
}

// R 追加写入字符
func (my BufferImpl) R(values ...rune) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.r(values...)
}

func (my BufferImpl) WhenR(condition bool, values ...rune) Buffer {
	my.lock.Lock()
	defer my.lock.Unlock()
	if condition {
		return my.r(values...)
	}
	return my
}

// String 获取字符串
func (my BufferImpl) String() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer my.original.Reset()
	return my.original.String()
}

// Bytes 获取字节码
func (my BufferImpl) Bytes() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer my.original.Reset()
	return my.original.Bytes()
}

// Ptr 获取字符串指针
func (my BufferImpl) Ptr() *string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer my.original.Reset()
	ret := my.original.String()
	return &ret
}

func (my BufferImpl) Sha256Sum256() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer my.original.Reset()
	return digest.Sha256Sum256(my.original.Bytes())
}

// Copy 复制当前对象
func (my BufferImpl) Copy() Buffer {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return BufferImpl{original: bytes.NewBuffer(my.original.Bytes()), lock: &sync.RWMutex{}}
}
