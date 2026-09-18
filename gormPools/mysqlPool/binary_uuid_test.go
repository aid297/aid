package mysqlPool_test

import (
	"testing"

	"github.com/aid297/aid/v2/gormPools/mysqlPool"
)

func Test1(t *testing.T) {
	a := "01a08e63-ed3f-7870-9e0f-d153e4ea99cc"
	// b := "01a08e22-dc6e-7abe-9c32-7642662cefd8"

	b, err := mysqlPool.BinaryFromString(a)
	if err != nil {
		t.Fatalf("%v", err)
	}

	c, err := mysqlPool.BinaryFromAny(b)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("%s", c.String())
}
