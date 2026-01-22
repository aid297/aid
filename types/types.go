package types

import "encoding/json"

// StructToOther 结构体通过 json 转其他格式
func StructToOther[K any, V any](params K) (ret V, err error) {
	var b []byte

	if b, err = json.Marshal(params); err != nil {
		return
	}

	if err = json.Unmarshal(b, &ret); err != nil {
		return
	}

	return
}
