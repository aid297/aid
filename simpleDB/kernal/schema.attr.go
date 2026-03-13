package kernal

import "github.com/aid297/aid/simpleDB/plugin"

type (
	SchemaAttributer interface{ RegisterAttr(db *SimpleDB) }

	AttrMaxMemoryBytes  struct{ MaxMemoryBytes uint64 }
	AttrMaxMemoryKB     struct{ MaxMemoryKB uint64 }
	AttrMaxMemoryMB     struct{ MaxMemoryMB uint64 }
	AttrMaxMemoryGB     struct{ MaxMemoryGB uint64 }
	AttrMaxCPUCores     struct{ MaxCPUCores uint8 }
	AttrUUIDVersion     struct{ UUIDVersion uint8 }
	AttrUUIDWithHyphen  struct{ UUIDWithHyphen *bool }
	AttrUUIDUpper       struct{ UUIDUpper *bool }
	AttrCascadeMaxDepth struct{ CascadeMaxDepth int }
	AttrPersistence     struct {
		WindowSeconds int
		WindowBytes   uint64
		Threshold     uint64
	}
	AttrSecurity struct {
		CompressAlgo string
		EncryptAlgo  string
		EncryptKey   string
	}
)

func MaxMemoryBytes(volume uint64) AttrMaxMemoryBytes { return AttrMaxMemoryBytes{volume} }
func (my AttrMaxMemoryBytes) RegisterAttr(db *SimpleDB) {
	if my.MaxMemoryBytes > 0 {
		db.config.MaxMemoryBytes = normalizeMemoryBytes(my.MaxMemoryBytes)
	}
}
func MaxMemoryKB(volume uint64) AttrMaxMemoryKB { return AttrMaxMemoryKB{volume} }
func (my AttrMaxMemoryKB) RegisterAttr(db *SimpleDB) {
	if my.MaxMemoryKB > 0 {
		db.config.MaxMemoryBytes = normalizeMemoryBytes(my.MaxMemoryKB * 1024)
	}
}
func MaxMemoryMB(volume uint64) AttrMaxMemoryMB { return AttrMaxMemoryMB{volume} }
func (my AttrMaxMemoryMB) RegisterAttr(db *SimpleDB) {
	if my.MaxMemoryMB > 0 {
		db.config.MaxMemoryBytes = normalizeMemoryBytes(my.MaxMemoryMB * 1024 * 1024)
	}
}
func MaxMemoryGB(volume uint64) AttrMaxMemoryGB { return AttrMaxMemoryGB{volume} }
func (my AttrMaxMemoryGB) RegisterAttr(db *SimpleDB) {
	if my.MaxMemoryGB > 0 {
		db.config.MaxMemoryBytes = normalizeMemoryBytes(my.MaxMemoryGB * 1024 * 1024 * 1024)
	}
}

func MaxCPUCores(cores uint8) AttrMaxCPUCores { return AttrMaxCPUCores{cores} }
func (my AttrMaxCPUCores) RegisterAttr(db *SimpleDB) {
	if my.MaxCPUCores > 0 {
		db.config.MaxCPUCores = normalizeCPUCores(int(my.MaxCPUCores))
	}
}

func UUIDVersion(version uint8) AttrUUIDVersion { return AttrUUIDVersion{version} }
func (my AttrUUIDVersion) RegisterAttr(db *SimpleDB) {
	if my.UUIDVersion >= 1 && my.UUIDVersion <= 8 {
		db.config.DefaultUUIDVersion = int(my.UUIDVersion)
	}
}

func UUIDWithHyphen(withHyphen bool) AttrUUIDWithHyphen { return AttrUUIDWithHyphen{&withHyphen} }
func (my AttrUUIDWithHyphen) RegisterAttr(db *SimpleDB) {
	if my.UUIDWithHyphen != nil {
		db.config.DefaultUUIDWithHyphen = my.UUIDWithHyphen
	}
}

func UUIDUpper(withUpper bool) AttrUUIDUpper { return AttrUUIDUpper{&withUpper} }
func (my AttrUUIDUpper) RegisterAttr(db *SimpleDB) {
	if my.UUIDUpper != nil {
		db.config.DefaultUUIDUppercase = my.UUIDUpper
	}
}

func CascadeMaxDepth(depth int) AttrCascadeMaxDepth { return AttrCascadeMaxDepth{depth} }
func (my AttrCascadeMaxDepth) RegisterAttr(db *SimpleDB) {
	if my.CascadeMaxDepth > 0 && my.CascadeMaxDepth <= HardCascadeMaxDepthLimit {
		db.config.DefaultCascadeMaxDepth = my.CascadeMaxDepth
	}
}

func Persistence(windowSeconds int, windowBytes uint64, threshold uint64) AttrPersistence {
	return AttrPersistence{
		WindowSeconds: windowSeconds,
		WindowBytes:   windowBytes,
		Threshold:     threshold,
	}
}
func (my AttrPersistence) RegisterAttr(db *SimpleDB) {
	if my.WindowSeconds > 0 {
		db.config.Persistence.WindowSeconds = my.WindowSeconds
	}
	if my.WindowBytes > 0 {
		db.config.Persistence.WindowBytes = my.WindowBytes
	}
	if my.Threshold > 0 {
		db.config.Persistence.Threshold = my.Threshold
	}
}

func Security(compressAlgo, encryptAlgo, encryptKey string) AttrSecurity {
	return AttrSecurity{compressAlgo, encryptAlgo, encryptKey}
}

func (my AttrSecurity) RegisterAttr(db *SimpleDB) {
	db.compressor = plugin.GetCompressor(my.CompressAlgo)
	db.encryptor = plugin.GetEncryptor(my.EncryptAlgo, my.EncryptKey)
}
