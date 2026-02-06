### AnySlice 使用说明
1. 初始化

   ```go
   a1 := anySlice.New[int]()
   // 或
   a1 := anySlice.New[int](anySlice.Len(5)) // make([]int, 5)
   // 或
   a1 := anySlice.New[int](anySlice.Cap(5)) // make([]int, 0, 5)
   
   a2 := anySlice.NewList([]int{1, 2, 3})
   // 或
   a2 := anySlice.New[int]().SetAttr(anySlice.List([]int{1, 2, 3}))
   // 或
   a2 := anySlice.New[int](anySlice.List([]int{1, 2, 3}))
   
   a3 := anySlice.NewItems("a", "b", "c")
   // 或
   a3 := anySlice.New[string]().SetAttr(anySlice.Items("a", "b", "c"))
   // 或
   a3 := anySlice.New[string](anySlice.Items("a", "b", "c"))
   ```

2. 获取内容
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New[int](anySlice.Items(1, 2, 3))
   	Printf("%#v", a1.ToSlice()) // []int{1, 2, 3}
   }
   ```

3. 判断是否为空
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New[int](anySlice.Items(1, 2, 3))
   
   	Printf("is empty: %v", a1.Empty())        // false
   	Printf("is not empty: %v", a1.NotEmpty()) // true
   }
   ```

4. 如果为空则执行回调
   ```go
   package main
   
   import (
   	"errors"
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New[int](anySlice.Items(1, 2, 3))
   
   	a1.IfEmpty(func(array anySlice.AnySlicer[int]) {
   		Printf("array is empty: %v", array) // do nothing
   	})
   
   	a1.IfNotEmpty(func(array anySlice.AnySlicer[int]) {
   		Printf("array is not empty: %+v", array) // array is not empty: &{data:[1 2 3] mu:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:{_:{} v:0} readerWait:{_:{} v:0}}}
   	})
   
   	a1.IfEmptyError(func(array anySlice.AnySlicer[int]) error {
   		Printf("array is empty: %v", array) // do nothing
   		return errors.New("array is empty")
   	})
   
   	a1.IfNotEmptyError(func(array anySlice.AnySlicer[int]) error {
   		Printf("array is not empty: %+v", array) // array is not empty: &{data:[1 2 3] mu:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:{_:{} v:0} readerWait:{_:{} v:0}}}
   		return errors.New("array is not empty")
   	})
   }
   ```

5. 判断`key`是否存在
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New[int](anySlice.Items(1, 2, 3))
   
   	Printf("the key %d has? = %v", 2, a1.Has(2)) // the key 2 has? = true
   }
   ```

6. 通过`index`设置`value`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New(anySlice.Len[string](5))
   	a1.SetValue(0, "111").SetValue(1, "222").SetValue(4, "333")
   
   	Printf("%#v", a1.ToSlice()) // []string{"111", "222", "", "", "333"}
   }
   
   // 这里严格遵守golang中slice的用法
   // anySlice.Len=make([]string, size)
   // anySlice.Cap=make([]string, 0, size)
   // Len和Cap不能同时使用，后者会覆盖前者对内存重新分配
   ```

7. 通过`key`获取`value`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New(anySlice.Len[string](5))
   	a1.SetValue(4, "999")
   
   	Printf("number 4 is %s", a1.GetValue(4)) // number 4 is 999
   }
   ```

8. 通过`index`获取`value`的指针
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   	"github.com/aid297/aid/ptr"
   )
   
   func main() {
   	a1 := anySlice.New(anySlice.Len[string](5))
   	a1.SetValue(1, "888").SetValue(4, "999")
   
   	Printf("number 4 is %v\n", a1.GetValuePtr(4)) // number 4 is 0x1400013c130
   	// 等价于：
   	v := a1.GetValuePtr(4)
   	Printf("number 4 is %v\n", &v) // number 4 is 0x140000a2230。这里需要注意，两个指针地址不一样，因为 v 是 GetValuePtr(4) 返回值的拷贝
   	// 等价于
   	Printf("number 1 is %v", ptr.New(a1.GetValue(1))) // number 1 is 0x1400010b470
   }
   ```

9. 通过`index`获取`value`并指定默认值
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/array/anySlice"
   )
   
   func main() {
   	a1 := anySlice.New(anySlice.Cap[string](5))
   	Printf("%#v\n", a1.GetValueOrDefault(0, "default string"))  // "default string"
   }
   ```

10. 通过`index`获取多个`value`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Cap[string](5))
    	a1.Append("first", "second", "third")
    	Printf("%#v\n", a1.GetValues(0, 1, 2)) // []string{"first", "second", "third"}
    }
    ```

11. 通过`index`获取`切片`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Cap[string](5))
    	a1.Append("first", "second", "third")
    	Printf("%#v\n", a1.GetValuesBySlices(0, 3)) // []string{"first", "second", "third"} 等价于：a1[0:3]
    }
    ```

12. 获取`第一`和`最后`的值
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Cap[string](5))
    	a1.Append("first", "second", "third")
    	Printf("第一个值：%#v\n", a1.First())
    	Printf("最后一个值：%#v\n", a1.Last())
    }
    ```

13. 获取`原始切片`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Cap[string](5))
    	a1.Append("first", "second", "third")
    	Printf("a1 to slice: %#v\n", a1.ToSlice()) // a1 to slice: []string{"first", "second", "third"}
    }
    ```

14. 通过`value`获取`index`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Cap[string](5))
    	a1.Append("first", "second", "third")
    	Printf("second of index: %#v\n", a1.GetIndexByValue("second"))           // second of index: 1
    	Printf("third of index: %#v\n", a1.GetIndexesByValues("third", "first")) // third of index: []int{2, 0}
    }
    ```

15. 乱序
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[string]("first", "second", "third", "fourth", "fifth"))
    	Printf("Shuffle: %#v\n", a1.Shuffle().ToSlice()) // Shuffle: []string{"first", "fifth", "second", "fourth", "third"}
    }
    ```

16. 长度
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[string]("first", "second", "third", "fourth", "fifth"))
    	Printf("length of a1: %d\n", a1.Length()) // length of a1: 5
    }
    ```

17. 纯长度
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[string]("first", "second", "third", "fourth", "fifth", "", "", ""))
    	Printf("length of a1 for not empty: %d\n", a1.LengthNotEmpty()) // length of a1: 5
    }
    ```

18. 过滤
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](-2, -1, 0, 1, 2))
    	Printf("filter: %#v\n", a1.Filter(func(item int) bool { return item > 0 }).ToSlice()) // filter: []int{1, 2}
    }
    ```

19. 去掉`空值`和`零值`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    	Printf("remove empty: %#v\n", a1.RemoveEmpty().ToSlice()) // remove empty: []int{1, 3, 2}
    }
    ```

20. 拼接`字符串`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    	Printf("join to string: %#v\n", a1.Join("、")) // join to string: "0、1、0、3、0、2"
    }
    ```

21. 去掉`空值`和`零值`后拼接`字符串`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    	Printf("join to string: %#v\n", a1.JoinNotEmpty("、")) // join to string: "1、3、2"
    }
    ```

22. `in`和`not int`

    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    	Printf("in : %#v\n", a1.In(1, 2))        // true
    	Printf("not in : %#v\n", a1.NotIn(1, 2)) // false
    }
    ```

23. 如果`存在`或`不存在`则执行回调
    ```go
    package main
    
    import (
    	"errors"
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    
    	a1.IfIn(func(array anySlice.AnySlicer[int]) {
    		Println("1,2,3 is in this array")
    	}, 1, 2, 3) // 1,2,3 is in this array
    
    	a1.IfNotIn(func(array anySlice.AnySlicer[int]) {
    		Println("4,5,6 is not in this array")
    	}, 4, 5, 6) // 4,5,6 is not in this array
    
    	err := a1.IfNotInError(func(array anySlice.AnySlicer[int]) error {
    		return errors.New("7,8,9 is not in this array")
    	}, 7, 8, 9)
    	Printf("err: %#v\n", err) // err: &errors.errorString{s:"7,8,9 is not in this array"}
    }
    ```

24. 判断是否`全部为空`或`任意为空`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](0, 1, 0, 3, 0, 2))
    
    	Printf("is all empty: %#v\n", a1.AllEmpty()) // is all empty: false
    	Printf("is any empty: %#v\n", a1.AnyEmpty()) // is any empty: true
    }
    ```

25. 分块
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.New(anySlice.Items[int](1, 2, 3, 4, 5, 6, 7, 8, 9))
    
    	Printf("chunked: %#v\n", a1.Chunk(3)) // chunked: [][]int{[]int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9}}
    }
    ```

26. 摘取
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	type User struct {
    		Name string
    		Age  int
    	}
    	a1 := anySlice.New(anySlice.Items[User](
    		User{Name: "张三", Age: 17},
    		User{Name: "李四", Age: 18},
    		User{Name: "王五", Age: 27},
    		User{Name: "赵六", Age: 28},
    	))
    
    	Printf("pluck user's name: %#v\n", a1.Pluck(func(item User) any { return item.Age }).ToSlice()) // pluck user's name: []interface {}{17, 18, 27, 28}
    	// 注意#1：pluck 的结果是 []any，如果需要 []int 需要再转换一次
    	// 注意#2：pluck 的结果不会覆盖原有的a1，需要a2来接收
    }
    ```

27. 取`交集`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.NewItems(1, 2, 3, 4, 5)
    	a2 := anySlice.NewItems(2, 3, 4)
    
    	Printf("intersection: %#v\n", a1.Intersection(a2).ToSlice())             // intersection: []int{2, 3, 4}
    	Printf("intersection: %#v\n", a1.IntersectionBySlice(4, 5, 6).ToSlice()) // intersection: []int{4, 5}
    }
    ```

28. 取`差集`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.NewItems(1, 2, 3, 4, 5)
    	a2 := anySlice.NewItems(2, 3, 4)
    
    	Printf("difference: %#v\n", a1.Difference(a2).ToSlice())             // difference: []int{1, 5}
    	Printf("difference: %#v\n", a1.DifferenceBySlice(4, 5, 6).ToSlice()) // difference: []int{1, 2, 3}
    }
    ```

29. 取`并集`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.NewItems(1, 2, 3, 4, 5)
    	a2 := anySlice.NewItems(2, 3, 4)
    
    	Printf("union: %#v\n", a1.Union(a2).ToSlice())             // union: []int{1, 2, 3, 4, 5}
    	Printf("union: %#v\n", a1.UnionBySlice(4, 5, 6).ToSlice()) // union: []int{1, 2, 3, 4, 5, 6}
    }
    ```

30. 通过`index`删除
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.NewItems(1, 2, 3, 4, 5)
    	a1.RemoveByIndex(1, 2)
    	Printf("remove by index: %#v\n", a1.ToSlice()) // remove by index: []int{1, 4, 5}
    }
    ```

31. `every`和`each`迭代

    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    type User struct{ FirstName, LastName, FullName string }
    
    func main() {
    	a1 := anySlice.NewItems(
    		User{FirstName: "三", LastName: "张", FullName: "三张"},
    		User{FirstName: "四", LastName: "李", FullName: "李四"},
    		User{FirstName: "五", LastName: "王", FullName: "王五"},
    		User{FirstName: "六", LastName: "赵", FullName: "赵六"},
    	)
    
    	a1.
    		Every(func(item User) User {
    			item.FullName = item.FirstName + item.LastName
    			return item
    		}).
    		Each(func(idx int, item User) { Printf("idx: %d, item: %#v\n", idx, item) })
    	// idx: 0, item: main.User{FirstName:"三", LastName:"张", FullName:"三张"}
    	// idx: 1, item: main.User{FirstName:"四", LastName:"李", FullName:"李四"}
    	// idx: 2, item: main.User{FirstName:"五", LastName:"王", FullName:"王五"}
    	// idx: 3, item: main.User{FirstName:"六", LastName:"赵", FullName:"赵六"}
    }
    ```

32. 排序
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    )
    
    func main() {
    	a1 := anySlice.NewItems(1, 2, 3, 4, 5)
    
    	Println(a1.Sort(func(i, j int) bool { return a1.GetValue(i) > a1.GetValue(j) }).ToString())
    	// Output: [5 4 3 2 1]
    	// 等价于：sort.Slice(slice, func(i, j int) bool)
    }
    ```
    
33. 其他简单方法：
    ```go
    func Lock() AnySlicer[T] {}              // 写锁：加锁
    func Unlock() AnySlicer[T] {}            // 写锁：解锁
    func RLock() AnySlicer[T] {}             // 读锁：加锁
    func RUnlock() AnySlicer[T] {}           // 读锁：解锁
    func Clean() AnySlicer[T]                // 清理数据
    func MarshalJSON() ([]byte, error)      // JSON 序列化 → 实现 JSON 序列化接口
    func UnmarshalJSON(data []byte) error   // JSON 反序列化 → 实现 JSON 序列化接口
    func ToString(formats ...string) string // 转字符串
    ```

34. 高级方法：填充
    ```go
    a1 := anySlice.FillFunc[User,string]([]User{
    		{Name: "张三"},
    		{Name: "李四"},
    		{Name: "王五"},
    		{Name: "赵六"},
    	}, func(_ int, value User) string { return value.Name })
    Println(a1.ToString()) // [张三 李四 王五 赵六]
    ```

35. 高级方法：转换
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/array/anySlice"
    	"github.com/aid297/aid/operation/operationV2"
    )
    
    func main() {
    	a1 := anySlice.Cast[int, string](
    		anySlice.NewItems(40, 50, 60, 70, 80),
    		func(value int) string {
    			return operationV2.NewTernary(operationV2.TrueValue("及格"), operationV2.FalseValue("不及格")).GetByValue(value >= 60)
    		},
    	).ToString()
    	Println(a1) // [不及格 不及格 及格 及格 及格]
    }
    ```
    
    