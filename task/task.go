package task

import (
	"container/heap"
	"fmt"
	"sync"
)

type (
	Task struct {
		ID       string
		Priority int          // 数值越大，优先级越高
		Index    int          // 用于 heap 内部维护
		Fn       func() error // 用于执行
		Error    error        // 执行结果
	}
	TaskPool struct {
		queue PriorityQueue
		mu    sync.Mutex
	}
)

var (
	taskPoolOnce sync.Once
	taskPoolIns  *TaskPool
)

func initTaskPool(tasks ...*Task) *TaskPool {
	tp := &TaskPool{queue: make(PriorityQueue, 0), mu: sync.Mutex{}}
	heap.Init(&tp.queue)
	tp.PushTask(tasks...)
	return tp
}

// OnceTaskPool 获取单例
func OnceTaskPool(tasks ...*Task) *TaskPool {
	taskPoolOnce.Do(func() { taskPoolIns = initTaskPool(tasks...) })
	return taskPoolIns
}

// PushTask 推送任务
func (my *TaskPool) PushTask(tasks ...*Task) *TaskPool {
	my.mu.Lock()
	defer my.mu.Unlock()

	if len(tasks) > 0 {
		for idx := range tasks {
			heap.Push(&my.queue, tasks[idx])
		}
	}

	return taskPoolIns
}

// GO 执行任务
func (my *TaskPool) GO() []*Task {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.queue.Len() > 0 {
		tasks := make([]*Task, 0, my.queue.Len())

		task := heap.Pop(&my.queue).(*Task)

		for my.queue.Len() > 0 {
			if task.Fn != nil {
				task.Error = task.Fn()
			}

			tasks = append(tasks, task)
		}
		return tasks
	}
	return nil
}

func Demo() {
	// 1. 初始化队列
	tp := OnceTaskPool()
	// 2. 推送任务并执行
	tasks := tp.PushTask(
		&Task{ID: "task-low", Priority: 1},
		&Task{ID: "task-high", Priority: 10},
		&Task{ID: "task-medium", Priority: 5},
	).GO()
	for idx := range tasks {
		fmt.Printf("执行ID：%s (优先级：%d)\n", tasks[idx].ID, tasks[idx].Priority)
	}
}
