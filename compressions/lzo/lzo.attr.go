package lzo

import (
	"errors"

	"github.com/cyberdelia/lzo"

	"github.com/aid297/aid/v2/anySlices"
	"github.com/aid297/aid/v2/compressions"
)

func Level(level int) compressions.CompressorAttr {
	return func(compressor compressions.Compressor) (err error) {
		if anySlices.New(anySlices.Items(lzo.BestSpeed, lzo.BestCompression, -1)).NotIn(level) {
			return errors.New("压缩等级错误，请使用 lzo.BestSpeed 或 lzo.BestCompression")
		}

		compressor.SetLevel(level)
		return nil
	}
}
