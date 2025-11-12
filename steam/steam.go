package steam

import (
	"bytes"
	"errors"
	"io"
)

type (
	Steam struct {
		Error      error
		readCloser io.ReadCloser
		copyFn     func(copied []byte) error
	}
)

// New 实例化：Steam
func (Steam) New(attrs ...Attributer) Steam {
	return Steam{}.Set(attrs...)
}

func (my Steam) Set(attrs ...Attributer) Steam {
	if len(attrs) > 0 {
		for _, attr := range attrs {
			attr.Register(&my)
		}
	}
	return my
}

// Copy 复制流
func (my Steam) Copy() (io.ReadCloser, error) {
	var (
		err    error
		copied = make([]byte, 0)
	)

	if my.readCloser == nil {
		return nil, errors.New("空内容")
	}

	if copied, err = io.ReadAll(my.readCloser); err != nil {
		return nil, err
	}

	if len(copied) > 0 {
		if err = my.copyFn(copied); err != nil {
			return nil, err
		}
	}

	return io.NopCloser(bytes.NewBuffer(copied)), err
}
