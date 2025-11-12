package digest

import (
	"github.com/tjfoc/gmsm/sm3"

	"encoding/hex"
)

// Sm3 生成sm3摘要
func Sm3(original []byte) string {
	h := sm3.New()
	if _, err := h.Write(original); err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}
