### HttpClientV2 使用说明

1. 初始化
   ```go
   package main
   
   import (
   	`fmt`
   
   	`github.com/aid297/aid/httpClient/httpClientV2`
   )
   
   func main() {
   	hc := httpClientV2.APP.HTTPClient.New(httpClientV2.URL("https://www.baidu.com"))
   	if err := hc.Send().OK(); err != nil {
   		panic(err)
   	}
   
   	res := hc.ToBytes()
   	fmt.Println(string(res))
   
   	// 输出结果
   	// <!DOCTYPE html>
   	// <!--STATUS OK--><html> <head><meta http-equiv=content-type content=text/html;charset=utf-8><meta http-equiv=X-UA-Compatible content=IE=Edge><meta content=always name=referrer><link rel=stylesheet type=text/css href=https://ss1.bdstatic.com/5eN1bjq8AAUYm2zgoY3K/r/www/cache/bdorz/baidu.min.css><title>百度一下，你就知道</title></head> <body link=#0000cc> <div id=wrapper> <div id=head> <div class=head_wrapper> <div class=s_form> <div class=s_form_wrapper> <div id=lg> <img hidefocus=true src=//www.baidu.com/img/bd_logo1.png width=270 height=129> </div> <form id=form name=f action=//www.baidu.com/s class=fm> <input type=hidden name=bdorz_come value=1> <input type=hidden name=ie value=utf-8> <input type=hidden name=f value=8> <input type=hidden name=rsv_bp value=1> <input type=hidden name=rsv_idx value=1> <input type=hidden name=tn value=baidu><span class="bg s_ipt_wr"><input id=kw name=wd class=s_ipt value maxlength=255 autocomplete=off autofocus=autofocus></span><span class="bg s_btn_wr"><input type=submit id=su value=百度一下 class="bg s_btn" autofocus></span> </form> </div> </div> <div id=u1> <a href=http://news.baidu.com name=tj_trnews class=mnav>新闻</a> <a href=https://www.hao123.com name=tj_trhao123 class=mnav>hao123</a> <a href=http://map.baidu.com name=tj_trmap class=mnav>地图</a> <a href=http://v.baidu.com name=tj_trvideo class=mnav>视频</a> <a href=http://tieba.baidu.com name=tj_trtieba class=mnav>贴吧</a> <noscript> <a href=http://www.baidu.com/bdorz/login.gif?login&amp;tpl=mn&amp;u=http%3A%2F%2Fwww.baidu.com%2f%3fbdorz_come%3d1 name=tj_login class=lb>登录</a> </noscript> <script>document.write('<a href="http://www.baidu.com/bdorz/login.gif?login&tpl=mn&u='+ encodeURIComponent(window.location.href+ (window.location.search === "" ? "?" : "&")+ "bdorz_come=1")+ '" name="tj_login" class="lb">登录</a>');
   	// </script> <a href=//www.baidu.com/more/ name=tj_briicon class=bri style="display: block;">更多产品</a> </div> </div> </div> <div id=ftCon> <div id=ftConw> <p id=lh> <a href=http://home.baidu.com>关于百度</a> <a href=http://ir.baidu.com>About Baidu</a> </p> <p id=cp>&copy;2017&nbsp;Baidu&nbsp;<a href=http://www.baidu.com/duty/>使用百度前必读</a>&nbsp; <a href=http://jianyi.baidu.com/ class=cp-feedback>意见反馈</a>&nbsp;京ICP证030173号&nbsp; <img src=//www.baidu.com/img/gs.gif> </p> </div> </div> </div> </body> </html>
   }
   ```

2. 其他参数
   ```go
   package main
   
   import (
   	`log`
   	`net/http`
   
   	`github.com/aid297/aid/httpClient/httpClientV2`
   	`github.com/aid297/aid/time`
   )
   
   func main() {
   	var wrongs []error
   
   	hc := httpClientV2.APP.HTTPClient.New(
   		httpClientV2.URL("https://", "这里", "为了方便", "可以直接", "使用", "可变参数", "的方式", "传入", "URL"),
   		httpClientV2.Method(http.MethodPost), // 设置访问方式
   		httpClientV2.SetHeaderValues(map[string][]any{}). // 设置请求头（如果有同一个 key 则进行值覆盖）
   			ContentType(httpClientV2.ContentTypeJSON). // 设置 Content-Type
   			Authorization("username", "password", "Basic"), // 设置认证
   	)
   
   	// 在这里可以追加或覆盖一些新的配置项，比如：
   	hc.SetAttrs(
   		httpClientV2.SetHeaderValue(map[string]any{}). // 设置请求投（如果有同一个 key 则进行值追加）
   			Accept(httpClientV2.AcceptJSON), // 设置 Accept
   		httpClientV2.Method(http.MethodPut), // 这里可以覆盖之前的属性
   		httpClientV2.AutoCopy(true),         // 设置自动备份响应体。默认情况下响应体不会自动读取流，所以如果需要在 Send 之后再次读取响应体，就需要设置这个属性为 true，来自动备份响应体内容，以便后续再次读取。
   	)
   
   	// 设置请求体：JSON
   	jsonBody := httpClientV2.JSON(map[string]any{})
   	hc.SetAttrs(jsonBody)
   
   	// 设置 form 表单数据
   	formBody := httpClientV2.Form(map[string]any{})
   	hc.SetAttrs(formBody)
   
   	// 设置 form-data 表单数据
   	formData := httpClientV2.FormData(map[string]string{}, map[string]string{})
   	hc.SetAttrs(formData)
   
   	// 这是一个很有用地请求体，我们可以把用户的请求例如 gin.Context.Request.Body 直接传入这个请求体中
   	readerBody := httpClientV2.Reader(nil)
   	hc.SetAttrs(readerBody)
   
   	if err := hc.Send().OK(); err != nil {
   		panic(err)
   	}
   
   	res := hc.ToBytes() // 读取 []byte 响应体
   	hc.ToJSON(&res)     // 读取 JSON 响应体并解析成 map[string]any
   	hc.ToXML(&res)      // 读取 XML 响应体并解析成 map[string]any
   	hc.ToWriter(nil)    // 通常与 Reader 配合使用，我们可以直接将用户请求流转发出去，也可以将响应流直接写入用户请求的响应流中
   
   	hc, wrongs = hc.SendWithRetry(5, 10*time.Second, func(statusCode int, err error) bool {
   		return statusCode != http.StatusOK || err != nil // 假设 status code 不等于 200 或者发生了错误都需要重试
   	})
   	// wrongs: 在重试过程中产生的错误
   	if len(wrongs) > 0 {
   		for _, wrong := range wrongs {
   			log.Printf("错误：%v\n", wrong)
   		}
   	}
   }
   ```

   