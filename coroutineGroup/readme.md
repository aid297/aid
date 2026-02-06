### 协程组

```go
	package main

	import (
		. "fmt"

		"github.com/aid297/aid/array/anySlice"
		"github.com/aid297/aid/coroutineGroup"
	)

	type (
		Page struct{ Page int }
		Res  struct{ Body string }
	)

	func main() {
		// 假设我们有一个分页组件，需要展示 1-9 的页的数据
		pagers := anySlice.NewItems(
			Page{Page: 1},
			Page{Page: 2},
			Page{Page: 3},
			Page{Page: 4},
			Page{Page: 5},
			Page{Page: 6},
			Page{Page: 7},
			Page{Page: 8},
			Page{Page: 9},
		)
		// 由于服务器压力限制，每次只能请求 3 页的数据
		capacities := 3
		pagersChunk := pagers.Chunk(capacities)

		cg := coroutineGroup.NewCoroutineGroup[Res]().
			SetBatchesByCapacities(pagers.Length(), capacities). // 设置最多 9 页，每次 3 页数据：相当于 3 批次
			GO(func(batch, capacity uint) (result coroutineGroup.Result[Res]) {
				result = coroutineGroup.Result[Res]{}

				if batch >= uint(len(pagersChunk)) { // 超出批次范围，标记为跳过
					result.IsSkip = true
					return
				}

				if capacity > uint(len(pagersChunk[batch])) { // 超出容量范围，标记为跳过
					result.IsSkip = true
					return
				}

				// 模拟请求数据
				result.Data = Res{Body: Sprintf("Batch %d, Capacity %d: Pages %v", batch+1, capacity, pagersChunk[batch][capacity].Page)}
				return
			})

		if !cg.OK {
			Printf("err: %v\n", cg.Error) // 这里是整个协程池错误
		}

		for _, res := range cg.Results {
			if res.IsSkip {
				continue
			}
			if res.Error != nil {
				Printf("partial err: %v\n", res.Error) // 这里是单次请求错误
				continue
			}
			Printf("result: %v\n", res.Data.Body) // 请求的单页结果
		}

		// 输出示例：
		// result: Batch 1, Capacity 2: Pages 3
		// result: Batch 1, Capacity 0: Pages 1
		// result: Batch 1, Capacity 1: Pages 2
		// result: Batch 2, Capacity 2: Pages 6
		// result: Batch 2, Capacity 1: Pages 5
		// result: Batch 2, Capacity 0: Pages 4
		// result: Batch 3, Capacity 2: Pages 9
		// result: Batch 3, Capacity 1: Pages 8
		// result: Batch 3, Capacity 0: Pages 7
    // 从输出结果中看到，分了三个批次，每一个批次三页请求。每个批次中的请求按照多协程进行。每个批次请求都结束后进行下一个批次的并发。
	}

```

