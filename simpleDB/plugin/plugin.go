package plugin

// Compressor defines the interface for data compression
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

// Encryptor defines the interface for data encryption
type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

// NoOpCompressor does nothing
type NoOpCompressor struct{}

func (n *NoOpCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (n *NoOpCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

// NoOpEncryptor does nothing
type NoOpEncryptor struct{}

func (n *NoOpEncryptor) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (n *NoOpEncryptor) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

// Registry
var (
	compressors = make(map[string]func() Compressor)
	encryptors  = make(map[string]func(key string) Encryptor)
)

func RegisterCompressor(name string, factory func() Compressor) {
	compressors[name] = factory
}

func GetCompressor(name string) Compressor {
	if factory, ok := compressors[name]; ok {
		return factory()
	}
	return &NoOpCompressor{}
}

func RegisterEncryptor(name string, factory func(key string) Encryptor) {
	encryptors[name] = factory
}

func GetEncryptor(name string, key string) Encryptor {
	if factory, ok := encryptors[name]; ok {
		return factory(key)
	}
	return &NoOpEncryptor{}
}
