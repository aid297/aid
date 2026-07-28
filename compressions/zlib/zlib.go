package zlib

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/aid297/aid/v2/compressions"
)

var _ compressions.Compressor = (*Zlib)(nil)

type Zlib struct {
	data     []byte
	level    int
	levelSet bool
}

// New 实例化
func New(attrs ...compressions.CompressorAttr) (compressions.Compressor, error) {
	ins := &Zlib{}
	err := ins.SetAttrs(attrs...)
	return ins, err
}

func (my *Zlib) SetAttrs(attrs ...compressions.CompressorAttr) (err error) {
	for _, attr := range attrs {
		if err = attr(my); err != nil {
			return
		}
	}
	return
}

func (my *Zlib) SetData(data []byte) { my.data = data }
func (my *Zlib) SetLevel(level int)  { my.level = level; my.levelSet = true }

// Encode 压缩
func (my *Zlib) Encode() (compressed []byte, err error) {
	var (
		buffer bytes.Buffer
		writer *zlib.Writer
	)

	if len(my.data) == 0 {
		return nil, nil
	}

	// 创建 Zlib 压缩器，支持压缩等级:
	// -1=DefaultCompression(默认), 0=NoCompression, 1=BestSpeed ~ 9=BestCompression, -2=HuffmanOnly
	if my.levelSet {
		if writer, err = zlib.NewWriterLevel(&buffer, my.level); err != nil {
			return nil, err
		}

		defer func() { _ = writer.Close() }()
	} else {
		writer = zlib.NewWriter(&buffer)

		defer func() { _ = writer.Close() }()
	}

	// 写入数据到压缩器
	if _, err = writer.Write(my.data); err != nil {
		return nil, err
	}

	// 压缩后的数据存储在b的缓冲区中
	return buffer.Bytes(), nil
}

// Decode 解压缩
func (my *Zlib) Decode() (decompressed []byte, err error) {
	var (
		buffer bytes.Buffer
		reader io.ReadCloser
	)

	if len(my.data) == 0 {
		return nil, nil
	}

	if reader, err = zlib.NewReader(bytes.NewReader(my.data)); err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	// 读取解压缩后的数据到缓冲区
	if _, err = io.Copy(&buffer, reader); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
