package task

// PriorityQueue 实现了 heap.Interface 接口
type PriorityQueue []*Task

func (my PriorityQueue) Len() int { return len(my) }

// Less 决定了优先级的排序规则
// 这里我们定义 Priority 越大，优先级越高
func (my PriorityQueue) Less(i, j int) bool { return my[i].Priority > my[j].Priority }

func (my PriorityQueue) Swap(i, j int) {
	my[i], my[j] = my[j], my[i]
	my[i].Index = i
	my[j].Index = j
}

func (my *PriorityQueue) Push(x any) {
	task := x.(*Task)
	task.Index = len(*my)
	*my = append(*my, task)
}

func (my *PriorityQueue) Pop() any {
	old := *my
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	task.Index = -1 // 标记为已移除
	*my = old[:n-1]
	return task
}
