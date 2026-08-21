
### HttpClient 使用说明

1. 初始化
---
```go
   package main
   
   import (
   	`fmt`
   
   	`github.com/aid297/aid/v3/httpClients`
   )
   
   func main() {
   	hc,err := httpClients.New(httpClients.URL("https://www.baidu.com"))
    if err != nil{
        panci(err) // 初始化客户端错误
    }
   	if err := hc.Send().OK(); err != nil {
   		panic(err) // 发送消息错误
   	}
   
   	res := hc.ToBytes()
   	debugLogs.Print(string(res))
   
   	// 输出结果
   	// <!DOCTYPE html>
   	// <!--STATUS OK--><html> <head><meta http-equiv=content-type content=text/html;charset=utf-8><meta http-equiv=X-UA-Compatible content=IE=Edge><meta content=always name=referrer><link rel=stylesheet type=text/css href=https://ss1.bdstatic.com/5eN1bjq8AAUYm2zgoY3K/r/www/cache/bdorz/baidu.min.css><title>百度一下，你就知道</title></head> <body link=#0000cc> <div id=wrapper> <div id=head> <div class=head_wrapper> <div class=s_form> <div class=s_form_wrapper> <div id=lg> <img hidefocus=true src=//www.baidu.com/img/bd_logo1.png width=270 height=129> </div> <form id=form name=f action=//www.baidu.com/s class=fm> <input type=hidden name=bdorz_come value=1> <input type=hidden name=ie value=utf-8> <input type=hidden name=f value=8> <input type=hidden name=rsv_bp value=1> <input type=hidden name=rsv_idx value=1> <input type=hidden name=tn value=baidu><span class="bg s_ipt_wr"><input id=kw name=wd class=s_ipt value maxlength=255 autocomplete=off autofocus=autofocus></span><span class="bg s_btn_wr"><input type=submit id=su value=百度一下 class="bg s_btn" autofocus></span> </form> </div> </div> <div id=u1> <a href=http://news.baidu.com name=tj_trnews class=mnav>新闻</a> <a href=https://www.hao123.com name=tj_trhao123 class=mnav>hao123</a> <a href=http://map.baidu.com name=tj_trmap class=mnav>地图</a> <a href=http://v.baidu.com name=tj_trvideo class=mnav>视频</a> <a href=http://tieba.baidu.com name=tj_trtieba class=mnav>贴吧</a> <noscript> <a href=http://www.baidu.com/bdorz/login.gif?login&amp;tpl=mn&amp;u=http%3A%2F%2Fwww.baidu.com%2f%3fbdorz_come%3d1 name=tj_login class=lb>登录</a> </noscript> <script>document.write('<a href="http://www.baidu.com/bdorz/login.gif?login&tpl=mn&u='+ encodeURIComponent(window.location.href+ (window.location.search === "" ? "?" : "&")+ "bdorz_come=1")+ '" name="tj_login" class="lb">登录</a>');
   	// </script> <a href=//www.baidu.com/more/ name=tj_briicon class=bri style="display: block;">更多产品</a> </div> </div> </div> <div id=ftCon> <div id=ftConw> <p id=lh> <a href=http://home.baidu.com>关于百度</a> <a href=http://ir.baidu.com>About Baidu</a> </p> <p id=cp>&copy;2017&nbsp;Baidu&nbsp;<a href=http://www.baidu.com/duty/>使用百度前必读</a>&nbsp; <a href=http://jianyi.baidu.com/ class=cp-feedback>意见反馈</a>&nbsp;京ICP证030173号&nbsp; <img src=//www.baidu.com/img/gs.gif> </p> </div> </div> </div> </body> </html>
   }
```
---

2. 其他参数
---
```go
   package main
   
   import (
   	`log`
   	`net/http`
   
   	`github.com/aid297/aid/v3/httpClients`
   	`github.com/aid297/aid/v3/time`
   )
   
   func main() {
   	var wrongs []error
   
   	hc,err := httpClients.New(
   		httpClients.URL("https://", "这里", "为了方便", "可以直接", "使用", "可变参数", "的方式", "传入", "URL"),
   		httpClients.Method(http.MethodPost), // 设置访问方式
      httpClients.ContentType(httpClients.ConentTypeJSON), // 设置 Content-Type 
    httpClients.Authorization("username","password","Basic"),// 设置 认证信息
   	)
    if err != nil{
        // ...
    }
   
   	// 在这里可以追加或覆盖一些新的配置项，比如：
   	hc.SetAttrs(
      Accept(httpClients.AcceptJSON), // 设置 Accept
   		httpClients.Method(http.MethodPut), // 这里可以覆盖之前的属性
   		httpClients.AutoCopy(true),         // 设置自动备份响应体。默认情况下响应体不会自动读取流，所以如果需要在 Send 之后再次读取响应体，就需要设置这个属性为 true，来自动备份响应体内容，以便后续再次读取。
   	)
   
   	// 设置请求体：JSON
    hc.SetAttrs(httpClients.JSON(map[string]any{}))
   
   	// 设置 form 表单数据
   	hc.SetAttrs(httpClients.Form(map[string]any{}))
   
   	// 设置 form-data 表单数据
   	hc.SetAttrs(httpClients.FormData(map[string]string{}, map[string]string{}))
   
   	// 这是一个很有用地请求体，我们可以把用户的请求例如 gin.Context.Request.Body 直接传入这个请求体中
   	hc.SetAttrs(httpClients.ReadCloser(nil))
   
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
   			debugLogs.Print("错误：%v\n", wrong)
   		}
   	}
   }
```
---

3. TLS
---
```go
   func getCAPool() *x509.CertPool {
   	caFile := "ca.crt"
   	caPEM, err := os.ReadFile(caFile)
   	if err != nil {
   		log.Fatalf("读取 CA %s: %v", caFile, err)
   	}
   	caPool := x509.NewCertPool()
   	if !caPool.AppendCertsFromPEM(caPEM) {
   		log.Fatalf("解析 CA PEM 失败: %s", caFile)
   	}
   	return caPool
   }
   
   func Test1(t *testing.T) {
   	t.Run("http client request init", func(t *testing.T) {
   		caPool := getCAPool() // 获取 CA 证书池
   		clientBCrt := "client-b.crt"
   		clientBPriv := "client-b.key"
   		clientBKeyPair, err := tls.LoadX509KeyPair(clientBCrt, clientBPriv)
   		if err != nil {
   			log.Fatalf("加载节点证书 %s + %s: %v", clientBCrt, clientBPriv, err)
   		}
   
   		clientBTLSConfig := &tls.Config{
   			Certificates: []tls.Certificate{clientBKeyPair}, // 设置客户端证书
   			RootCAs:      caPool,                            // 设置 CA 证书池(所有出站请求都必须在此信任链中)
   			MinVersion:   tls.VersionTLS12,                  // 设置最小 TLS 版本
   		}
   
   		transport := &http.Transport{
   			DisableKeepAlives:   true,             // 禁用连接复用
   			MaxIdleConns:        100,              // 最大空闲连接数
   			IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
   			TLSHandshakeTimeout: 10 * time.Second, // TLS 握手超时时间
   			TLSClientConfig:     clientBTLSConfig, // TLS 证书
   		}
   
   		hc := httpClients.New(
   			httpClients.URL("http://www", ".baidu", ".com"),
   			httpClients.Method(http.MethodGet),
   			httpClients.Queries(map[string]any{"name": "张三", "age": 18}),
   			httpClients.Authorization("username", "password", "Basic"),
        httpClients.Accept(AcceptJSON),
        httpClients.ContentType(ContentTypeJSON),
   			httpClients.JSON(map[string]any{"李四": 20, "王五": 30, "赵六": 40}),
   			httpClients.Timeout(5*time.Minute),
   			httpClients.Transport(transport),
   			httpClients.Cert(nil),
   			httpClients.AutoCopy(false),
   		)
   
   		t.Logf("%+v\n", hc)
   		t.Logf("url: %s\n", hc.GetURL())
   		t.Logf("method: %s\n", hc.GetMethod())
   		t.Logf("queries: %+v\n", hc.GetQueries())
   		t.Logf("headers: %+v\n", hc.GetHeaders())
   		t.Logf("body: %s\n", string(hc.GetBody()))
   		t.Logf("timeout: %s\n", hc.GetTimeout())
   		t.Logf("transport: %+v\n", hc.GetTransport())
   		t.Logf("error: %+v\n", hc.Error())
   
   		t.Logf("response: %s\n", hc.Send().ToBytes())
   	})
   }
```
---

4. 加密
```go
func Test3(t *testing.T) {
	t.Run("encrypt request body with AES", func(t *testing.T) {
		// 创建加密密钥和 IV
		key := []byte{} // 16 字节 = 128 位
		iv := []byte{}  // 16 字节 IV

		// 创建 AES 加密器（使用 CBC 模式）
		aesCipher, err := aes.New(
			aes.RandKey(&key),
			aes.RandIV(&iv),
			aes.AlgorithmCBC(),
		)
		if err != nil {
			t.Fatalf("创建 AES 加密器失败: %v", err)
		}

		// 原始请求体内容
		originalBody := map[string]any{"username": "admin", "password": "secret123"}

		// 创建 HTTPClient 并应用加密
		hc := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/login"),
			httpClients.Method(http.MethodPost),
			httpClients.JSON(originalBody),
			httpClients.Encrypt(aesCipher),
		)

		// 验证没有错误
		if err := hc.Error(); err != nil {
			t.Fatalf("HTTPClient 创建失败: %v", err)
		}

		// 验证加密后的请求体
		encryptedBody := hc.GetBody()
		if len(encryptedBody) == 0 {
			t.Fatal("加密后的请求体为空")
		}

		// 验证加密后的内容与原始内容不同（简单验证）
		originalBodyStr := `{"password":"secret123","username":"admin"}`
		if string(encryptedBody) == originalBodyStr {
			t.Fatal("请求体未被加密")
		}

		// 验证 Content-Encoding 头已设置
		headers := hc.GetHeaders()
		if contentEncoding, ok := headers["Content-Encoding"]; !ok {
			t.Fatal("Content-Encoding 头未设置")
		} else if len(contentEncoding) == 0 || contentEncoding[0] != "encrypted" {
			t.Fatalf("Content-Encoding 头值不正确: %v", contentEncoding)
		}

		t.Logf("✓ 加密成功")
		t.Logf("  原始大小: %d 字节", len(originalBodyStr))
		t.Logf("  加密后大小: %d 字节", len(encryptedBody))
		t.Logf("  Content-Encoding: %v", headers["Content-Encoding"])
	})

	t.Run("encrypt with GCM mode", func(t *testing.T) {
		// 使用 GCM 模式进行加密（更安全，包含认证）
		key := []byte{}
		gcmNonce := []byte{} // 12 字节 GCM nonce

		aesCipher, err := aes.New(
			aes.KeySize(aes.AES128),
			aes.RandKey(&key),
			aes.RandGCMNonce(&gcmNonce),
			aes.AlgorithmGCM(),
		)
		if err != nil {
			t.Fatalf("创建 AES GCM 加密器失败: %v", err)
		}

		hc := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/secure"),
			httpClients.Method(http.MethodPost),
			httpClients.JSON(map[string]any{"data": "sensitive information"}),
			httpClients.Encrypt(aesCipher),
		)

		if err := hc.Error(); err != nil {
			t.Fatalf("GCM 模式加密失败: %v", err)
		}

		encryptedBody := hc.GetBody()
		if len(encryptedBody) == 0 {
			t.Fatal("GCM 模式加密后的请求体为空")
		}

		t.Logf("✓ GCM 模式加密成功")
		t.Logf("  加密后大小: %d 字节 (包含 nonce 12 字节 + tag 16 字节)", len(encryptedBody))
	})

	t.Run("encrypt with CTR mode", func(t *testing.T) {
		// 使用 CTR 模式进行流加密
		key := []byte{}
		ctrNonce := []byte{} // 16 字节 CTR nonce

		aesCipher, err := aes.New(
			aes.KeySize(aes.AES128),
			aes.RandKey(&key),
			aes.RandCTRNonce(&ctrNonce), // CTR 模式使用 nonce 作为 IV
			aes.AlgorithmCTR(),
		)
		if err != nil {
			t.Fatalf("创建 AES CTR 加密器失败: %v", err)
		}

		hc := httpClients.APP.HTTPClient.New(
			httpClients.URL("http://127.0.0.1:8080/api/stream"),
			httpClients.Method(http.MethodPost),
			httpClients.Plain("This is a large streaming message that needs encryption"),
			httpClients.Encrypt(aesCipher),
		)

		if err := hc.Error(); err != nil {
			t.Fatalf("CTR 模式加密失败: %v", err)
		}

		encryptedBody := hc.GetBody()
		if len(encryptedBody) == 0 {
			t.Fatal("CTR 模式加密后的请求体为空")
		}

		t.Logf("✓ CTR 模式加密成功")
	t.Logf("  加密后大小: %d 字节 (包含 nonce 16 字节)", len(encryptedBody))
})
```
---

5. 文件上传
---
```go
package main

import (
	"github.com/aid297/aid/v3/compressions/zlib"
	"github.com/aid297/aid/v3/httpClients"
	"github.com/aid297/aid/v3/secrets/symmetric/aes"
)

func main() {
	// 普通文件上传（不使用切块）
	hc, err := httpClients.New(
		httpClients.URL("http://127.0.0.1:8080/api/upload"),
		httpClients.Method(http.MethodPost),
		httpClients.File("path/to/file.txt", 0), // goroutineCount = 0，不使用切块
	)
	if err != nil {
		panic(err)
	}
	hc.Send()

	// 大文件切块上传（4个协程并行）
	hc, err = httpClients.New(
		httpClients.URL("http://127.0.0.1:8080/api/upload"),
		httpClients.Method(http.MethodPost),
		httpClients.File("path/to/large_file.zip", 4), // goroutineCount = 4，使用4个协程切块上传
	)
	if err != nil {
		panic(err)
	}
	hc.Send()

	// 文件上传 + 压缩
	compressor, _ := zlib.New()
	hc, err = httpClients.New(
		httpClients.URL("http://127.0.0.1:8080/api/upload"),
		httpClients.Method(http.MethodPost),
		httpClients.File("path/to/file.txt", 0),
		httpClients.Compressor(compressor), // 启用压缩
	)
	if err != nil {
		panic(err)
	}
	hc.Send()

	// 文件上传 + 加密
	key := []byte{}
	iv := []byte{}
	aesCipher, _ := aes.New(
		aes.RandKey(&key),
		aes.RandIV(&iv),
		aes.AlgorithmCBC(),
	)
	hc, err = httpClients.New(
		httpClients.URL("http://127.0.0.1:8080/api/upload"),
		httpClients.Method(http.MethodPost),
		httpClients.File("path/to/file.txt", 0),
		httpClients.Encrypt(aesCipher), // 启用加密
	)
	if err != nil {
		panic(err)
	}
	hc.Send()

	// 完整配置：大文件切块上传 + 压缩 + 加密
	hc, err = httpClients.New(
		httpClients.URL("http://127.0.0.1:8080/api/upload"),
		httpClients.Method(http.MethodPost),
		httpClients.File("path/to/large_file.zip", 4),    // 4个协程切块上传
		httpClients.Compressor(compressor),              // 启用压缩
		httpClients.Encrypt(aesCipher),                  // 启用加密
	)
	if err != nil {
		panic(err)
	}
	hc.Send()
}
```

**File 函数参数说明**：
- `filename`: 文件路径
- `goroutineCount`: 协程数（0 或 1 = 不使用切块，> 1 = 使用切块并行上传）

**处理流程**：
1. **不使用切块**（goroutineCount = 0 或 1）：直接读取整个文件到内存，然后根据配置进行压缩和/或加密
2. **使用切块**（goroutineCount > 1）：将文件分成多个块，使用指定数量的协程并行读取和处理每个块，然后合并发送

**注意事项**：
- 压缩和加密会按顺序应用：先压缩，后加密
- 切块模式下，每个块会独立进行压缩和加密处理
- 切块模式会自动设置以下请求头：
  - `X-File-Size`: 文件总大小
  - `X-Chunk-Count`: 块总数
  - `X-Chunk-Size`: 每块大小