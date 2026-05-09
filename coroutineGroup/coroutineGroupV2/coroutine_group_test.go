package coroutineGroupV2

import "testing"

func Test(t *testing.T) {
	var results []Result[int]
	results = New[int](4).
		GO(
			func() Result[int] { return Result[int]{Value: 1} },
			func() Result[int] { return Result[int]{Value: 2} },
			func() Result[int] { return Result[int]{Value: 3} },
			func() Result[int] { return Result[int]{Value: 4} },
			func() Result[int] { return Result[int]{Value: 6} },
			func() Result[int] { return Result[int]{Value: 7} },
			func() Result[int] { return Result[int]{Value: 8} },
			func() Result[int] { return Result[int]{Value: 9} },
			func() Result[int] { return Result[int]{Value: 10} },
		)

	t.Logf("results: %+v", results)
}
