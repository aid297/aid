### 权重队列

```go
package main

import (
	`fmt`

	`github.com/aid297/aid/task`
)

func main() {
	// 1. 初始化队列
	tp := task.OnceTaskPool()
	// 2. 推送任务并执行
	tasks := tp.PushTask(
		&task.Task{ID: "task-low", Priority: 1},
		&task.Task{ID: "task-high", Priority: 10},
		&task.Task{ID: "task-medium", Priority: 5},
	).GO()
	for idx := range tasks {
		fmt.Printf("执行ID：%s (优先级：%d)\n", tasks[idx].ID, tasks[idx].Priority)
	}
}
```

