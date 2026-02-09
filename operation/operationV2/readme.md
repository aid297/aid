### 多元操作

1. 三元操作：简单用法
   ```go
   package main
   
   import (
   	. `fmt`
   	`github.com/aid297/aid/operation/operationV2`
   )
   
   func main() {
   	res := operationV2.NewTernary(
   		operationV2.TrueValue("及格"),
   		operationV2.FalseValue("不及格"),
   	).GetByValue(true)
   
   	Printf("是否及格：%s", res)
   }
   ```

2. 三元操作：复杂用法。有些时候我们没有办法直接获取某一个内容的值，因为可能这个值在一个`Map`或者`Struct`中保存，而这个`Map`可能没有这个`KEY`或者这个`Struct`指针本身为空。如果使用`TrueValue`或``FalseValue`方法，在赋值时就需要直接获取到内容，这样会报`空指针错误`，于是我们就需要用到`TrueFn`或`FalseFn`方法。
   ```go
   package main
   
   import (
   	. `fmt`
   	`github.com/aid297/aid/operation/operationV2`
   )
   
   type Condition struct {
   	Score int
   }
   
   func main() {
   	condition := &Condition{Score: 85}
   
   	res := operationV2.NewTernary(
   		operationV2.TrueValue("及格"),
   		operationV2.FalseValue("不及格"),
   	).GetByValue(condition != nil && condition.Score >= 60)
   
   	Printf("是否及格：%s", res)
   
   	conditionMap := map[string]int{"分数": 70}
   
   	res2 := operationV2.NewTernary(
   		operationV2.TrueFn(func() int {
   			score, exist := conditionMap["分数"]
   			if !exist {
   				return 0
   			}
   			return score
   		}),
   		operationV2.FalseValue(0),
   	).GetByFunc(func() bool { return condition != nil })
   
   	Printf("是否及格2：%s\n", res2)
   }
   ```

3. 多元操作：多元操作通常场景在于启动程序时，我们需要确定程序的配置文件。如下所示，我们设计优先级为：终端参数 > 环境变量 > 默认参数。判断条件为：参数不为空。
   ```go
   package main
   
   import (
   	. `fmt`
   
   	`github.com/aid297/aid/operation/operationV2`
   )
   
   func main() {
   	m := operationV2.NewMultivariate[string]().
   		Append(operationV2.MultivariateAttr[string]{Item: "a", HitFunc: func(_ int, _ string) { Printf("采用高级") }}). // A 最高优先级：终端命令
   		Append(operationV2.MultivariateAttr[string]{Item: "b", HitFunc: func(idx int, item string) { Printf("采用次高级") }}). // B 次高优先级：全局变量
   		SetDefault(operationV2.MultivariateAttr[string]{Item: "c"}) // 设置默认值
   
   	_, f := m.Finally(func(item string) bool { return item != "" })
   
   	if f != "a" {
   		panic(Errorf("错误：%s", f))
   	}
   	Printf("成功：%s", f)
   }
   ```

   