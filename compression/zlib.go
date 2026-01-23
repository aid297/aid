package compression

import (
	"bytes"
	"compress/zlib"
	"io"
)

type Zlib struct{}

// NewZlib 实例化Zlib
func NewZlib() *Zlib { return &Zlib{} }

// Compress 压缩
func (*Zlib) Compress(originalData []byte) ([]byte, error) {
	var (
		err    error
		buffer bytes.Buffer
		writer *zlib.Writer
	)

	// 创建一个新的Zlib压缩器
	writer = zlib.NewWriter(&buffer)

	// 写入数据到压缩器
	if _, err = writer.Write(originalData); err != nil {
		return nil, err
	}

	// 记住要关闭Writer以完成压缩
	if err = writer.Close(); err != nil {
		return nil, err
	}

	// 压缩后的数据存储在b的缓冲区中
	return buffer.Bytes(), nil
}

// Decompress 解压缩
func (*Zlib) Decompress(data []byte) ([]byte, error) {
	var (
		err    error
		buffer bytes.Buffer
		reader io.ReadCloser
	)
	if reader, err = zlib.NewReader(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	defer reader.Close()

	// 读取解压缩后的数据到缓冲区
	if _, err = io.Copy(&buffer, reader); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
