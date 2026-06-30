package secret

import "strings"

func PaddingBase64(encoded string) string {
	sub := ""
	if m := len(encoded) % 4; m != 0 {
		sub += strings.Repeat("=", 4-m)
	}

	return sub
}
