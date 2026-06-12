package zstd4

import (
	"errors"

	"github.com/klauspost/compress/zstd"

	"github.com/aid297/aid/v2/compressions"
)

// Level 设置 zstd 压缩等级:
//
//	1 = SpeedFastest(最快速度)
//	2 = SpeedDefault(默认)
//	3 = SpeedBetterCompression(更好的压缩率)
//	4 = SpeedBestCompression(最高压缩率)
func Level(level int) compressions.CompressorAttr {
	return func(compressor compressions.Compressor) (err error) {
		switch zstd.EncoderLevel(level) {
		case zstd.SpeedFastest, zstd.SpeedDefault, zstd.SpeedBetterCompression, zstd.SpeedBestCompression:
		default:
			return errors.New("zstd 压缩等级错误，有效值为 1=SpeedFastest, 2=SpeedDefault, 3=SpeedBetterCompression, 4=SpeedBestCompression")
		}
		compressor.SetLevel(level)
		return nil
	}
}
