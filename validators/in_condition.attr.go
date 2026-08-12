package validators

import (
	"fmt"

	"github.com/aid297/aid/v2/anySlices"
)

func EqualIn[SRC, DST any](data []SRC, fn func(val SRC) DST) string {
	return fmt.Sprintf("==%s", anySlices.FillFunc(data, func(_ int, value SRC) DST { return fn(value) }).JoinNotEmpty(defaultSliceSplitChar))
}

func NotEqualIn[SRC, DST any](data []SRC, fn func(val SRC) DST) string {
	return fmt.Sprintf("!=%s", anySlices.FillFunc(data, func(_ int, value SRC) DST { return fn(value) }).JoinNotEmpty(defaultSliceSplitChar))
}
