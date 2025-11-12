package steam

import "io"

type (
	Attributer interface {
		Register(steam *Steam)
	}

	AttrReadCloser struct{ readCloser io.ReadCloser }
	AttrCopyFn     struct{ fn func(copied []byte) error }
)

func ReadCloser(readCloser io.ReadCloser) AttrReadCloser {
	return AttrReadCloser{readCloser: readCloser}
}

func (my AttrReadCloser) Register(steam *Steam) { steam.readCloser = my.readCloser }

func CopyFn(fn func(copied []byte) error) AttrCopyFn { return AttrCopyFn{fn: fn} }

func (my AttrCopyFn) Register(steam *Steam) { steam.copyFn = my.fn }
