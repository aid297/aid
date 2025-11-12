package array

import (
	"encoding/json"
	"reflect"
	"testing"
)

type A struct {
	Name string
}

func Test1(t *testing.T) {
	t.Run("test1 New", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.ToString() != "[1 2 3]" {
			t.Fatal("错误")
		}
	})
}

func Test2(t *testing.T) {
	t.Run("test2 IsEmpty", func(t *testing.T) {
		aa := New([]int{1, 2, 3})
		if aa.IsEmpty() {
			t.Fatal("错误")
		}

		aa2 := Make[int](0)
		if aa2.IsNotEmpty() {
			t.Fatalf("错误")
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 Has", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if !aa.Has(3) {
			t.Fatal("错误")
		}
	})
}

func Test4(t *testing.T) {
	t.Run("test4 Set", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Set(3, 100).ToString() != "[1 2 3 100]" {
			t.Fatal("错误")
		}
	})
}

func Test5(t *testing.T) {
	t.Run("test5 Get", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Get(2) != 3 {
			t.Fatal("错误")
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 GetByIndexes", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.GetByIndexes(0, 2).ToString() != "[1 3]" {
			t.Fatal("错误")
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 Append", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Append(4).ToString() != "[1 2 3 4]" {
			t.Fatal("错误")
		}
	})
}

func Test8(t *testing.T) {
	t.Run("test8 First", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.First() != 1 {
			t.Fatal("错误")
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test9 Last", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Last() != 3 {
			t.Fatal("错误")
		}
	})
}

func Test10(t *testing.T) {
	t.Run("test10 ToSlice", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if len(aa.ToSlice()) != 3 {
			t.Fatal("错误")
		}
	})
}

func Test11(t *testing.T) {
	t.Run("test11 GetIndexByValue", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.GetIndexByValue(2) != 1 {
			t.Fatal("错误")
		}
	})
}

func Test12(t *testing.T) {
	t.Run("test12 GetIndexesByValues", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.GetIndexesByValues(2, 3).ToString() != "[1 2]" {
			t.Fatal("错误")
		}
	})
}

func Test13(t *testing.T) {
	t.Run("test13 Copy", func(t *testing.T) {
		aa := New([]int{1, 2, 3})
		aa2 := aa.Copy()

		if !reflect.DeepEqual(aa, aa2) {
			t.Fatal("错误")
		}
	})
}

func Test14(t *testing.T) {
	t.Run("test14 Shuffle", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.ToString() == aa.Shuffle().ToString() {
			t.Fatal("错误")
		}
	})
}

func Test15(t *testing.T) {
	t.Run("test15 Len", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Len() != 3 {
			t.Fatal("错误")
		}
	})
}

func Test16(t *testing.T) {
	t.Run("test16 Filter", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		if aa.Filter(func(v int) bool { return v%2 == 0 }).ToString() != "[2 4 6 8 10]" {
			t.Fatal("错误")
		}
	})
}

func Test17(t *testing.T) {
	t.Run("test17 RemoveEmpty", func(t *testing.T) {
		aa := New([]string{"a", "", "c"})

		t.Log(aa.RemoveEmpty().ToString())

		if aa.RemoveEmpty().ToString() != "[a c]" {
			t.Fatal("错误")
		}
	})
}

func Test18(t *testing.T) {
	t.Run("test18 Join", func(t *testing.T) {
		aa := New([]string{"a", "", "c"})

		if aa.Join(";") != "a;;c" {
			t.Fatal("错误")
		}
	})
}

func Test19(t *testing.T) {
	t.Run("test19 JoinWithoutEmpty", func(t *testing.T) {
		aa := New([]string{"a", "", "c"})

		if aa.JoinWithoutEmpty(";") != "a;c" {
			t.Fatal("错误")
		}
	})
}

func Test20(t *testing.T) {
	t.Run("test20 In", func(t *testing.T) {
		aa := New([]string{"a", "", "c"})

		if !aa.In("a") {
			t.Fatal("错误")
		}
	})
}

func Test21(t *testing.T) {
	t.Run("test21 AllEmpty", func(t *testing.T) {
		aa := New([]string{"", "", ""})

		if !aa.AllEmpty() {
			t.Fatal("错误")
		}

		aa2 := New([]string{"", "a", ""})
		if aa2.AllEmpty() {
			t.Fatal("错误")
		}
	})
}

func Test22(t *testing.T) {
	t.Run("test22 AnyEmpty", func(t *testing.T) {
		aa := New([]string{"", "a", ""})

		if !aa.AnyEmpty() {
			t.Fatal("错误")
		}
	})
}

func Test23(t *testing.T) {
	t.Run("test23 Chunk", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		t.Logf("%v", aa.Chunk(3))
	})
}

func Test24(t *testing.T) {
	t.Run("test24 Pluck", func(t *testing.T) {
		aa := New([]A{{"a"}, {"b"}, {"c"}})
		t.Logf("%v", aa.Pluck(func(item A) any { return item.Name }))
	})
}

func Test25(t *testing.T) {
	t.Run("test25 Unique", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 1, 2, 3})

		if aa.Unique().ToString() != "[1 2 3]" {
			t.Fatal("错误")
		}
	})
}

func Test26(t *testing.T) {
	t.Run("test26 RemoveByIndexes", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		if aa.RemoveByIndexes(1, 2, 3).ToString() != "[1 5 6 7 8 9 10]" {
			t.Fatal("错误")
		}
	})
}

func Test27(t *testing.T) {
	t.Run("test27 RemoveByValues", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		if aa.RemoveByValues(1, 2, 3).ToString() != "[4 5 6 7 8 9 10]" {
			t.Fatal("错误")
		}
	})
}

func Test28(t *testing.T) {
	t.Run("test28 Every", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		if aa.Every(func(item int) int { return item * 2 }).ToString() != "[2 4 6]" {
			t.Fatal("错误")
		}
	})
}

func Test29(t *testing.T) {
	t.Run("test29 Each", func(t *testing.T) {
		aa := New([]string{"a", "b", "c", "d"})

		aa.Each(func(idx int, item string) {
			t.Log(idx, item)
		})
	})
}

func Test30(t *testing.T) {
	t.Run("test30 Clean", func(t *testing.T) {
		aa := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		if aa.Clean().ToString() != "[]" {
			t.Fatal("错误")
		}
	})
}

func Test31(t *testing.T) {
	t.Run("test31 json序列化", func(t *testing.T) {
		aa := New([]int{1, 2, 3})

		b, e := json.Marshal(aa)
		if e != nil {
			t.Fatalf("序列化错误：%v", e)
		}

		if string(b) != `[1,2,3]` {
			t.Fatal("错误")
		}
	})
}

func Test32(t *testing.T) {
	t.Run("test32 json反序列化", func(t *testing.T) {
		aa := Make[string](0)
		j := []byte(`["a","b","c"]`)

		if e := json.Unmarshal(j, &aa); e != nil {
			t.Fatalf("反序列化错误：%v", e)
		}

		if aa.ToString() != "[a b c]" {
			t.Fatal("错误")
		}
	})
}
