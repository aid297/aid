package zstd4

import (
	"github.com/klauspost/compress/zstd"

	"github.com/aid297/aid/v3/compressions"
)

var _ compressions.Compressor = (*Zstd4)(nil)

type Zstd4 struct {
	data     []byte
	level    zstd.EncoderLevel
	levelSet bool
}

func New(attrs ...compressions.CompressorAttr) (compressions.Compressor, error) {
	ins := &Zstd4{}
	err := ins.SetAttrs(attrs...)
	return ins, err
}

func (my *Zstd4) SetAttrs(attrs ...compressions.CompressorAttr) (err error) {
	for _, attr := range attrs {
		if err = attr(my); err != nil {
			return
		}
	}
	return
}

func (my *Zstd4) SetData(data []byte) { my.data = data }
func (my *Zstd4) SetLevel(level int) {
	my.level = zstd.EncoderLevel(level)
	my.levelSet = true
}

func (my *Zstd4) Encode() (compressed []byte, err error) {
	var encoder *zstd.Encoder

	if len(my.data) == 0 {
		return
	}

	// 创建压缩器，支持压缩等级:
	// 1=SpeedFastest, 2=SpeedDefault, 3=SpeedBetterCompression, 4=SpeedBestCompression
	opts := make([]zstd.EOption, 0)
	if my.levelSet {
		opts = append(opts, zstd.WithEncoderLevel(my.level))
	}
	if encoder, err = zstd.NewWriter(nil, opts...); err != nil {
		return nil, err
	}
	defer func() { _ = encoder.Close() }()

	// EncodeAll 直接返回压缩后的字节切片
	return encoder.EncodeAll(my.data, nil), nil
}

func (my *Zstd4) Decode() (decompressed []byte, err error) {
	var decoder *zstd.Decoder

	if len(my.data) == 0 {
		return
	}

	if decoder, err = zstd.NewReader(nil); err != nil {
		return nil, err
	}
	defer decoder.Close()

	return decoder.DecodeAll(my.data, nil)
}
