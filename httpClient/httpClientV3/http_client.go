package httpClientV3

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/aid297/aid/operation"
	"github.com/aid297/aid/str"
	jsonIter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

type HTTPClient struct {
	err                       error
	url                       string
	queries                   map[string]any
	method                    string
	headers                   map[string][]any
	requestBody, responseBody []byte
	timeout                   time.Duration
	transport                 *http.Transport
	cert                      []byte
	rawRequest                *http.Request
	rawResponse               *http.Response
	client                    *http.Client
	autoCopy                  bool
	autoLock                  bool
	lock                      *sync.RWMutex
}

// init 初始化
func (HTTPClient) init(method string, attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{lock: &sync.RWMutex{}, autoLock: true}.SetAttrs(Method(method), AppendHeaderValues(map[string][]any{})).SetAttrs(attrs...)
}

// New 实例化：HTTP客户端
func (HTTPClient) New(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodGet, attrs...)
}

// GET 实例化：HTTP客户端 - GET请求
func (HTTPClient) GET(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodGet, attrs...)
}

// POST 实例化：HTTP客户端 - POST请求
func (HTTPClient) POST(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodPost, attrs...)
}

// PUT 实例化：HTTP客户端 - PUT请求
func (HTTPClient) PUT(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodPut, attrs...)
}

// PATCH 实例化：HTTP客户端 - PATCH请求
func (HTTPClient) PATCH(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodPatch, attrs...)
}

// DELETE 实例化：HTTP客户端 - DELETE请求
func (HTTPClient) DELETE(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodDelete, attrs...)
}

// HEAD 实例化：HTTP客户端 - HEAD请求
func (HTTPClient) HEAD(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodHead, attrs...)
}

// OPTIONS 实例化：HTTP客户端 - OPTIONS请求
func (HTTPClient) OPTIONS(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodOptions, attrs...)
}

// TRACE 实例化：HTTP客户端 - TRACE请求
func (HTTPClient) TRACE(attrs ...HTTPClientAttributer) HTTPClient {
	return HTTPClient{}.init(http.MethodTrace, attrs...)
}

// setAttrs 设置属性
func (my HTTPClient) setAttrs(attrs ...HTTPClientAttributer) HTTPClient {
	if len(attrs) > 0 {
		for _, option := range attrs {
			option.Register(&my)
			if my.err != nil {
				return my
			}
		}
	}

	return my
}

// SetAttrs 设置属性（线程安全）
func (my HTTPClient) SetAttrs(attrs ...HTTPClientAttributer) HTTPClient {
	if my.autoLock {
		my.lock.Lock()
		defer my.lock.Unlock()
	}

	return my.setAttrs(attrs...)
}

// Lock 加锁 - 写锁
func (my HTTPClient) Lock() HTTPClient {
	my.lock.Lock()
	return my
}

// Unlock 解锁 - 写锁
func (my HTTPClient) Unlock() HTTPClient {
	my.lock.Unlock()
	return my
}

// RLock 加锁 - 读锁
func (my HTTPClient) RLock() HTTPClient {
	my.lock.RLock()
	return my
}

// RUnlock 解锁 - 读锁
func (my HTTPClient) RUnlock() HTTPClient {
	my.lock.RUnlock()
	return my
}

// GetURL 获取完整请求URL
func (my HTTPClient) GetURL() string {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getURL()
}

// getURL 获取完整请求URL
func (my HTTPClient) getURL() string {
	queries := url.Values{}
	if len(my.queries) > 0 {
		for k, v := range my.queries {
			queries.Add(k, cast.ToString(v))
		}
	}

	if len(queries) > 0 {
		return str.APP.Buffer.NewString(my.url).S("?").S(queries.Encode()).String()
	}

	return my.url
}

// GetQueries 获取请求查询参数
func (my HTTPClient) GetQueries() map[string]any {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getQueries()
}

// getQueries 获取请求查询参数
func (my HTTPClient) getQueries() map[string]any { return my.queries }

// GetMethod 获取请求方法
func (my HTTPClient) GetMethod() string {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getMethod()
}

// getMethod 获取请求方法
func (my HTTPClient) getMethod() string { return my.method }

// GetHeaders 获取请求头
func (my HTTPClient) GetHeaders() map[string][]any {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getHeaders()
}

// getHeaders 获取请求头
func (my HTTPClient) getHeaders() map[string][]any { return my.headers }

// GetBody 获取请求体
func (my HTTPClient) GetBody() []byte {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getBody()
}

// getBody 获取请求体
func (my HTTPClient) getBody() []byte { return my.requestBody }

// GetTimeout 获取请求超时时间
func (my HTTPClient) GetTimeout() time.Duration {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getTimeout()
}

// getTimeout 获取请求超时时间
func (my HTTPClient) getTimeout() time.Duration { return my.timeout }

// GetTransport 获取HTTP传输配置
func (my HTTPClient) GetTransport() *http.Transport {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getTransport()
}

// getTransport 获取HTTP传输配置
func (my HTTPClient) getTransport() *http.Transport { return my.transport }

// GetCert 获取TLS证书
func (my HTTPClient) GetCert() []byte {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getCert()
}

// getCert 获取TLS证书
func (my HTTPClient) getCert() []byte { return my.cert }

// GetRawRequest 获取原始HTTP请求
func (my HTTPClient) GetRawRequest() *http.Request {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getRawRequest()
}

// getRawRequest 获取原始HTTP请求
func (my HTTPClient) getRawRequest() *http.Request { return my.rawRequest }

// GetRawResponse 获取原始HTTP响应
func (my HTTPClient) GetRawResponse() *http.Response {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getRawResponse()
}

// getRawResponse 获取原始HTTP响应
func (my HTTPClient) getRawResponse() *http.Response { return my.rawResponse }

// GetClient 获取HTTP客户端
func (my HTTPClient) GetClient() *http.Client {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	return my.getClient()
}

// getClient 获取HTTP客户端
func (my HTTPClient) getClient() *http.Client { return my.client }

// send 发送HTTP请求（不加锁）
func (my HTTPClient) send() HTTPClient {
	if my.err != nil {
		return my
	}

	if my.rawRequest, my.err = http.NewRequest(my.method, my.getURL(), bytes.NewReader(my.requestBody)); my.err != nil {
		return my
	}

	for key, values := range my.headers {
		v := make([]string, 0, len(values))
		for idx := range values {
			v = append(v, cast.ToString(values[idx]))
		}
		my.rawRequest.Header[key] = append(my.rawRequest.Header[key], v...)
	}

	if len(my.cert) > 0 {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(my.cert) {
			my.err = errors.New("生成TLS证书失败")
			return my
		}

		my.transport.TLSClientConfig = &tls.Config{RootCAs: certPool}
	}

	my.client = &http.Client{}

	// 发送新的请求
	if my.transport != nil {
		my.client.Transport = my.transport
	}

	// 设置超时
	if my.timeout > 0 {
		my.client.Timeout = my.timeout
	}

	if my.rawResponse, my.err = my.client.Do(my.rawRequest); my.err != nil {
		return my
	}

	if my.autoCopy {
		my = my.parseBody()
		my.rawResponse.Body = io.NopCloser(bytes.NewBuffer(my.responseBody)) // 还原响应体
	}

	return my
}

func (my HTTPClient) SendWithRetry(count uint, interval time.Duration, condition func(statusCode int, err error) bool) HTTPClient {
	if my.autoLock {
		my.lock.Lock()
		defer my.lock.Unlock()
	}

	if my.send().Error() != nil {
		return my
	}

	if count > 0 && interval > 0 {
		var (
			maxAttempts uint = count + 1 // 首次尝试 + 重试次数
			shouldRetry      = false
		)

		for attempt := uint(0); attempt < maxAttempts; attempt++ {
			time.Sleep(interval)

			if condition == nil {
				condition = func(statusCode int, err error) bool { return statusCode > 399 || err != nil }
			}
			shouldRetry = condition(my.rawResponse.StatusCode, my.err)

			if !shouldRetry || attempt == maxAttempts-1 {
				break
			}

			if my.rawResponse != nil && my.rawResponse.Body != nil {
				_ = my.rawResponse.Body.Close()
				my.rawResponse = nil
			}
		}
	}

	return my
}

// Send 发送HTTP请求（线程安全）
func (my HTTPClient) Send() HTTPClient {
	if my.autoLock {
		my.lock.Lock()
		defer my.lock.Unlock()
	}

	return my.send()
}

// parseBody 解析响应体
func (my HTTPClient) parseBody() HTTPClient {
	var (
		buffer  = bytes.NewBuffer([]byte{})
		written int64
	)

	if my.err != nil {
		return my
	}

	my.responseBody = []byte{}

	if my.rawResponse == nil {
		return my
	}

	// 读取新地响应的主体
	if my.rawResponse.ContentLength > 1*1024*1024 { // 1MB
		if written, my.err = io.Copy(buffer, my.rawResponse.Body); my.err != nil {
			return my
		}
		if written < 1 {
			return my
		}
		if buffer.Len() == 0 {
			return my
		}
		my.responseBody = buffer.Bytes()
		return my
	} else {
		if my.responseBody, my.err = io.ReadAll(my.rawResponse.Body); my.err != nil {
			return my
		}
	}

	return my
}

// ToJSON 将响应体转换为JSON对象
func (my HTTPClient) ToJSON(target any, keys ...any) HTTPClient {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		my = my.parseBody()
	}

	if len(my.responseBody) == 0 {
		return my
	}

	return operation.TernaryFuncAll(
		func() bool { return len(keys) > 0 },
		func() HTTPClient {
			jsonIter.Get(my.responseBody, keys...).ToVal(&target)
			return my
		}, func() HTTPClient {
			my.err = json.Unmarshal(my.responseBody, &target)
			return my
		},
	)
}

// ToXML 将响应体转换为XML对象
func (my HTTPClient) ToXML(target any) HTTPClient {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		my = my.parseBody()
	}

	if len(my.responseBody) == 0 {
		return my
	}

	my.err = xml.Unmarshal(my.responseBody, &target)

	return my
}

// ToBytes 将响应体转换为字节切片
func (my HTTPClient) ToBytes() []byte {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

	if my.err != nil {
		return []byte{}
	}

	if my.responseBody == nil {
		my = my.parseBody()
	}

	if len(my.responseBody) == 0 {
		return nil
	}

	return my.responseBody
}

// ToWriter 将响应体写入HTTP响应写入器
func (my HTTPClient) ToWriter(writer http.ResponseWriter) HTTPClient {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

	if my.err != nil {
		return my
	}

	_, my.err = io.Copy(writer, my.rawResponse.Body)
	return my
}

// Error 获取错误信息
func (my HTTPClient) Error() (err error) {
	if my.autoLock {
		my.lock.RLock()
		defer my.lock.RUnlock()
	}

	defer func() { my.err = nil }()

	err = my.err
	return
}
