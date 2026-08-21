package httpClients_test

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aid297/aid/v3/compressions/zlib"
	"github.com/aid297/aid/v3/httpClients"
	"github.com/aid297/aid/v3/secrets/symmetric/aes"
	"github.com/aid297/aid/v3/texts/volumer"
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

		hc, err := httpClients.New(
			httpClients.URL("http://www", ".baidu", ".com"),
			httpClients.Method(http.MethodGet),
			httpClients.Queries(map[string]any{"name": "张三", "age": 18}),
			httpClients.Authorization("username", "password", "Basic"),
			httpClients.Accept(httpClients.AcceptJSON),
			httpClients.ContentType(httpClients.ContentTypeJSON),
			httpClients.JSON(map[string]any{"李四": 20, "王五": 30, "赵六": 40}),
			httpClients.Timeout(5*time.Minute),
			httpClients.Transport(transport),
			httpClients.Cert(nil),
			httpClients.AutoCopy(false),
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
	hc, err := httpClients.New(
		httpClients.TransportDefault(),
		httpClients.URL("http://127.0.0.1:19003/project/list"),
		httpClients.Header("User-Info", "eyJ1dWlkIjoiYmIwZjJhYjItYTRiNy00ZjYxLWIzYTQtMWRlNTMzOGNkNmNkIiwibmlja25hbWUiOiLkvZkiLCJ1c2VybmFtZSI6Inl1aml6aG91IiwiZW1haWwiOiIxMzUyMjE3ODA1N0BlbWFpbC5jb20iLCJpc0FkbWluIjpmYWxzZSwidGVhbUlkIjoyLCJvd25lclRlYW1JZHMiOlszM10sImhhc1RlYW1JZHMiOlsyLDMzLDM3XSwidG9rZW4iOiIvZWREdnBCRGhManlaVWJ0TC9iVkdaZktRWlRoajJNdDc1bVBxSVduTWxQdFFGOGdxbWpJQjEzMG5MVDllelF2SmJvN1dGbG9YVzU3SW5JZkQvNkRrMU1ERmtpWVJ5aHdZRENSanZJVnArZzY2bnFwSTd5bDBCcWpDN0FBU2NrT3cyQzFWSmY1emtjaGVqbWxIMDhJNnB5Ylk2NmtzaEwwOWxGNlJMZVIzd0xNQ3l5RGNjSVpmclJQQS9IOUtNM3YvNWdVTFk3UGpKL1BSR0NzSzJlYkMyTHlEdGpTMk02MmJ1N0FRekhRQmhkSFEvN3hsQXlUTk1aT0NPU0tyWlRQWmVWK0V5b2NMMGNqNEozQ0ZzMGRFZHkvdG5iVWV3SFhIa0grU0lvUG0vQ2pPeVN4SFFEV3FCS0plemRiSFdCN2ZUMjZZcWljcjBJVmROOTB3UE1SZGtQSHJvSmVkVHNGQSs2SGhVZ1hsVHc9In0="),
		httpClients.Timeout(time.Second),
		httpClients.Method(http.MethodPost),
		httpClients.JSON(map[string]any{"projectUUID": "1f06786e-07bd-6868-8ef5-355bce72ed9b"}),
		httpClients.AutoCopy(true),
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
		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/login"),
			httpClients.Method(http.MethodPost),
			httpClients.JSON(originalBody),
			httpClients.Encrypt(aesCipher),
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

		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/secure"),
			httpClients.Method(http.MethodPost),
			httpClients.JSON(map[string]any{"data": "sensitive information"}),
			httpClients.Encrypt(aesCipher),
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

		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/stream"),
			httpClients.Method(http.MethodPost),
			httpClients.Plain("This is a large streaming message that needs encryption"),
			httpClients.Encrypt(aesCipher),
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

func Test4(t *testing.T) {
	t.Run("file upload without chunking", func(t *testing.T) {
		// 创建临时测试文件
		tempFile, err := os.CreateTemp("", "test_*.txt")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		tempFileName := tempFile.Name()
		defer os.Remove(tempFileName)

		testContent := "This is a test file content for HTTP client file upload"
		if _, err = tempFile.WriteString(testContent); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		tempFile.Close()

		// 不使用切块模式（goroutineCount = 0）
		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/upload"),
			httpClients.Method(http.MethodPost),
			httpClients.File(tempFileName, 0), // 不使用切块
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		// 验证请求体不为空
		body := hc.GetBody()
		if len(body) == 0 {
			t.Fatal("文件内容未被读取")
		}

		t.Logf("✓ 普通文件上传成功")
		t.Logf("  文件大小: %d 字节", len(body))
	})

	t.Run("file upload with chunking", func(t *testing.T) {
		// 创建临时测试文件（模拟大文件）
		tempFile, err := os.CreateTemp("", "test_chunk_*.dat")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		tempFileName := tempFile.Name()
		defer os.Remove(tempFileName)

		// 写入足够大的内容来测试切块
		largeContent := make([]byte, int(6*volumer.MB)) // 6MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}
		if _, err = tempFile.Write(largeContent); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		tempFile.Close()

		// 可以通过 SetDefaultFileSplitSize 来设置默认的切块大小
		httpClients.SetDefaultFileSplitSize(int64(2 * volumer.MB))

		// 使用切块模式（4个协程）
		_, err = httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/upload"),
			httpClients.Method(http.MethodPost),
			httpClients.File(tempFileName, 4), // 使用4个协程切块上传
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		t.Logf("✓ 切块文件上传配置成功")
		t.Logf("  文件大小: %d 字节", len(largeContent))
		t.Logf("  协程数: 4")
	})

	t.Run("file upload with compression", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_compress_*.txt")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		tempFileName := tempFile.Name()
		defer os.Remove(tempFileName)

		testContent := "This is a test file for compression. Compression is useful for large text files."
		if _, err = tempFile.WriteString(testContent); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		tempFile.Close()

		// 创建压缩器
		compressor, err := zlib.New()
		if err != nil {
			t.Fatalf("创建压缩器失败: %v", err)
		}

		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/upload"),
			httpClients.Method(http.MethodPost),
			httpClients.File(tempFileName, 0),
			httpClients.Compressor(compressor), // 启用压缩
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		compressedBody := hc.GetBody()
		t.Logf("✓ 文件上传+压缩成功")
		t.Logf("  原始大小: %d 字节", len(testContent))
		t.Logf("  压缩后大小: %d 字节", len(compressedBody))
	})

	t.Run("file upload with encryption", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_encrypt_*.txt")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		tempFileName := tempFile.Name()
		defer os.Remove(tempFileName)

		testContent := "Sensitive data that needs encryption"
		if _, err = tempFile.WriteString(testContent); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		tempFile.Close()

		// 创建 AES 加密器
		key := []byte{}
		iv := []byte{}
		aesCipher, err := aes.New(
			aes.RandKey(&key),
			aes.RandIV(&iv),
			aes.AlgorithmCBC(),
		)
		if err != nil {
			t.Fatalf("创建 AES 加密器失败: %v", err)
		}

		hc, err := httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/upload"),
			httpClients.Method(http.MethodPost),
			httpClients.File(tempFileName, 0),
			httpClients.Encrypt(aesCipher), // 启用加密
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		encryptedBody := hc.GetBody()
		t.Logf("✓ 文件上传+加密成功")
		t.Logf("  原始大小: %d 字节", len(testContent))
		t.Logf("  加密后大小: %d 字节", len(encryptedBody))
	})

	t.Run("file upload with chunking, compression and encryption", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test_full_*.dat")
		if err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
		tempFileName := tempFile.Name()
		defer os.Remove(tempFileName)

		largeContent := make([]byte, int(7*volumer.MB)) // 7MB
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}
		if _, err = tempFile.Write(largeContent); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		tempFile.Close()

		// 创建压缩器和加密器
		compressor, err := zlib.New()
		if err != nil {
			t.Fatalf("创建压缩器失败: %v", err)
		}

		key := []byte{}
		iv := []byte{}
		aesCipher, err := aes.New(
			aes.RandKey(&key),
			aes.RandIV(&iv),
			aes.AlgorithmCBC(),
		)
		if err != nil {
			t.Fatalf("创建 AES 加密器失败: %v", err)
		}

		_, err = httpClients.New(
			httpClients.URL("http://127.0.0.1:8080/api/upload"),
			httpClients.Method(http.MethodPost),
			httpClients.File(tempFileName, 4),  // 4个协程切块上传
			httpClients.Compressor(compressor), // 启用压缩
			httpClients.Encrypt(aesCipher),     // 启用加密
		)
		if err != nil {
			t.Fatalf("初始化 HTTP 客户端失败：%v", err)
		}

		t.Logf("✓ 完整配置文件上传成功")
		t.Logf("  文件大小: %d 字节", len(largeContent))
		t.Logf("  协程数: 4")
		t.Logf("  启用压缩: 是")
		t.Logf("  启用加密: 是")
	})
}
