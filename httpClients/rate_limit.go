package httpClients

import (
	"io"
	"time"
)

type rateLimitReader struct {
	reader io.Reader
	rate   uint64
}

func NewUploadRateReader(reader io.Reader, rate uint64) io.Reader {
	return &rateLimitReader{reader: reader, rate: rate}
}

func (r *rateLimitReader) Read(p []byte) (int, error) {
	if r.rate == 0 {
		return r.reader.Read(p)
	}

	if len(p) > int(r.rate) {
		p = p[:int(r.rate)]
	}

	n, err := r.reader.Read(p)
	if n > 0 {
		sleep := time.Duration(int64(n)) * time.Second / time.Duration(r.rate)
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}

	return n, err
}
