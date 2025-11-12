package digest

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256 摘要算法
func Sha256(original []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(original); err != nil {
		return "", err
	}

	shaString := hex.EncodeToString(hash.Sum(nil))

	return shaString, nil
}

// Sha256Sum256 摘要算法
func Sha256Sum256(original []byte) string {
	hash := sha256.Sum256(original)
	return hex.EncodeToString(hash[:])
}
