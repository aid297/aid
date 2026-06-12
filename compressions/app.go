package compressions

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
