package coroutineGroups

var (
	CoroutineGroup CoroutineGrouper[any] = (*CoroutineGroupImpl[any])(nil)
)
