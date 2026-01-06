package coroutineGroup

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func Test1(t *testing.T) {
	cg := New[[]string]().
		SetBatches(3).
		SetCapacity(4).
		Run(func(batch, capacity uint) *Result[[]string] {
			return &Result[[]string]{
				Data:  []string{"轮数:", cast.ToString(batch), "次数:", cast.ToString(capacity)},
				Error: nil,
			}
		})

	for idx := range cg.Results {
		t.Logf("%v", cg.Results[idx])
	}
}

func Test2(t *testing.T) {
	cg := New[[]string]().
		SetBatches(3).
		SetCapacity(4).
		Run(func(batch, capacity uint) *Result[[]string] {
			return &Result[[]string]{
				Data:  []string{"轮数:", cast.ToString(batch), "次数:", cast.ToString(capacity)},
				Error: fmt.Errorf("模拟错误 %d:%d", batch, capacity),
			}
		})

	for idx := range cg.Results {
		t.Logf("%v", cg.Results[idx])
	}
}

func Test3(t *testing.T) {
	cg := New[[]string]().
		SetBatches(3).
		SetCapacity(4).
		RunUntilError(func(batch, capacity uint) *Result[[]string] {
			return &Result[[]string]{
				Data:  []string{"轮数:", cast.ToString(batch), "次数:", cast.ToString(capacity)},
				Error: fmt.Errorf("模拟错误 %d:%d", batch, capacity),
			}
		})

	for idx := range cg.Results {
		t.Logf("%v", cg.Results[idx])
	}
}
