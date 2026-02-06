package digestV2

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm3"
)

type (
	SM3Encoder interface{ Encode() (string, error) }

	SM3 struct{ original []byte }
)

func NewSM3(original string) SM3Encoder { return &SM3{original: []byte(original)} }

// Sm3 生成sm3摘要
func (my *SM3) Encode() (string, error) {
	h := sm3.New()
	if _, err := h.Write(my.original); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
