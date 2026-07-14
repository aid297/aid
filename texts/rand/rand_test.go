package rand_test

import (
	"testing"

	"github.com/aid297/aid/v2/texts/rand"
)

func Test1(t *testing.T) {
	s := new(rand.RandomImpl).New().Strings(80)
	t.Logf("结果：%s", s)
}
