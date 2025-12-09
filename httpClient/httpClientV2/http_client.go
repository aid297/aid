package httpClientV2

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

	jsonIter "github.com/json-iterator/go"
	"github.com/spf13/cast"

	"github.com/aid297/aid/operation/operationV2"
	"github.com/aid297/aid/str"
)

type (
	HTTPClient struct {
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
		lock                      sync.RWMutex
	}

	HTTPClientBuilder struct {
		attrs []HTTPClientAttributer
	}
)

func (*HTTPClient) init(method string, attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).SetAttrs(Method(method), AppendHeaderValues(map[string][]any{})).SetAttrs(attrs...)
}

func (*HTTPClientBuilder) New(attrs ...HTTPClientAttributer) *HTTPClientBuilder {
	return &HTTPClientBuilder{attrs: attrs}
}

func (my *HTTPClientBuilder) GetClient() *HTTPClient {
	return new(HTTPClient).init(http.MethodGet, my.attrs...)
}

func (*HTTPClient) New(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodGet, attrs...)
}

func (*HTTPClient) GET(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodGet, attrs...)
}

func (*HTTPClient) POST(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodPost, attrs...)
}

func (*HTTPClient) PUT(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodPut, attrs...)
}

func (*HTTPClient) PATCH(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodPatch, attrs...)
}

func (*HTTPClient) DELETE(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodDelete, attrs...)
}

func (*HTTPClient) HEAD(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodHead, attrs...)
}

func (*HTTPClient) OPTIONS(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodOptions, attrs...)
}

func (*HTTPClient) TRACE(attrs ...HTTPClientAttributer) *HTTPClient {
	return new(HTTPClient).init(http.MethodTrace, attrs...)
}

func (my *HTTPClient) set(attrs ...HTTPClientAttributer) {
	if len(attrs) > 0 {
		for _, option := range attrs {
			option.Register(my)
			if my.err != nil {
				return
			}
		}
	}
}

func (my *HTTPClient) SetAttrs(attrs ...HTTPClientAttributer) *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.set(attrs...)

	return my
}

func (my *HTTPClient) GetURL() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getURL()
}

func (my *HTTPClient) getURL() string {
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

func (my *HTTPClient) GetQueries() map[string]any {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getQueries()
}

func (my *HTTPClient) getQueries() map[string]any { return my.queries }

func (my *HTTPClient) GetMethod() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getMethod()
}

func (my *HTTPClient) getMethod() string { return my.method }

func (my *HTTPClient) GetHeaders() map[string][]any {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getHeaders()
}

func (my *HTTPClient) getHeaders() map[string][]any { return my.headers }

func (my *HTTPClient) GetBody() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getBody()
}

func (my *HTTPClient) getBody() []byte { return my.requestBody }

func (my *HTTPClient) GetTimeout() time.Duration {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTimeout()
}

func (my *HTTPClient) getTimeout() time.Duration { return my.timeout }

func (my *HTTPClient) GetTransport() *http.Transport {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTransport()
}

func (my *HTTPClient) getTransport() *http.Transport { return my.transport }

func (my *HTTPClient) GetCert() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getCert()
}

func (my *HTTPClient) getCert() []byte { return my.cert }

func (my *HTTPClient) GetRawRequest() *http.Request {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawRequest()
}

func (my *HTTPClient) getRawRequest() *http.Request { return my.rawRequest }

func (my *HTTPClient) GetRawResponse() *http.Response {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawResponse()
}

func (my *HTTPClient) getRawResponse() *http.Response { return my.rawResponse }

func (my *HTTPClient) GetClient() *http.Client {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getClient()
}

func (my *HTTPClient) getClient() *http.Client { return my.client }

func (my *HTTPClient) send() *HTTPClient {
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
		my.parseBody()
		my.rawResponse.Body = io.NopCloser(bytes.NewBuffer(my.responseBody)) // 还原响应体
	}

	return my
}

func (my *HTTPClient) SendWithRetry(count uint, interval time.Duration, condition func(statusCode int, err error) bool) *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

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

func (my *HTTPClient) Send() *HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	return my.send()
}

func (my *HTTPClient) parseBody() {
	var (
		buffer  = bytes.NewBuffer([]byte{})
		written int64
	)

	if my.err != nil {
		return
	}

	my.responseBody = []byte{}

	if my.rawResponse == nil {
		return
	}

	// 读取新地响应的主体
	if my.rawResponse.ContentLength > 1*1024*1024 { // 1MB
		if written, my.err = io.Copy(buffer, my.rawResponse.Body); my.err != nil {
			return
		}
		if written < 1 {
			return
		}
		if buffer.Len() == 0 {
			return
		}
		my.responseBody = buffer.Bytes()
	} else {
		if my.responseBody, my.err = io.ReadAll(my.rawResponse.Body); my.err != nil {
			return
		}
	}
}

func (my *HTTPClient) ToJSON(target any, keys ...any) *HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() { _ = my.rawResponse.Body.Close() }()

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		my.parseBody()
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
	defer func() { _ = my.rawResponse.Body.Close() }()

	if my.err != nil {
		return my
	}

	if my.responseBody == nil {
		my.parseBody()
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
	defer func() { _ = my.rawResponse.Body.Close() }()

	if my.err != nil {
		return []byte{}
	}

	if my.responseBody == nil {
		my.parseBody()
	}

	if len(my.responseBody) == 0 {
		return nil
	}

	return my.responseBody
}

func (my *HTTPClient) ToWriter(writer http.ResponseWriter) *HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() { _ = my.rawResponse.Body.Close() }()

	if my.err != nil {
		return my
	}

	_, my.err = io.Copy(writer, my.rawResponse.Body)
	return my
}

func (my *HTTPClient) Error() error {
	var err error
	defer func() { my.err = nil }()

	err = my.err
	return err
}

func (my *HTTPClient) GetStatusCode() int {
	return operationV2.NewTernary(operationV2.TrueFn(func() int { return my.GetRawResponse().StatusCode })).GetByValue(my.GetRawResponse() != nil)
}

func (my *HTTPClient) GetStatus() string {
	return operationV2.NewTernary(operationV2.TrueFn(func() string { return my.GetRawResponse().Status })).GetByValue(my.GetRawResponse() != nil)
}
