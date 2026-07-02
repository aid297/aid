package coroutineGroups

import `github.com/aid297/aid/v2/coroutineGroups_delete/coroutineGroupV2`

var (
	CoroutineGroupV1 CoroutineGroup[any]                  = (*CoroutineGroupImpl[any])(nil)
	CoroutineGroupV2 coroutineGroupV2.CoroutineGroup[any] = (*coroutineGroupV2.CoroutineGroupImpl[any])(nil)
)
