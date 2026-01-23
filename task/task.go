package task

import (
	"container/heap"
	"fmt"
	"sync"
)

type (
	Task struct {
		ID       string
		Priority int // 数值越大，优先级越高
		Index    int // 用于 heap 内部维护
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

func NewTaskPool(tasks ...*Task) *TaskPool {
	taskPoolOnce.Do(func() { taskPoolIns = initTaskPool(tasks...) })
	return taskPoolIns
}

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

func (my *TaskPool) GO() []*Task {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.queue.Len() > 0 {
		tasks := make([]*Task, 0, my.queue.Len())
		for my.queue.Len() > 0 {
			tasks = append(tasks, heap.Pop(&my.queue).(*Task))
		}
		return tasks
	}
	return nil
}

func Demo() {
	// 1. 初始化队列
	tp := NewTaskPool()
	// 2. 推送任务并执行
	tasks := tp.PushTask(
		&Task{ID: "task-low", Priority: 1},
		&Task{ID: "task-high", Priority: 10},
		&Task{ID: "task-medium", Priority: 5},
	).GO()
	for idx := range tasks {
		fmt.Printf("执行ID：%s (优先级：%d)\n", tasks[idx].ID, tasks[idx].Priority)
	}

	// // 1. 初始化队列
	// pq := make(PriorityQueue, 0)
	// heap.Init(&pq)

	// // 2. 投递任务（模拟不同优先级的任务）
	// heap.Push(&pq, &Task{ID: "task-low", Priority: 1})
	// heap.Push(&pq, &Task{ID: "task-high", Priority: 10})
	// heap.Push(&pq, &Task{ID: "task-medium", Priority: 5})

	// // 3. 消费任务
	// for pq.Len() > 0 {
	// 	task := heap.Pop(&pq).(*Task)
	// 	fmt.Printf("Processing task: %s (Priority: %d)\n", task.ID, task.Priority)
	// }
}
