package httpClientV3

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"

	"github.com/aid297/aid/dict/anyMap"
	"github.com/aid297/aid/steam"
	"gopkg.in/yaml.v3"
)

type (
	HTTPClient struct {
		Error        error
		httpBuilder  HTTPBuilder
		headers      map[string][]string
		responseBody []byte
		rawRequest   *http.Request
		rawResponse  *http.Response
		autoCopy     *bool
		lock         sync.Mutex
		OK           *bool
	}

	HTTPClientGetParam struct {
		Queries     HTTPQuery
		Headers     map[string][]string
		Accept      string
		Body        io.Reader
		ContentType string
	}
)

// NewHTTPClientBasic http 客户端：实例化 → 基础
func NewHTTPClientBasic(httpBuilder HTTPBuilder, url string, headers map[string][]string) *HTTPClient {
	return &HTTPClient{
		lock:        sync.Mutex{},
		rawRequest:  nil,
		headers:     headers,
		httpBuilder: httpBuilder,
	}
}

// NewHTTPClientBasic http 客户端：实例化 → 基础，非 GET
func NewHTTPClientNoGET(httpBuilder HTTPBuilder, url string, headers map[string][]string, method string, body io.Reader) (httpClient *HTTPClient, err error) {
	httpClient = NewHTTPClientBasic(httpBuilder, url, headers)
	httpClient.rawRequest, err = http.NewRequest(method, url, body)
	return
}

// NewHTTPClientGET http 客户端：实例化 → GET
func NewHTTPClientGET(httpBuilder HTTPBuilder, url string, headers map[string][]string) (httpClient *HTTPClient, err error) {
	httpClient = NewHTTPClientBasic(httpBuilder, url, headers)
	httpClient.rawRequest, err = http.NewRequest(http.MethodGet, url, nil)
	return
}

// NewHTTPClientPOST http 客户端：实例化 → POST
func NewHTTPClientPOST(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodPost, body)
}

// NewHTTPClientPUT http 客户端：实例化 → PUT
func NewHTTPClientPUT(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodPut, body)
}

// NewHTTPClientPATCH http 客户端：实例化 → PATCH
func NewHTTPClientPATCH(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodPatch, body)
}

// NewHTTPClientDELETE http 客户端：实例化 → DELETE
func NewHTTPClientDELETE(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodDelete, body)
}

// NewHTTPClientHEAD http 客户端：实例化 → HEAD
func NewHTTPClientHEAD(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodHead, body)
}

// NewHTTPClientOPTIONS http 客户端：实例化 → OPTIONS
func NewHTTPClientOPTIONS(httpBuilder HTTPBuilder, url string, headers map[string][]string, body io.Reader) (httpClient *HTTPClient, err error) {
	return NewHTTPClientNoGET(httpBuilder, url, headers, http.MethodOptions, body)
}

// SetAutoCopy 设置超时
func (MY *HTTPClient) SetAutoCopy(autoCopy *bool) *HTTPClient { MY.autoCopy = autoCopy; return MY }

// Send 发送请求
func (my *HTTPClient) Send() *HTTPClient {
	if len(my.headers) > 0 {
		maps.Copy(my.rawRequest.Header, my.headers)
	}

	client := &http.Client{}
	client.Transport = my.httpBuilder.GetTransport()
	client.Timeout = my.httpBuilder.GetTimeout()
	my.rawResponse, my.Error = client.Do(my.rawRequest)

	return my
}
func (my *HTTPClient) Plain() []byte {
	var body []byte

	if my.Error != nil {
		return nil
	}

	if my.rawResponse == nil {
		my.Error = errors.New("响应体为空")
		return nil
	}

	if my.autoCopy != nil && (*my.autoCopy || my.httpBuilder.GetAutoCopy()) {
		my.rawResponse.Body, my.Error = steam.APP.Steam.New(
			steam.ReadCloser(my.rawResponse.Body),
			steam.CopyFn(func(copied []byte) error { body = copied; return nil }),
		).Copy()
		return nil
	}

	defer my.rawResponse.Body.Close()
	if body, my.Error = io.ReadAll(my.rawResponse.Body); my.Error != nil {
		return nil
	}

	return body
}

// JSON 获取 JSON 结果
func (my *HTTPClient) JSON(ret any) *HTTPClient {
	var body []byte

	if my.Error != nil {
		return my
	}

	if my.rawResponse == nil {
		my.Error = errors.New("响应体为空")
		return my
	}

	if my.autoCopy != nil && (*my.autoCopy || my.httpBuilder.GetAutoCopy()) {
		my.rawResponse.Body, my.Error = steam.APP.Steam.New(
			steam.ReadCloser(my.rawResponse.Body),
			steam.CopyFn(func(copied []byte) error { return json.Unmarshal(copied, ret) }),
		).Copy()
		return my
	}

	defer my.rawResponse.Body.Close()
	if body, my.Error = io.ReadAll(my.rawResponse.Body); my.Error != nil {
		return my
	}

	my.Error = json.Unmarshal(body, ret)
	return my
}

func NewHTTPClientParam() *HTTPClientGetParam {
	return &HTTPClientGetParam{Queries: anyMap.New[string, string](), Headers: map[string][]string{}}
}

func (my *HTTPClientGetParam) SetQueries(queries HTTPQuery) *HTTPClientGetParam {
	my.Queries = queries
	return my
}

func (my *HTTPClientGetParam) SetHeaders(headers map[string][]string) *HTTPClientGetParam {
	my.Headers = headers
	return my
}

func (my *HTTPClientGetParam) SetJSON(body any) (*HTTPClientGetParam, error) {
	var (
		err error
		bb  []byte
	)
	if bb, err = json.Marshal(body); err != nil {
		return my, err
	}
	my.Body = bytes.NewBuffer(bb)
	my.ContentType = "application/json"

	return my, nil
}

func (my *HTTPClientGetParam) SetXML(body any) (*HTTPClientGetParam, error) {
	var (
		err error
		bb  []byte
	)
	if bb, err = xml.Marshal(body); err != nil {
		return my, err
	}
	my.Body = bytes.NewBuffer(bb)
	my.ContentType = "application/xml"

	return my, nil
}

func (my *HTTPClientGetParam) SetYAML(body any) (*HTTPClientGetParam, error) {
	var (
		err error
		bb  []byte
	)
	if bb, err = yaml.Marshal(body); err != nil {
		return my, err
	}
	my.Body = bytes.NewBuffer(bb)
	my.ContentType = "application/yaml"

	return my, nil
}

func (my *HTTPClientGetParam) SetForm(body HTTPQuery) (*HTTPClientGetParam, error) {
	if body.IsNotEmpty() {
		my.ContentType = "application/x-www-form-urlencoded"
		params := url.Values{}
		body.Each(func(key, value string) { params.Add(key, value) })
		my.Body = bytes.NewBuffer([]byte(params.Encode()))
	}

	return my, nil
}

func (my *HTTPClientGetParam) SetFormData(body map[string]string) (*HTTPClientGetParam, error) {
	var (
		err    error
		buffer bytes.Buffer
	)

	writer := multipart.NewWriter(&buffer)
	if len(body) > 0 {
		for k, v := range body {
			if err = writer.WriteField(k, v); err != nil {
				return my, errors.New("组织 form-data 参数错误")
			}
		}
	}
	writer.Close()

	my.Body = &buffer
	my.ContentType = writer.FormDataContentType()

	return my, nil
}
