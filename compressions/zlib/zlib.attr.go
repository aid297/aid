package zlib

import (
	"compress/zlib"
	"errors"

	"github.com/aid297/aid/v2/compressions"
)

// Level 设置 zlib 压缩等级:
//
//	-1 = DefaultCompression(默认推荐)
//	 0 = NoCompression(无压缩)
//	 1 = BestSpeed(最快速度)
//	 9 = BestCompression(最高压缩率)
//	-2 = HuffmanOnly(仅 Huffman 编码)
func Level(level int) compressions.CompressorAttr {
	return func(compressor compressions.Compressor) (err error) {
		if level < zlib.HuffmanOnly || level > zlib.BestCompression {
			return errors.New("zlib 压缩等级错误，有效范围为 -2~9，其中 -1 为默认，0 为无压缩，1~9 逐级增强，-2 为 HuffmanOnly")
		}
		compressor.SetLevel(level)
		return nil
	}
}
