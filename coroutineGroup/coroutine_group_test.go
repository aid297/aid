package coroutineGroup

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func Test1(t *testing.T) {
	cg := NewCoroutineGroup[[]string]().
		SetBatches(3).
		SetCapacity(4).
		GO(func(batch, capacity uint) Result[[]string] {
			return Result[[]string]{
				Data:  []string{"轮数:", cast.ToString(batch), "次数:", cast.ToString(capacity)},
				Error: nil,
			}
		})

	for idx := range cg.Results {
		t.Logf("%v", cg.Results[idx])
	}
}

func Test2(t *testing.T) {
	cg := NewCoroutineGroup[[]string]().
		SetBatches(3).
		SetCapacity(4).
		GO(func(batch, capacity uint) (result Result[[]string]) {
			result = Result[[]string]{}

			result.Data = []string{"轮数:", cast.ToString(batch), "次数:", cast.ToString(capacity)}
			result.Error = fmt.Errorf("模拟错误 %d:%d", batch, capacity)
			result.IsSkip = false
			return
		})

	for idx := range cg.Results {
		t.Logf("%v", cg.Results[idx])
	}
}
