### 表单验证器

1. 简单使用
   ```go
   package main
   
   import (
   	`fmt`
   
   	`github.com/aid297/aid/validator/validatorV3`
   )
   
   type UserRequest struct {
   	Firstname string  `json:"firstname" v-rule:"(required)(not-empty)" v-name:"firstname"`
   	Lastname  *string `json:"lastname" v-rule:"(required)" v-name:"lastname"`
   }
   
   func main() {
   	userRequest := UserRequest{}
   
   	wrong := validatorV3.APP.Validator.Once().Checker(&userRequest).Validate().WrongToString("\n")
   	fmt.Printf("错误：\n%v", wrong)
   	// 错误：
     // 问题1：[firstname] 不能为空 注意：如果字符串，数字这种类型的变量零值不会认为是空，只有字段不存在时才会认为是空，所以需要加入not-empty或者min>0的限制
   	// 问题2：[lastname] 不能为空
   }
   ```

   