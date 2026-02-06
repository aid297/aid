### AnyMap 使用说明

1. 初始化
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]string{"姓名": "张三", "年龄": "18"}))
   	Println(m1.ToString()) //map[姓名:张三 年龄:18]
   
   	m2 := anyMap.New(anyMap.Cap[string, int](10))
   	m2.SetAttrs(anyMap.Map(map[string]int{"张三": 18, "李四": 19, "王五": 20, "赵六": 21}))
   	Println(m2.ToString()) // map[张三:18 李四:19 王五:20 赵六:21]
   
   	m3 := anyMap.New(anyMap.Cap[string, bool](5))
   	m3.SetValue("张三", true).SetValue("李四", false).SetValue("王五", true).SetValue("赵六", false)
   	Println(m3.ToString()) // map[张三:true 李四:false 王五:true 赵六:false]
   }
   ```

2. 查看`key`是否存在
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"张三": 18, "李四": 19, "王五": 20, "赵六": 21}))
   	Printf("张三 exist: %v\n", m1.Has("张三")) // 张三 exist: true
   	Printf("孙七 exist: %v\n", m1.Has("孙七")) // 孙七 exist: false
   }
   ```

3. 获取`原始数据`
   ```go
   package main
   
    import (
        . "fmt"
   
        "github.com/aid297/aid/dict/anyMap"
    )
   
    func main() {
        m1 := anyMap.New(anyMap.Map(map[string]any{"张三": 18, "李四": 19, "王五": 20, "赵六": 21}))
        Printf("all : %v\n", m1.ToMap()) // all : map[张三:18 李四:19 王五:20 赵六:21]
    }
    // 原始数据的遍历的区别在于遍历时顺序不唯一
   ```

4. 判断是否为空`anyMapper`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	Println(anyMap.New(anyMap.Map(map[string]any{})).IsEmpty())          // true
   	Println(anyMap.New(anyMap.Map(map[string]any{"a": 1})).IsEmpty())    // false
   	Println(anyMap.New(anyMap.Map(map[string]any{"a": 1})).IsNotEmpty()) // true
   }
   ```

5. 通过`key`获取`value`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"a": 1, "b": 2, "c": 3}))
   	v, exist := m1.GetValueByKey("b")
   	Printf("b is exist: %v\n", exist) // b is exist: true
   	Printf("b's value is: %v\n", v)   // b's value is: 2
   
   	Printf("get many values by keys: %+v\n", m1.GetValuesByKeys("a", "c")) // get many values by keys: &{data:[1 3] mu:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:{_:{} v:0} readerWait:{_:{} v:0}}}
     
     // GetValuesByKeys方法返回的是一个anySlice.AnySlicer
   }
   ```

6. 判断`key`和`value`是否存在
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"a": 1, "b": 2, "c": 3}))
   	Printf("`a` is in keys of m1: %v\n", m1.HasKey("a"))           // `a` is in keys of m1: true
   	Printf("`a` or `d` in keys of m1: %v\n", m1.HasKeys("a", "d")) // `a` or `d` in keys of m1: true
   
   	Printf("`1` is in values of m1: %v\n", m1.HasValue(1))         // `1` is in values of m1: true
   	Printf("`1` or `4` in values of m1: %v\n", m1.HasValues(1, 4)) // `1` or `4` in values of m1: true
   }
   ```

7. `default`处理

   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"a": 1, "b": 2, "c": 3}))
   	m1.HasKeyDefault(
   		"d",
   		func(v any) any { return v.(int) + 111 }, // 如果存在 d，则回调 d 的值并处理再返回
   		func() any { return 4 },                  // 如果不存在 d，则设置 d 的值为4
   	)
   	Printf("%+v\n", m1.ToString()) // map[a:1 b:2 c:3 d:4]
   }
   ```

8. 获取`所有keys`和所有`values`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"a": 1, "b": 2, "c": 3}))
   	Printf("keys: %v\n", m1.GetKeys().ToSlice())     // keys: [c a b]
   	Printf("values: %v\n", m1.GetValues().ToSlice()) // values: [3 1 2]
     // 注意：从这个案例中可以看出，anyMap.Map初始化这种方式并不能保证 `Key` 的顺序
   }
   ```

9. 判断`长度`或`纯长度`
   ```go
   package main
   
   import (
   	. "fmt"
   
   	"github.com/aid297/aid/dict/anyMap"
   )
   
   func main() {
   	m1 := anyMap.New(anyMap.Map(map[string]any{"a": 1, "b": 2, "c": 3, "d": 0, "e": nil, "f": nil, "g": struct{}{}}))
   	Printf("length: %d\n", m1.Length())                        // length: 7
   	Printf("length with not empty: %d\n", m1.LengthNotEmpty()) // length with not empty: 3
   }
   ```

10. 过滤
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/dict/anyMap"
    )
    
    func main() {
    	m1 := anyMap.New(anyMap.Map(map[string]uint8{"张三": 40, "李四": 50, "王五": 60, "赵六": 70}))
    	m1.Filter(func(item uint8) bool { return item >= 60 })
    	Printf("filted map: %+v\n", m1.ToMap()) // filted map: map[王五:60 赵六:70] Filter方法会修改原map
    }
    ```

11. 剩余`简单方法`或与`AnySlicer`类似方法
    ```go
    RemoveEmpty() AnyMapper[K, V]
    Join(sep string) string
    JoinNotEmpty(sep string) string
    InKey(keys ...K) bool
    NotInKey(keys ...K) bool
    InValue(values ...V) bool
    NotInValue(values ...V) bool
    AllEmpty() bool
    AnyEmpty() bool
    RemoveByKey(key K) AnyMapper[K, V]
    RemoveByKeys(keys ...K) AnyMapper[K, V]
    RemoveByValue(value V) AnyMapper[K, V]
    RemoveByValues(values ...V) AnyMapper[K, V]
    Every(fn func(key K, value V) V) AnyMapper[K, V]
    Each(fn func(key K, value V)) AnyMapper[K, V]
    Clean() AnyMapper[K, V]
    MarshalJSON() ([]byte, error)
    UnmarshalJSON(data []byte) error
    ```

12. 高级方法：`转换`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/dict/anyMap"
    	"github.com/spf13/cast"
    )
    
    func main() {
    	m2 := anyMap.Cast(
    		anyMap.New(anyMap.Map(map[string]uint8{"张三": 40, "李四": 50, "王五": 60, "赵六": 70})),
    		func(k string, v uint8) string { return cast.ToString(v) },
    	)
    	Printf("m2: %#v\n", m2.ToMap()) // m2: map[string]string{"张三":"40", "李四":"50", "王五":"60", "赵六":"70"}
    }
    ```

13. 高级用法：`压缩`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/dict/anyMap"
    )
    
    func main() {
    	Println(
    		anyMap.Zip(
    			[]string{"张三", "李四", "王五", "赵六"},
    			[]bool{true, false, true, false},
    		).ToString(),
    	) // map[张三:true 李四:false 王五:true 赵六:false]
    }
    ```

14. 高级用法：`struct`转`map[string]any`
    ```go
    package main
    
    import (
    	. "fmt"
    
    	"github.com/aid297/aid/dict/anyMap"
    )
    
    type User struct {
    	Name   string `json:"name"`
    	Age    int    `json:"age"`
    	Gender string `json:"gender"`
    }
    
    func main() {
    	u := User{Name: "Alice", Age: 30, Gender: "female"}
    	r, e := anyMap.StructToOther[User, map[string]any](u)
    	if e != nil {
    		panic(e)
    	}
    	Printf("%v\n", r) // map[age:30 gender:female name:Alice]
    }
    ```

    