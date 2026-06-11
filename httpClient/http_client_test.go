package httpClient_test

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aid297/aid/v2/httpClient"
	"github.com/aid297/aid/v2/secret/symmetric/aes"
)

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

		hc, err := httpClient.New(
			httpClient.URL("http://www", ".baidu", ".com"),
			httpClient.Method(http.MethodGet),
			httpClient.Queries(map[string]any{"name": "张三", "age": 18}),
			httpClient.Authorization("username", "password", "Basic"),
			httpClient.Accept(httpClient.AcceptJSON),
			httpClient.ContentType(httpClient.ContentTypeJSON),
			httpClient.JSON(map[string]any{"李四": 20, "王五": 30, "赵六": 40}),
			httpClient.Timeout(5*time.Minute),
			httpClient.Transport(transport),
			httpClient.Cert(nil),
			httpClient.AutoCopy(false),
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

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

func Test2(t *testing.T) {
	hc, err := httpClient.New(
		httpClient.TransportDefault(),
		httpClient.URL("http://127.0.0.1:19003/project/list"),
		httpClient.Header("User-Info", "eyJ1dWlkIjoiYmIwZjJhYjItYTRiNy00ZjYxLWIzYTQtMWRlNTMzOGNkNmNkIiwibmlja25hbWUiOiLkvZkiLCJ1c2VybmFtZSI6Inl1aml6aG91IiwiZW1haWwiOiIxMzUyMjE3ODA1N0BlbWFpbC5jb20iLCJpc0FkbWluIjpmYWxzZSwidGVhbUlkIjoyLCJvd25lclRlYW1JZHMiOlszM10sImhhc1RlYW1JZHMiOlsyLDMzLDM3XSwidG9rZW4iOiIvZWREdnBCRGhManlaVWJ0TC9iVkdaZktRWlRoajJNdDc1bVBxSVduTWxQdFFGOGdxbWpJQjEzMG5MVDllelF2SmJvN1dGbG9YVzU3SW5JZkQvNkRrMU1ERmtpWVJ5aHdZRENSanZJVnArZzY2bnFwSTd5bDBCcWpDN0FBU2NrT3cyQzFWSmY1emtjaGVqbWxIMDhJNnB5Ylk2NmtzaEwwOWxGNlJMZVIzd0xNQ3l5RGNjSVpmclJQQS9IOUtNM3YvNWdVTFk3UGpKL1BSR0NzSzJlYkMyTHlEdGpTMk02MmJ1N0FRekhRQmhkSFEvN3hsQXlUTk1aT0NPU0tyWlRQWmVWK0V5b2NMMGNqNEozQ0ZzMGRFZHkvdG5iVWV3SFhIa0grU0lvUG0vQ2pPeVN4SFFEV3FCS0plemRiSFdCN2ZUMjZZcWljcjBJVmROOTB3UE1SZGtQSHJvSmVkVHNGQSs2SGhVZ1hsVHc9In0="),
		httpClient.Timeout(time.Second),
		httpClient.Method(http.MethodPost),
		httpClient.JSON(map[string]any{"projectUUID": "1f06786e-07bd-6868-8ef5-355bce72ed9b"}),
		httpClient.AutoCopy(true),
	)
	if err != nil {
		t.Fatalf("初始化 HTTP 客户端失败：%v", err)
	}
	if err = hc.Send().OK(); err != nil {
		t.Fatalf("发送 HTTP 请求失败：%v", err)
	}

	t.Logf("结果：%s\n", hc.ToBytes())
}

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
		hc, err := httpClient.New(
			httpClient.URL("http://127.0.0.1:8080/api/login"),
			httpClient.Method(http.MethodPost),
			httpClient.JSON(originalBody),
			httpClient.Encrypt(aesCipher),
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
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

		hc, err := httpClient.New(
			httpClient.URL("http://127.0.0.1:8080/api/secure"),
			httpClient.Method(http.MethodPost),
			httpClient.JSON(map[string]any{"data": "sensitive information"}),
			httpClient.Encrypt(aesCipher),
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

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

		hc, err := httpClient.New(
			httpClient.URL("http://127.0.0.1:8080/api/stream"),
			httpClient.Method(http.MethodPost),
			httpClient.Plain("This is a large streaming message that needs encryption"),
			httpClient.Encrypt(aesCipher),
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		encryptedBody := hc.GetBody()
		if len(encryptedBody) == 0 {
			t.Fatal("CTR 模式加密后的请求体为空")
		}

		t.Logf("✓ CTR 模式加密成功")
		t.Logf("  加密后大小: %d 字节 (包含 nonce 16 字节)", len(encryptedBody))
	})
}
