### 配置文件

```go
package main

import (
	`log`

	`github.com/aid297/aid/v2/settings`
	`github.com/aid297/aid/v2/str`
	`github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global`
	`github.com/fsnotify/fsnotify`
	`github.com/spf13/viper`
)

func main() {
	var config any

	_, err := settings.APP.New(
		settings.Filename("local.yaml"),   // 文件名
		settings.EnvName("CONFIG"),        // 环境变量
		settings.DefaultName("test.yaml"), // 默认名
		settings.Content(&config),         // 获取配置文件
		settings.OnChange(func(v *viper.Viper, e fsnotify.Event) {
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

