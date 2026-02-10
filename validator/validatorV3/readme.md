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
8. 
