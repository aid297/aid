package httpClientV4

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
	"strings"
	"sync"
	"time"

	jsonIter "github.com/json-iterator/go"
)

// 使用对象池减少内存分配
var (
	bufferPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

	clientPool = sync.Pool{
		New: func() any {
			return &HTTPClient{
				queries: make(map[string]string, 8),
				headers: make(http.Header, 8),
			}
		},
	}
)

type HTTPClient struct {
	err                       error
	url                       string
	queries                   map[string]string // 使用 string 替代 any 避免接口装箱
	method                    string
	headers                   http.Header // 使用标准库类型,已优化
	requestBody, responseBody []byte
	timeout                   time.Duration
	transport                 *http.Transport
	cert                      []byte
	rawRequest                *http.Request
	rawResponse               *http.Response
	client                    *http.Client
	autoCopy                  bool
	lock                      sync.RWMutex
}

type HTTPClientBuilder struct {
	attrs []HTTPClientAttributer
}

// 获取客户端实例(从池中)
func Acquire() *HTTPClient {
	c := clientPool.Get().(*HTTPClient)
	c.reset()
	return c
}

// 释放客户端实例(归还池)
func (my *HTTPClient) Release() {
	if my == nil {
		return
	}
	my.reset()
	clientPool.Put(my)
}

// 重置客户端状态
func (my *HTTPClient) reset() {
	my.err = nil
	my.url = ""
	my.method = http.MethodGet
	my.timeout = 0
	my.autoCopy = false
	my.requestBody = nil
	my.responseBody = nil
	my.cert = nil
	my.transport = nil
	my.rawRequest = nil
	my.rawResponse = nil
	my.client = nil

	// 清空 map 但保留容量
	for k := range my.queries {
		delete(my.queries, k)
	}
	for k := range my.headers {
		delete(my.headers, k)
	}
}

func (my *HTTPClient) init(method string, attrs ...HTTPClientAttributer) *HTTPClient {
	my.method = method
	if my.headers == nil {
		my.headers = make(http.Header, 8)
	}
	if my.queries == nil {
		my.queries = make(map[string]string, 8)
	}
	my.SetAttrs(attrs...)
	return my
}

func (*HTTPClientBuilder) New(attrs ...HTTPClientAttributer) *HTTPClientBuilder {
	return &HTTPClientBuilder{attrs: attrs}
}

func (my *HTTPClientBuilder) GetClient() *HTTPClient {
	c := Acquire()
	return c.init(http.MethodGet, my.attrs...)
}

func New(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodGet, attrs...)
}

func GET(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodGet, attrs...)
}

func POST(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodPost, attrs...)
}

func PUT(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodPut, attrs...)
}

func PATCH(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodPatch, attrs...)
}

func DELETE(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodDelete, attrs...)
}

func HEAD(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodHead, attrs...)
}

func OPTIONS(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodOptions, attrs...)
}

func TRACE(attrs ...HTTPClientAttributer) *HTTPClient {
	c := Acquire()
	return c.init(http.MethodTrace, attrs...)
}

// 内部设置方法,不加锁(调用者负责加锁)
func (my *HTTPClient) setAttrs(attrs ...HTTPClientAttributer) {
	if len(attrs) > 0 {
		for _, option := range attrs {
			option.Apply(my)
			if my.err != nil {
				return
			}
		}
	}
}

func (my *HTTPClient) SetAttrs(attrs ...HTTPClientAttributer) *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setAttrs(attrs...)
	return my
}

func (my *HTTPClient) GetURL() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.buildURL()
}

// 构建 URL(不加锁)
func (my *HTTPClient) buildURL() string {
	if len(my.queries) == 0 {
		return my.url
	}

	queries := url.Values{}
	for k, v := range my.queries {
		queries.Add(k, v)
	}

	// 使用 strings.Builder 减少字符串拼接的内存分配
	var sb strings.Builder
	sb.Grow(len(my.url) + len(queries.Encode()) + 1)
	sb.WriteString(my.url)
	sb.WriteByte('?')
	sb.WriteString(queries.Encode())

	return sb.String()
}

func (my *HTTPClient) GetQueries() map[string]string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	// 返回副本避免外部修改
	result := make(map[string]string, len(my.queries))
	for k, v := range my.queries {
		result[k] = v
	}
	return result
}

func (my *HTTPClient) GetMethod() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.method
}

func (my *HTTPClient) GetHeaders() http.Header {
	my.lock.RLock()
	defer my.lock.RUnlock()

	// 返回副本
	result := make(http.Header, len(my.headers))
	for k, v := range my.headers {
		result[k] = append([]string(nil), v...)
	}
	return result
}

func (my *HTTPClient) GetBody() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.requestBody
}

func (my *HTTPClient) GetTimeout() time.Duration {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.timeout
}

func (my *HTTPClient) GetTransport() *http.Transport {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.transport
}

func (my *HTTPClient) GetCert() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.cert
}

func (my *HTTPClient) GetRawRequest() *http.Request {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.rawRequest
}

func (my *HTTPClient) GetRawResponse() *http.Response {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.rawResponse
}

func (my *HTTPClient) GetClient() *http.Client {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return my.client
}

// 内部发送方法(不加锁)
func (my *HTTPClient) doSend() error {
	if my.err != nil {
		return my.err
	}

	// 创建请求
	req, err := http.NewRequest(my.method, my.buildURL(), bytes.NewReader(my.requestBody))
	if err != nil {
		my.err = err
		return err
	}
	my.rawRequest = req

	// 设置 headers
	req.Header = my.headers

	// 配置 TLS
	if len(my.cert) > 0 {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(my.cert) {
			my.err = errors.New("生成TLS证书失败")
			return my.err
		}

		if my.transport == nil {
			my.transport = &http.Transport{}
		}
		my.transport.TLSClientConfig = &tls.Config{RootCAs: certPool}
	}

	// 创建客户端
	my.client = &http.Client{}

	if my.transport != nil {
		my.client.Transport = my.transport
	}

	if my.timeout > 0 {
		my.client.Timeout = my.timeout
	}

	// 发送请求
	resp, err := my.client.Do(req)
	if err != nil {
		my.err = err
		return err
	}
	my.rawResponse = resp

	// 自动复制响应体
	if my.autoCopy {
		if err := my.parseBody(); err != nil {
			return err
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(my.responseBody))
	}

	return nil
}

func (my *HTTPClient) SendWithRetry(count uint, interval time.Duration, condition func(statusCode int, err error) bool) *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	if err := my.doSend(); err != nil {
		return my
	}

	if count > 0 && interval > 0 {
		maxAttempts := count + 1

		if condition == nil {
			condition = func(statusCode int, err error) bool {
				return statusCode > 399 || err != nil
			}
		}

		for attempt := uint(1); attempt < maxAttempts; attempt++ {
			shouldRetry := condition(my.rawResponse.StatusCode, my.err)
			if !shouldRetry {
				break
			}

			time.Sleep(interval)

			if my.rawResponse != nil && my.rawResponse.Body != nil {
				_ = my.rawResponse.Body.Close()
				my.rawResponse = nil
			}

			if err := my.doSend(); err != nil {
				break
			}
		}
	}

	return my
}

func (my *HTTPClient) Send() *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	_ = my.doSend()
	return my
}

// 解析响应体
func (my *HTTPClient) parseBody() error {
	if my.err != nil {
		return my.err
	}

	my.responseBody = nil

	if my.rawResponse == nil || my.rawResponse.Body == nil {
		return nil
	}

	// 从池中获取 buffer
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// 根据大小选择读取策略
	if my.rawResponse.ContentLength > 1*1024*1024 { // 1MB
		written, err := io.Copy(buf, my.rawResponse.Body)
		if err != nil {
			my.err = err
			return err
		}
		if written < 1 {
			return nil
		}
		if buf.Len() == 0 {
			return nil
		}
		// 复制数据(因为 buffer 会被归还池)
		my.responseBody = make([]byte, buf.Len())
		copy(my.responseBody, buf.Bytes())
	} else {
		data, err := io.ReadAll(my.rawResponse.Body)
		if err != nil {
			my.err = err
			return err
		}
		my.responseBody = data
	}

	return nil
}

func (my *HTTPClient) ToJSON(target any, keys ...any) *HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil && my.rawResponse.Body != nil {
		defer func() { _ = my.rawResponse.Body.Close() }()
	}

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		if err := my.parseBody(); err != nil {
			return my
		}
	}

	if len(my.responseBody) == 0 {
		return my
	}

	if len(keys) > 0 {
		jsonIter.Get(my.responseBody, keys...).ToVal(&target)
	} else {
		my.err = json.Unmarshal(my.responseBody, &target)
	}

	return my
}

func (my *HTTPClient) ToXML(target any) *HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil && my.rawResponse.Body != nil {
		defer func() { _ = my.rawResponse.Body.Close() }()
	}

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		if err := my.parseBody(); err != nil {
			return my
		}
	}

	if len(my.responseBody) == 0 {
		return my
	}

	my.err = xml.Unmarshal(my.responseBody, &target)
	return my
}

func (my *HTTPClient) ToBytes() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil && my.rawResponse.Body != nil {
		defer func() { _ = my.rawResponse.Body.Close() }()
	}

	if my.err != nil {
		return []byte{}
	}

	if my.responseBody == nil {
		if err := my.parseBody(); err != nil {
			return nil
		}
	}

	if len(my.responseBody) == 0 {
		return nil
	}

	return my.responseBody
}

func (my *HTTPClient) ToWriter(writer http.ResponseWriter) *HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil && my.rawResponse.Body != nil {
		defer func() { _ = my.rawResponse.Body.Close() }()
	}

	if my.err != nil {
		return my
	}

	_, my.err = io.Copy(writer, my.rawResponse.Body)
	return my
}

func (my *HTTPClient) Error() error {
	err := my.err
	my.err = nil
	return err
}

func (my *HTTPClient) GetStatusCode() int {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil {
		return my.rawResponse.StatusCode
	}
	return 0
}

func (my *HTTPClient) GetStatus() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if my.rawResponse != nil {
		return my.rawResponse.Status
	}
	return ""
}
