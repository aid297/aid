### 表单验证器

1. 概览
   ```go
   package main

   import (
   	"fmt"

   	"github.com/aid297/aid/v2/validators"
   )

   type UserRequest struct {
   	Firstname string
   	Lastname  *string
   	Age       int
   	IsActive  bool
   }

   func main() {
   	userRequest := UserRequest{
   		Firstname: "",
   		Lastname:  nil,
   		Age:       180,
   	}

   	v := validators.NewValidator(userRequest).Validate(
   		validators.NewChecker("Firstname", "姓").Min(">0"),
   		validators.NewChecker("Lastname", "名").Required().Min(">2"),
   		validators.NewChecker("Age", "年龄").Min(">0").Max("<100"),
   		validators.NewChecker("IsActive", "启用").True(),
   	)

   	fmt.Printf("验证是否通过：%v\n", v.Invalid())
   	for _, err := range v.GetErrors() {
   		fmt.Printf("%v\n", err)
   	}
   	// 验证是否通过：true
   	// 『姓』长度不能小于 0
   	// 『名』必填
   	// 『年龄』值不能大于 100
   	// 『启用』需要是 true
   }
   ```

2. 规则说明：string
   ```go
   v := validators.NewValidator(UserRequest{Firstname: "张"}).Validate(
   	validators.NewChecker("Firstname", "姓").Min(">0").Max("<64"),
   )
   ```

   * `string`类型支持验证规则：`Required`、`Min`、`Max`、`Size`、`In`、`Regex`、`Format`
   * `Min(">2")`：字符串`长度`必须`大于2`。`Max("<10")`同理，长度必须`小于10`。
   * `Size("==5")`：字符串长度必须`等于5`；`Size("!=5")` 必须`不等于5`；`Size("2~8")` 表示长度区间 `(2, 8)`。
   * `In("==a|b|c")`：字符串必须`在 a、b、c 之中`；`In("!=a|b|c")` 必须在`a、b、c 之外`。
   * `Regex("^[a-z]+$")`：字符串必须符合正则表达式。
   * `Format("DateOnly")`：字符串必须符合指定时间格式，支持：`RFC3339`、`RFC3339Nano`、`DateTime`、`DateOnly`、`TimeOnly`、`RFC822`、`RFC850`、`RFC1123`、`RFC1123Z`、`Kitchen`、`Stamp`、`StampMilli`、`StampMicro`、`StampNano`、`ANSIC`、`UnixDate`、`RubyDate`、`ReferenceLayout`、`SonarQubeDatetime`。

3. 规则说明：`int`、`int8`、`int16`、`int32`、`int64`、`uint`、`uint8`、`uint16`、`uint32`、`uint64`、`float32`、`float64`
   ```go
   v := validators.NewValidator(UserRequest{Age: 18, Level: 1, Score: 99.5}).Validate(
   	validators.NewChecker("Age", "年龄").Min(">0").Max("<100"),
   	validators.NewChecker("Level", "等级").Min(">0"),
   	validators.NewChecker("Score", "分数").Min(">60").Max("<100"),
   )
   ```

   * 数字类型支持验证规则：`Min`、`Max`、`Size`、`In`、`Regex`、`Format`
   * `Min(">0")`：数值必须`大于0`。`Max("<100")`同理。
   * `Size("==18")`：数值必须`等于18`；`Size("!=18")` 必须`不等于18`；区间写法 `Size("2~8")` 与 string 一致。
   * `In("==1|2|3")`：数值必须在`1、2、3 之中`。

4. 规则说明：`bool`
   ```go
   v := validators.NewValidator(UserRequest{IsActive: false}).Validate(
   	validators.NewChecker("IsActive", "启用").True(),
   )
   ```

   * `bool`类型支持验证规则：`Required`、`True`、`False`
   * `True()`：值必须为`true`。`False()`：值必须为`false`。
   * `True/False` 对其他基础类型同样生效（按布尔语义转换）：`int 1` 视为 `true`、`int 0` 视为 `false`、字符串 `"true"/"false"` 同理。
   * `Required`：`bool` 的零值 `false` 会被视为未填写而报错，需要 `true` 才能通过。

5. 规则说明：`slice`、`array`、`map`
   ```go
   v := validators.NewValidator(UserRequest{Tags: []string{"a"}}).Validate(
   	validators.NewChecker("Tags", "标签").Min(">1"),
   )
   ```

   * 集合类型支持验证规则：`Min`、`Max`、`Size`（按`长度`校验）
   * `Min(">1")`：元素个数必须`大于1`。`Max("<10")`、`Size("==3")` 同理。

6. 规则说明：`time.Time`
   ```go
   v := validators.NewValidator(UserRequest{Birthday: time.Now()}).Validate(
   	validators.NewChecker("Birthday", "生日").Format("DateOnly"),
   )
   ```

   * `time.Time`类型支持验证规则：`Format`
   * 校验时先按对应布局格式化时间，再匹配格式规则，支持的格式名与 string 的 `Format` 一致。

7. 指针字段与自定义基础类型
   ```go
   type (
   	AAA   string
   	MyInt int
   	MyBool bool

   	UserRequest struct {
   		Lastname *string  // 指针字段自动解引用后按 string 校验
   		Code     AAA      // 自定义类型按底层类型校验
   		Num      MyInt
   		Flag     MyBool
   	}
   )

   v := validators.NewValidator(UserRequest{
   	Lastname: ptr.New("张三"),
   	Code:     "ABC",
   	Num:      5,
   	Flag:     true,
   }).Validate(
   	validators.NewChecker("Lastname", "名").Required().Min(">2"),
   	validators.NewChecker("Code", "编码").Min(">2"),
   	validators.NewChecker("Num", "数字").Min(">3"),
   	validators.NewChecker("Flag", "标记").True(),
   )
   ```

   * 指针字段（如 `*string`）自动解引用后按底层类型校验；`nil` 指针在 `Required` 时按未填写报错。
   * 自定义基础类型（`type AAA string`、`type MyInt int`、`type MyBool bool` 等）按底层类型执行对应规则。

8. 校验器独立使用
   ```go
   // 不依赖 Validator，可对任意值直接校验
   c := validators.NewChecker("Name", "姓名").Required().Min(">2")
   c = c.Check("ab")
   if len(c.GetErrors()) > 0 {
   	fmt.Printf("%v\n", c.GetErrors())
   	// 『姓名』长度不能小于 2
   }
   ```

   * `NewChecker(field, name)`：`field` 必须与结构体字段名一致，`name` 用于错误提示。
   * 链式规则：`Required`、`Min`、`Max`、`Size`、`In`、`Regex`、`Format`、`True`、`False`，可任意组合。
   * `Check(original any)`：直接传入待校验的值，返回自身后可继续调用 `GetErrors()` 获取错误列表。
