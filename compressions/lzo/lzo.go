package lzo

import (
	"bytes"
	"io"

	"github.com/cyberdelia/lzo"

	"github.com/aid297/aid/v2/compressions"
)

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
		writer, err = lzo.NewWriterLevel(&buf, my.level)
	} else {
		writer, err = lzo.NewWriterLevel(&buf, lzo.BestCompression)
	}
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(my.data); err != nil {
		return nil, err
	}

	// ⚠️ 必须关闭以刷新缓冲区并写入尾部数据
	if err := writer.Close(); err != nil {
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
