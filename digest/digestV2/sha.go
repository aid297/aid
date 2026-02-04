package digestV2

import (
	"crypto/sha256"
	"encoding/hex"
)

type (
	SHAEncoder interface {
		Encode256() (string, error)
		Encode256Sum256() string
	}

	SHA struct{ original []byte }
)

func NewSHA(original string) SHAEncoder { return &SHA{original: []byte(original)} }

// Sha256 摘要算法
func (my *SHA) Encode256() (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(my.original); err != nil {
		return "", err
	}

	shaString := hex.EncodeToString(hash.Sum(nil))

	return shaString, nil
}

// Sha256Sum256 摘要算法
func (my *SHA) Encode256Sum256() string {
	h := sha256.Sum256(my.original)
	return hex.EncodeToString(h[:])
}
