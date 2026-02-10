### 表单验证器

1. 概览
   ```go
   package main
   
   import (
   	`fmt`
   	`time`
   
   	`github.com/aid297/aid/ptr`
   	`github.com/aid297/aid/validator/validatorV3`
   )
   
   type UserRequest struct {
   	Firstname string    `json:"firstname" v-rule:"(min>0)" v-name:"姓"`
   	Lastname  *string   `json:"lastname" v-rule:"(required)(min>2)" v-name:"名"`
   	Age       int       `json:"int" v-rule:"(min>=18)(max<100)" v-name:"年龄"`
   	Score     int       `json:"score" v-rule:"(required)" v-name:"分数"`
   	Birthday  time.Time `json:"birthday" v-rule:"(required)" v-name:"生日"`
   }
   
   func main() {
   	userRequest := UserRequest{
   		Firstname: "",
   		Lastname:  ptr.New("张"),
   		Age:       180,
   	}
   
   	checker := validatorV3.APP.Validator.Once().Checker(&userRequest)
   
   	checker.Validate()
   	fmt.Printf("验证结果：%v\n", checker.OK())
   
   	for _, wrong := range checker.Wrongs() {
   		fmt.Printf("%v\n", wrong)
   	}
   	// 验证结果：false
   	// [姓] 长度错误 期望：> 0
   	// [名] 长度错误 期望：> 2
   	// [年龄] 长度错误 期望：< 100
   	// [生日] 不能为空
   }
   ```

2. 规则说明：string
   ```go
   type UserRequest struct {
   	Firstname string    `json:"firstname" v-rule:"(required)(min>0)(max<64)" v-name:"姓"`
   }
   ```

   * `string`类型支持验证规则：`required`、`min>`、`min>=`、`max<`、`max<=`、`in:a,b,c`、`not-in:a,b,c`、`size=`、`size!=`、`ex:fn-a,fn-b,fn-c`
   * `required`：当前字符串`长度`不能等于`0`。
   * `min>`：当前字符串长度不能`小于x`。`max<`同理。
   * `min>=`：当前字符串长度必须`大于x`。`max<=`同理。
   * `in`、`not-in` ：当前字符必须`在x范围`内或`在x范围`外。
   * `size=`、`size!=`：当前字符串长度必须`等于x`或必须`不等于x`
   * `ex`：额外执行的验证程序，需要提前在验证器中注册。

3. 规则说明：*string
   ```go
   type UserRequest struct {
   	Firstname *string `json:"firstname" v-rule:"(required)(min>0)(max<64)" v-name:"姓"`
   }
   ```

   * `*string`类型支持验证规则：与`string`一致
   * `required`：当前字符串不能为`nil`且`长度`不能等于`0`

4. 规则说明：`int`、`int8`、`int16`、`int32`、`int64`
   ```go
   type UserRequest struct {
   	Age int `v-rule:"(min>0)(max<100)" v-name:"年龄"`
   }
   ```

   * `int`类型支持验证规则：`min`>、`min>=`、`max<`、`max<=`、`in`、`not-in`、`size=`、`size!=`、`ex:fn-a,fn-b,fn-c`

5. 规则说明：`*int`、`*int8`、`*int16`、`*int32`、`*int64`
   ```go
   type UserRequest struct {
   	Age *int `v-rule:"(min>0)(max<100)" v-name:"年龄"`
   }
   ```

   * `*int`类型支持验证规则：`required`、`min`>、`min>=`、`max<`、`max<=`、`in`、`not-in`、`size=`、`size!=`、`ex:fn-a,fn-b,fn-c`
6. 其他普通类型支持：`uint`、`uint8`、`uint16`、`uint32`、`uint64`、`*uint`、`*uint8`、`*uint16`、`*uint32`、`*uint64`、`float32`、`float64`、`*float32`、`*float64`
7. 特殊类型支持：`time.Time`、`*time.Time`

   * 规则支持：`required`、`min`>、`min>=`、`max<`、`max<=`、`in`、`not-in`、`ex:fn-a,fn-b,fn-c`
8. 基础检查方法：
   ```go
   package validatorV3
   
   import (
   	"testing"
   	"time"
   
   	"github.com/aid297/aid/ptr"
   )
   
   type (
   	UserRequest struct {
   		Firstname string            `v-rule:"(min>0)" v-name:"姓"`
   		Lastname  *string           `v-rule:"(required)(min>2)" v-name:"名"`
   		Birthday  time.Time         `v-rule:"required" v-name:"生日"`
   		Age       int               `v-rule:"(min>0)(max<100)" v-name:"年龄"`
   		Level     float64           `v-rule:"(in:1.0,2.0,3.0,4.4)" v-name:"级别"`
   		Articles  []ArticleRequest  `v-rule:"(required)(min>1)" v-name:"文章"`
   		Articles2 []*ArticleRequest `v-rule:"(required)" v-name:"文章2"`
   	}
   	ArticleRequest struct {
   		Title string `v-rule:"(required)(min>5)" v-name:"标题"`
   	}
   )
   
   func Test1(t *testing.T) {
   	ur := &UserRequest{
   		Firstname: "",
   		Lastname:  ptr.New("三"),
   		Birthday:  time.Time{},
   		Age:       101,
   		Level:     4,
   		Articles:  []ArticleRequest{{Title: "123456"}},
   		Articles2: []*ArticleRequest{{Title: "李四"}},
   	}
   
   	checker := APP.Validator.Once().Checker(ur)
   	checker.Validate()
   	t.Logf("验证是否通过：%v\n", checker.OK())
   	for _, wrong := range checker.Wrongs() {
   		t.Logf("%v\n", wrong)
   	}
   
   	// 验证是否通过：false
   	// [姓] 长度错误 期望：> 0
   	// [名] 长度错误 期望：> 2
   	// [年龄] 长度错误 期望：< 100
   	// [级别] 内容错误 期望：在 [1.0 2.0 3.0 4.4] 之中
   	// [文章] 长度错误 期望：> 1
   	// [文章2.标题] 长度错误 期望：> 5
   }
   ```
9. 提前注册额外测试方法：
   ```go
   package validatorV3
   
   import (
   	"fmt"
   	"testing"
   )
   
   type (
   	UserRequest struct {
   		Firstname string `v-rule:"ex:some-ex-check-fn" v-name:"姓"`
   	}
   )
   
   func someExCheckFn(value any) (err error) {
   	if str, ok := value.(string); !ok || str == "" {
   		return ErrInvalidType
   	} else {
   		if str == "张三" {
   			return fmt.Errorf("名字必须不能是：%s", str)
   		}
   	}
   
   	return nil
   }
   
   func Test1(t *testing.T) {
   	ur := &UserRequest{
   		Firstname: "张三",
   	}
   
   	validator := APP.Validator.Once().RegisterExFn("some-ex-check-fn", someExCheckFn) // validator 是单例，只需要在程序初始化时提前注册好检查方法即可
   
   	checker := validator.Checker(ur)
   	checker.Validate()
   	t.Logf("验证是否通过：%v\n", checker.OK())
   	for _, wrong := range checker.Wrongs() {
   		t.Logf("%v\n", wrong)
   	}
   
   	// 验证是否通过：false
   	// 名字必须不能是：张三
   }
   ```
10. 一次性额外注册：
    ```go
    package validatorV3
    
    import (
    	"fmt"
    	"testing"
    )
    
    type (
    	UserRequest struct {
    		Firstname string `v-rule:"ex:some-ex-check-fn" v-name:"姓"`
    	}
    )
    
    func Test1(t *testing.T) {
    	ur := &UserRequest{
    		Firstname: "张三",
    	}
    
    	validator := APP.Validator.Once()
    
    	checker := validator.Checker(ur)
    	checker.Validate(func(form any) (err error) {
    		// 这里是一次性自定义验证（模拟去数据库中进行验证）
    		if form.(*UserRequest).Firstname != "王五" {
    			err = fmt.Errorf("名字必须是：王五")
    		}
    
    		return
    	})
    	t.Logf("验证是否通过：%v\n", checker.OK())
    	for _, wrong := range checker.Wrongs() {
    		t.Logf("%v\n", wrong)
    	}
    
    	// 验证是否通过：false
    	// 名字必须是：王五
    }
    ```

    
