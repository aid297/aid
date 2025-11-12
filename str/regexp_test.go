package str

import (
	"testing"
)

func Test1(t *testing.T) {
	a := "team1.abc"

	b := APP.Regexp.New("team1.", RegexpTargetString(a)).ReplaceAllString("")

	t.Logf("b: %v", b)
}
