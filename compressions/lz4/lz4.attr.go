package lz4

import (
	"errors"

	"github.com/aid297/aid/v3/compressions"
)

// Level 设置 lz4 压缩等级: 0=Fast(默认), 1~9=Level1~Level9
// 等级越高压缩率越好，但速度越慢
func Level(level int) compressions.CompressorAttr {
	return func(compressor compressions.Compressor) (err error) {
		if level < 0 || level > 9 {
			return errors.New("lz4 压缩等级错误，有效范围为 0~9，其中 0 为 Fast 模式")
		}
		compressor.SetLevel(level)
		return nil
	}
}
