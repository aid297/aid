### 配置文件

```go
package main

import (
	`log`

	`github.com/aid297/aid/setting`
	`github.com/aid297/aid/str`
	`github.com/aid297/aid/web-site/backend/aid-web-backend/src/global`
	`github.com/fsnotify/fsnotify`
	`github.com/spf13/viper`
)

func main() {
	var config any

	_, err := setting.NewSetting(
		setting.Filename("local.yaml"),   // 文件名
		setting.EnvName("CONFIG"),        // 环境变量
		setting.DefaultName("test.yaml"), // 默认名
		setting.Content(&config),         // 获取配置文件
		setting.OnChange(func(v *viper.Viper, e fsnotify.Event) {
			log.Println("配置文件发生变化")
			if err := v.Unmarshal(&global.CONFIG); err != nil {
				global.LOG.Error(str.APP.Buffer.JoinString("更新配置文件失败：", err.Error()))
			}
		}),
	)
	if err != nil {
		panic(err)
	}
}
```

