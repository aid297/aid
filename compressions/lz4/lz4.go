package lz4

import (
	"bytes"
	"io"

	lz4 "github.com/pierrec/lz4/v4"

	"github.com/aid297/aid/v2/compressions"
)

var _ compressions.Compressor = (*LZ4)(nil)

type LZ4 struct {
	data  []byte
	level int
}

func New(attrs ...compressions.CompressorAttr) (compressions.Compressor, error) {
	ins := &LZ4{}
	err := ins.SetAttrs(attrs...)
	return ins, err
}

func (my *LZ4) SetAttrs(attrs ...compressions.CompressorAttr) (err error) {
	for _, attr := range attrs {
		if err = attr(my); err != nil {
			return
		}
	}
	return
}

func (my *LZ4) SetData(data []byte) { my.data = data }

func (my *LZ4) SetLevel(level int) { my.level = level }

func (my *LZ4) Encode() (compressed []byte, err error) {
	var buf bytes.Buffer

	if len(my.data) == 0 {
		return
	}

	w := lz4.NewWriter(&buf)
	defer func() { _ = w.Close() }()

	// 设置压缩等级: 0=Fast(默认), 1~9=Level1~Level9
	if my.level > 0 && my.level <= 9 {
		if err = w.Apply(lz4.CompressionLevelOption(lz4.CompressionLevel(1 << (8 + my.level)))); err != nil {
			return nil, err
		}
	}

	if _, err = w.Write(my.data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (my *LZ4) Decode() (decompressed []byte, err error) {
	if len(my.data) == 0 {
		return
	}

	return io.ReadAll(lz4.NewReader(bytes.NewReader(my.data)))
}
