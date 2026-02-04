package digestV2

import (
	"crypto/md5"
	"encoding/hex"
)

type (
	MD5HashEncoder interface{ Encode() (string, error) }
	MD5            struct{ original []byte }
)

// NewMD5 创建MD5编码器
func NewMD5(original string) MD5HashEncoder { return &MD5{original: []byte(original)} }

// Md5 编码
func (my *MD5) Encode() (string, error) {
	hash := md5.New()
	if _, err := hash.Write(my.original); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
