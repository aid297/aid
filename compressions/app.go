package compressions

import (
	`github.com/aid297/aid/v2/compressions/lz4`
	`github.com/aid297/aid/v2/compressions/lzo`
	`github.com/aid297/aid/v2/compressions/zlib`
	zstd4 `github.com/aid297/aid/v2/compressions/zstd4`
)

var (
	LZ4   Compressor = (*lz4.LZ4)(nil)
	LZO   Compressor = (*lzo.LZO)(nil)
	Zlib  Compressor = (*zlib.Zlib)(nil)
	ZSTD4 Compressor = (*zstd4.Zstd4)(nil)
)

type (
	Compressor interface {
		SetAttrs(attrs ...CompressorAttr) (err error)
		SetData(data []byte)
		SetLevel(level int)
		Encode() (compressed []byte, err error)
		Decode() (decompressed []byte, err error)
	}

	CompressorAttr func(compressor Compressor) (err error)
)

func Data(data []byte) CompressorAttr {
	return func(compressor Compressor) (err error) {
		compressor.SetData(data)
		return nil
	}
}
