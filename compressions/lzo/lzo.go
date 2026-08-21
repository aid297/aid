package lzo

import (
	"bytes"
	"io"

	"github.com/cyberdelia/lzo"

	"github.com/aid297/aid/v3/compressions"
)

var _ compressions.Compressor = (*LZO)(nil)

type LZO struct {
	data     []byte
	level    int
	levelSet bool
}

func New(attrs ...compressions.CompressorAttr) (compressions.Compressor, error) {
	ins := &LZO{}
	err := ins.SetAttrs(attrs...)
	return ins, err
}

func (my *LZO) SetAttrs(attrs ...compressions.CompressorAttr) (err error) {
	for _, attr := range attrs {
		if err = attr(my); err != nil {
			return
		}
	}
	return
}

func (my *LZO) SetData(data []byte) { my.data = data }
func (my *LZO) SetLevel(level int)  { my.level = level; my.levelSet = true }

func (my *LZO) Encode() (compressed []byte, err error) {
	if len(my.data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer

	// 创建 lzo Writer，支持指定压缩级别
	var writer *lzo.Writer
	if my.levelSet {
		if writer, err = lzo.NewWriterLevel(&buf, my.level); err != nil {
			return nil, err
		}

		defer func() { _ = writer.Close() }()
	} else {
		if writer, err = lzo.NewWriterLevel(&buf, lzo.BestCompression); err != nil {
			return nil, err
		}

		defer func() { _ = writer.Close() }()
	}

	if _, err = writer.Write(my.data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (my *LZO) Decode() (decompressed []byte, err error) {
	var (
		reader    io.Reader
		lzoReader *lzo.Reader
	)

	if len(my.data) == 0 {
		return nil, nil
	}

	reader = bytes.NewReader(my.data)
	if lzoReader, err = lzo.NewReader(reader); err != nil {
		return nil, err
	}

	// 读取所有解压后的数据
	if decompressed, err = io.ReadAll(lzoReader); err != nil {
		return nil, err
	}

	return decompressed, nil
}
