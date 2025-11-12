package dict

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func Test1(t *testing.T) {
	t.Run("test1 NewAnyDict", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if fmt.Sprintf("%#v", d.ToMap()) != `map[string]int{"分数":18, "年龄":20}` {
			t.Fatal("错误")
		}
	})
}

func Test2(t *testing.T) {
	t.Run("test2 GetKeyByIndex方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetKeyByIndex(0) != "分数" {
			t.Fatal("错误")
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 GetKeysByIndexes方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetKeysByIndexes(0, 1).ToString() != `[分数 年龄]` {
			t.Fatal("错误")
		}
	})
}

func Test4(t *testing.T) {
	t.Run("test4 GetKeyByValue方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetKeyByValue(18) != "分数" {
			t.Fatal("错误")
		}
	})
}

func Test5(t *testing.T) {
	t.Run("test5 GetKeysByValues方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetKeysByValues(18, 100).ToString() != "[分数 年龄]" {
			t.Fatal("错误")
		}
	})
}

func Test6(t *testing.T) {
	t.Run("test6 GetValueByIndex方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetValueByIndex(0) != 18 {
			t.Fatal("错误")
		}
	})
}

func Test7(t *testing.T) {
	t.Run("test7 GetValuesByIndexes方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetValuesByIndexes(0, 1).ToString() != "[18 100]" {
			t.Fatal("错误")
		}
	})
}

func Test8(t *testing.T) {
	t.Run("test8 GetValueByKey方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetValueByKey("分数") != 18 {
			t.Fatal("错误")
		}
	})
}

func Test9(t *testing.T) {
	t.Run("test9 GetValuesByKeys方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetValuesByKeys("分数", "年龄").ToString() != "[18 100]" {
			t.Fatal("错误")
		}
	})
}

func Test10(t *testing.T) {
	t.Run("test10 GetIndexByKey方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetIndexByKey("分数") != 0 {
			t.Fatal("错误")
		}
	})
}

func Test11(t *testing.T) {
	t.Run("test11 GetIndexesByKeys方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetIndexesByKeys("分数", "年龄").ToString() != "[0 1]" {
			t.Fatal("错误")
		}
	})
}

func Test12(t *testing.T) {
	t.Run("test12 GetIndexByValue方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetIndexByValue(18) != 0 {
			t.Fatal("错误")
		}
	})
}

func Test13(t *testing.T) {
	t.Run("test13 GetIndexesByValues方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.GetIndexesByValues(18, 100).ToString() != "[0 1]" {
			t.Fatal("错误")
		}
	})
}

func Test14(t *testing.T) {
	t.Run("test14 Len方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.Len() != 2 {
			t.Fatal("错误")
		}
	})
}

func Test15(t *testing.T) {
	t.Run("test15 IsEmpty方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		if d.IsEmpty() {
			t.Fatal("错误")
		}
	})
}

func Test16(t *testing.T) {
	t.Run("test16 Copy方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		d2 := d.Copy()

		if !reflect.DeepEqual(d, d2) {
			t.Fatal("错误")
		}
	})
}

func Test17(t *testing.T) {
	t.Run("test17 GetKeys方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.GetKeys().ToString() != `[分数 年龄]` {
			t.Fatal("错误")
		}
	})
}

func Test18(t *testing.T) {
	t.Run("test18 GetValues方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.GetValues().ToString() != "[18 100]" {
			t.Fatal("错误")
		}
	})
}

func Test19(t *testing.T) {
	t.Run("test19 GetIndexes", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.getIndexes().ToString() != `[0 1]` {
			t.Fatal("错误")
		}
	})
}

func Test20(t *testing.T) {
	t.Run("test20 FirstKey方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.FirstKey() != "分数" {
			t.Fatal("错误")
		}
	})
}

func Test21(t *testing.T) {
	t.Run("test21 FirstValue方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.FirstValue() != 18 {
			t.Fatal("错误")
		}
	})
}

func Test22(t *testing.T) {
	t.Run("test22 LastKey方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.LastKey() != "年龄" {
			t.Fatal("错误")
		}
	})
}

func Test23(t *testing.T) {
	t.Run("test23 LastValue方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)
		if d.LastValue() != 100 {
			t.Fatal("错误")
		}
	})
}

func Test24(t *testing.T) {
	t.Run("test24 Filter方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		d.Filter(func(key string, value int) bool { return value > 18 })
		if d.GetValues().ToString() != "[100]" {
			t.Fatal("错误")
		}
	})
}

func Test25(t *testing.T) {
	t.Run("test25 RemoveByKey方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100)

		d.RemoveByKey("分数")
		if d.GetKeys().ToString() != "[年龄]" {
			t.Fatal("错误")
		}
	})
}

func Test26(t *testing.T) {
	t.Run("test26 RemoveByValue方法", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100).Set("税务", 0)
		d.RemoveByValue(0).RemoveByValue(100)

		if d.GetValues().ToString() != "[18]" {
			t.Fatal("错误")
		}
	})
}

func Test27(t *testing.T) {
	t.Run("test27 RemoveEmpty", func(t *testing.T) {
		d := Make[string, string]().Set("分数", "18").Set("年龄", "100").Set("税务", "")
		d.RemoveEmpty()

		if d.GetValues().ToString() != "[18 100]" {
			t.Fatal("错误")
		}
	})
}

func Test28(t *testing.T) {
	t.Run("test28 Join", func(t *testing.T) {
		d := Make[string, string]().Set("分数", "18").Set("年龄", "100").Set("税务", "")

		if d.Join(";") != `18;100;` {
			t.Fatal("错误")
		}
	})
}

func Test29(t *testing.T) {
	t.Run("test29 JoinWithoutEmpty", func(t *testing.T) {
		d := Make[string, string]().Set("分数", "18").Set("年龄", "100").Set("税务", "")

		if d.JoinWithoutEmpty(";") != `18;100` {
			t.Fatal("错误")
		}
	})
}

func Test30(t *testing.T) {
	t.Run("test30 InKeys", func(t *testing.T) {
		d := Make[string, string]().Set("分数", "18").Set("年龄", "100").Set("税务", "")
		if !d.InKeys("分数", "年龄") {
			t.Fatal("错误")
		}
	})
}

func Test31(t *testing.T) {
	t.Run("test31 InValues", func(t *testing.T) {
		d := Make[string, string]().Set("分数", "18").Set("年龄", "100").Set("税务", "")
		if !d.InValues("18", "100") {
			t.Fatal("错误")
		}
	})
}

func Test32(t *testing.T) {
	t.Run("test32 Every", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100).Set("税务", 0)
		d.Every(func(key string, value int) (string, int) { return key, value + 1 })

		if d.GetValues().ToString() != "[19 101 1]" {
			t.Fatal("错误")
		}
	})
}

func Test33(t *testing.T) {
	t.Run("test33 Each", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100).Set("税务", 0)
		d.Each(func(key string, value int) { fmt.Println(key, value) })
	})
}

func Test34(t *testing.T) {
	t.Run("test34 Clean", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100).Set("税务", 0)
		d.Clean()

		if d.GetValues().ToString() != `[]` {
			t.Fatal("错误")
		}
	})
}

func Test35(t *testing.T) {
	t.Run("test35 json序列化", func(t *testing.T) {
		d := Make[string, int]().Set("分数", 18).Set("年龄", 100).Set("税务", 0)
		b, e := json.Marshal(d)
		if e != nil {
			t.Fatalf("序列化错误：%v", e)
		}

		t.Logf("%s", b)
	})
}

func Test36(t *testing.T) {
	t.Run("test36 json反序列化", func(t *testing.T) {
		j := []byte(`[{"分数":18,"年龄":100,"税务":0}]`)
		var d []*AnyDict[string, int]

		e := json.Unmarshal(j, &d)
		if e != nil {
			t.Fatalf("反序列化错误：%v", e)
		}

		for _, a := range d {
			t.Logf("%+v\n", a)
		}
	})
}
