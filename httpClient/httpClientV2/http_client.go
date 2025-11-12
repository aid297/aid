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

	"github.com/aid297/aid/operation"
	"github.com/aid297/aid/str"
)

type (
	HttpClient struct {
		err                       error
		url                       string
		queries                   map[string]any
		method                    string
		headers                   map[string][]string
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

	HttpClientBuilder struct {
		options []HttpClientAttributer
	}
)

func (*HttpClient) init(method string, attrs ...HttpClientAttributer) *HttpClient {
	ins := new(HttpClient)
	ins.SetAttrs(Method(method))
	ins.SetAttrs(AppendHeaderValues(map[string][]string{}))

	return ins.SetAttrs(attrs...)
}

func (*HttpClientBuilder) New(options ...HttpClientAttributer) *HttpClientBuilder {
	return &HttpClientBuilder{options: options}
}

func (my *HttpClientBuilder) GetClient() *HttpClient {
	return new(HttpClient).init(http.MethodGet, my.options...)
}

func (*HttpClient) New(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodGet, attrs...)
}

func (*HttpClient) NewGet(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodGet, attrs...)
}

func (*HttpClient) NewPost(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodPost, attrs...)
}

func (*HttpClient) NewPut(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodPut, attrs...)
}

func (*HttpClient) NewPatch(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodPatch, attrs...)
}

func (*HttpClient) NewDelete(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodDelete, attrs...)
}

func (*HttpClient) NewHead(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodHead, attrs...)
}

func (*HttpClient) NewOptions(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodOptions, attrs...)
}

func (*HttpClient) NewTrace(attrs ...HttpClientAttributer) *HttpClient {
	return new(HttpClient).init(http.MethodTrace, attrs...)
}

func (my *HttpClient) set(attrs ...HttpClientAttributer) {
	if len(attrs) > 0 {
		for _, option := range attrs {
			option.Register(my)
			if my.err != nil {
				return
			}
		}
	}
}

func (my *HttpClient) SetAttrs(attrs ...HttpClientAttributer) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.set(attrs...)

	return my
}

func (my *HttpClient) GetURL() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getURL()
}

func (my *HttpClient) getURL() string {
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

func (my *HttpClient) GetQueries() map[string]any {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getQueries()
}

func (my *HttpClient) getQueries() map[string]any { return my.queries }

func (my *HttpClient) GetMethod() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getMethod()
}

func (my *HttpClient) getMethod() string { return my.method }

func (my *HttpClient) GetHeaders() map[string][]string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getHeaders()
}

func (my *HttpClient) getHeaders() map[string][]string { return my.headers }

func (my *HttpClient) GetBody() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getBody()
}

func (my *HttpClient) getBody() []byte { return my.requestBody }

func (my *HttpClient) GetTimeout() time.Duration {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTimeout()
}

func (my *HttpClient) getTimeout() time.Duration { return my.timeout }

func (my *HttpClient) GetTransport() *http.Transport {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTransport()
}

func (my *HttpClient) getTransport() *http.Transport { return my.transport }

func (my *HttpClient) GetCert() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getCert()
}

func (my *HttpClient) getCert() []byte { return my.cert }

func (my *HttpClient) GetRawRequest() *http.Request {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawRequest()
}

func (my *HttpClient) getRawRequest() *http.Request { return my.rawRequest }

func (my *HttpClient) GetRawResponse() *http.Response {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawResponse()
}

func (my *HttpClient) getRawResponse() *http.Response { return my.rawResponse }

func (my *HttpClient) GetClient() *http.Client {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getClient()
}

func (my *HttpClient) getClient() *http.Client { return my.client }

func (my *HttpClient) send() *HttpClient {
	if my.err != nil {
		return my
	}

	if my.rawRequest, my.err = http.NewRequest(my.method, my.getURL(), bytes.NewReader(my.requestBody)); my.err != nil {
		return my
	}

	for key, values := range my.headers {
		my.rawRequest.Header[key] = append(my.rawRequest.Header[key], values...)
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

func (my *HttpClient) Send() *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	return my.send()
}

func (my *HttpClient) parseBody() {
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

func (my *HttpClient) ToJson(target any, keys ...any) *HttpClient {
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

	return operation.TernaryFuncAll(
		func() bool { return len(keys) > 0 },
		func() *HttpClient {
			jsonIter.Get(my.responseBody, keys...).ToVal(&target)
			return my
		}, func() *HttpClient {
			my.err = json.Unmarshal(my.responseBody, &target)
			return my
		},
	)
}

func (my *HttpClient) ToXml(target any) *HttpClient {
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

func (my *HttpClient) ToBytes() []byte {
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

func (my *HttpClient) ToWriter(writer http.ResponseWriter) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() { _ = my.rawResponse.Body.Close() }()

	if my.err != nil {
		return my
	}

	_, my.err = io.Copy(writer, my.rawResponse.Body)
	return my
}

func (my *HttpClient) Error() error {
	var err error
	defer func() { my.err = nil }()

	err = my.err
	return err
}
