package types

import (
	"github.com/bytedance/sonic"
)

// ToAny 转换到任意格式
func ToAny[DST any](src any, dst *DST) (err error) {
	var b []byte

	if b, err = sonic.Marshal(src); err != nil {
		return
	}

	return sonic.Unmarshal(b, &dst)
}
