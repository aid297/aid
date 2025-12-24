package httpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	jsonIter "github.com/json-iterator/go"

	"github.com/aid297/aid/operation"
	"github.com/aid297/aid/str"
)

type (
	// HttpClient http客户端
	HttpClient struct {
		Err                error
		requestUrl         string
		requestQueries     map[string]string
		requestMethod      string
		requestBody        []byte
		requestHeaders     map[string][]string
		request            *http.Request
		response           *http.Response
		responseBody       []byte
		responseBodyBuffer *bytes.Buffer
		isReady            bool
		cert               []byte
		transport          *http.Transport
		timeoutSecond      int64
		lock               sync.RWMutex
	}
)

var App HttpClient

func (*HttpClient) New(url string) *HttpClient       { return NewHttpClient(url) }
func (*HttpClient) NewGet(url string) *HttpClient    { return NewGet(url) }
func (*HttpClient) NewPost(url string) *HttpClient   { return NewPost(url) }
func (*HttpClient) NewPut(url string) *HttpClient    { return NewPut(url) }
func (*HttpClient) NewDelete(url string) *HttpClient { return NewDelete(url) }

// NewHttpClient 实例化：http客户端
//
//go:fix 推荐使用New方法
func NewHttpClient(urls ...string) *HttpClient {
	return &HttpClient{
		requestUrl:         str.APP.Buffer.NewString(operation.TernaryFuncAll(func() bool { return len(urls) == 0 }, func() string { return "" }, func() string { return urls[0] })).S(urls[1:]...).String(),
		requestQueries:     map[string]string{},
		requestHeaders:     map[string][]string{"Accept": {}, "Content-Type": {}},
		responseBody:       []byte{},
		responseBodyBuffer: bytes.NewBuffer([]byte{}),
		transport:          &http.Transport{
			// DisableKeepAlives:   true,             // 禁用连接复用
			// MaxIdleConns:        100,              // 最大空闲连接数
			// IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
			// TLSHandshakeTimeout: 10 * time.Second, // TLS握手超时时间
		},
	}
}

// NewGet 实例化：http客户端get请求
//
//go:fix 推荐使用NewGet方法
func NewGet(urls ...string) *HttpClient {
	return NewHttpClient(urls...).SetMethod(http.MethodGet)
}

// NewPost 实例化：http客户端post请求
//
//go:fix 推荐使用NewPost方法
func NewPost(urls ...string) *HttpClient {
	return NewHttpClient(urls...).SetMethod(http.MethodPost)
}

// NewPut 实例化：http客户端put请求
//
//go:fix 推荐使用NewPut方法
func NewPut(urls ...string) *HttpClient {
	return NewHttpClient(urls...).SetMethod(http.MethodPut)
}

// NewDelete 实例化：http客户端delete请求
//
//go:fix 推荐使用NewDelete方法
func NewDelete(urls ...string) *HttpClient {
	return NewHttpClient(urls...).SetMethod(http.MethodDelete)
}

// SetCert 设置SSL证书
func (my *HttpClient) SetCert(filename string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	var e error

	// 读取证书文件
	if my.cert, e = os.ReadFile(filename); e != nil {
		my.Err = e
	}

	return my
}

// SetUrl 设置请求地址
func (my *HttpClient) SetUrl(urls ...string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestUrl = str.APP.Buffer.NewString(operation.TernaryFuncAll(func() bool { return len(urls) == 0 }, func() string { return "" }, func() string { return urls[0] })).S(urls[1:]...).String()
	return my
}

// SetMethod 设置请求方法
func (my *HttpClient) SetMethod(method string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestMethod = method
	return my
}

// SetHeaders 设置请求头
func (my *HttpClient) SetHeaders(headers map[string][]string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestHeaders = headers
	return my
}

// AddHeaders 追加请求头
func (my *HttpClient) AddHeaders(headers map[string][]string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	for k, v := range headers {
		my.requestHeaders[k] = append(my.requestHeaders[k], v...)
	}

	return my
}

// SetQueries 设置请求参数
func (my *HttpClient) SetQueries(queries map[string]string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestQueries = queries
	return my
}

// SetAuthorization 设置认证
func (my *HttpClient) SetAuthorization(username, password, title string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestHeaders["Authorization"] = []string{str.APP.Buffer.NewString(title).S(" ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String()}
	return my
}

// SetBody 设置请求体
func (my *HttpClient) SetBody(body []byte) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.requestBody = body
	return my
}

// SetJsonBody 设置json请求体
func (my *HttpClient) SetJsonBody(body any) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeJson)
	my.requestBody, my.Err = json.Marshal(body)
	if my.Err != nil {
		my.Err = SetJsonBodyErr.Wrap(my.Err)
	}

	return my
}

// SetXmlBody 设置xml请求体
func (my *HttpClient) SetXmlBody(body any) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeXml)
	my.requestBody, my.Err = xml.Marshal(body)
	if my.Err != nil {
		my.Err = SetXmlBodyErr.Wrap(my.Err)
	}

	return my
}

// SetFormBody 设置表单请求体
func (my *HttpClient) SetFormBody(body map[string]string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeForm)
	params := url.Values{}
	for k, v := range body {
		params.Add(k, v)
	}
	my.requestBody = []byte(params.Encode())

	return my
}

// SetFormDataBody 设置表单数据请求体
func (my *HttpClient) SetFormDataBody(texts map[string]string, files map[string]string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	var (
		e      error
		buffer bytes.Buffer
	)

	my.setHeaderContentType("form-data")
	writer := multipart.NewWriter(&buffer)
	if len(texts) > 0 {
		for k, v := range texts {
			e = writer.WriteField(k, v)
			if e != nil {
				my.Err = SetFormBodyErr.Wrap(e)
				return my
			}
		}
	}

	if len(files) > 0 {
		for k, v := range files {
			fileWriter, _ := writer.CreateFormFile("fileField", k)
			file, e := os.Open(v)
			if e != nil {
				my.Err = SetFormBodyErr.Wrap(e)
				return my
			}
			_, e = io.Copy(fileWriter, file)
			if e != nil {
				my.Err = SetFormBodyErr.Wrap(e)
				return my
			}

			_ = file.Close()
		}
	}

	my.requestBody = []byte(writer.FormDataContentType())

	return my
}

// SetPlainBody 设置纯文本请求体
func (my *HttpClient) SetPlainBody(text string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypePlain)
	my.requestBody = []byte(text)

	return my
}

// SetHtmlBody 设置html请求体
func (my *HttpClient) SetHtmlBody(text string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeHtml)
	my.requestBody = []byte(text)

	return my
}

// SetCssBody 设置Css请求体
func (my *HttpClient) SetCssBody(text string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeCss)
	my.requestBody = []byte(text)

	return my
}

// SetJavascriptBody 设置Javascript请求体
func (my *HttpClient) SetJavascriptBody(text string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeJavascript)
	my.requestBody = []byte(text)

	return my
}

// SetSteamBodyByReader 设置字节码内容：通过readCloser接口
func (my *HttpClient) SetSteamBodyByReader(reader io.ReadCloser) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(ContentTypeSteam)

	if reader == nil {
		my.Err = SetSteamBodyErr.Panic()
		return my
	}

	// 创建RequestBodyReader用于读取文件内容
	if my.responseBodyBuffer.Len() > 1*1024*1024 { // 1MB
		_, my.Err = io.Copy(my.responseBodyBuffer, reader)
		if my.Err != nil {
			my.Err = ReadResponseErr.Wrap(my.Err)
			return my
		}
		my.requestBody = my.responseBodyBuffer.Bytes()
	} else {
		my.requestBody, my.Err = io.ReadAll(reader)
		if my.Err != nil {
			my.Err = ReadResponseErr.Wrap(my.Err)
			return my
		}
	}

	return my
}

// SetSteamBodyByFile 设置字节码内容：通过文件
func (my *HttpClient) SetSteamBodyByFile(filename string) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	var (
		err  error
		file *os.File
	)

	my.setHeaderContentType(ContentTypeSteam)

	file, err = os.Open(filename)
	if err != nil {
		my.Err = SetSteamBodyErr.Wrap(err)
		return my
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}(file)

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	// 创建RequestBodyReader用于读取文件内容
	if size > 1*1024*1024 {
		_, my.Err = io.Copy(my.responseBodyBuffer, file)
		if my.Err != nil {
			my.Err = ReadResponseErr.Wrap(my.Err)
			return my
		}
		my.requestBody = my.responseBodyBuffer.Bytes()
	} else {
		my.requestBody, err = io.ReadAll(file)
		if err != nil {
			my.Err = ReadResponseErr.Wrap(err)
			return my
		}
	}

	// my.request.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	return my
}

func (my *HttpClient) SetHeaderContentType(key ContentType) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderContentType(key)
	return my
}

// setHeaderContentType 设置请求头内容类型
func (my *HttpClient) setHeaderContentType(key ContentType) {
	if val, ok := ContentTypes[key]; ok {
		my.requestHeaders["Content-Type"] = []string{val}
	}
}

func (my *HttpClient) AppendHeaderContentType(keys ...ContentType) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.appendHeaderContentType(keys...)

	return my
}

// appendHeaderContentType 追加请求头内容类型
func (my *HttpClient) appendHeaderContentType(keys ...ContentType) {

	values := make([]string, len(keys))
	for k, v := range keys {
		if val, ok := ContentTypes[v]; ok {
			values[k] = val
		}
	}

	if len(my.requestHeaders["Content-Type"]) == 0 {
		my.requestHeaders["Content-Type"] = values
	} else {
		my.requestHeaders["Content-Type"] = append(my.requestHeaders["Content-Type"], values...)
	}
}

func (my *HttpClient) SetHeaderAccept(key Accept) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.setHeaderAccept(key)

	return my
}

// setHeaderAccept 设置请求头接受内容类型
func (my *HttpClient) setHeaderAccept(key Accept) *HttpClient {
	if val, ok := Accepts[key]; ok {
		my.requestHeaders["Accept"] = []string{val}
	}

	return my
}

func (my *HttpClient) AppendHeaderAccept(keys ...Accept) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.appendHeaderAccept(keys...)

	return my
}

// appendHeaderAccept 追加请求头接受内容类型
func (my *HttpClient) appendHeaderAccept(keys ...Accept) {
	values := make([]string, len(keys))
	for k, v := range keys {
		if val, ok := Accepts[v]; ok {
			values[k] = val
		}
	}

	if len(my.requestHeaders["Accept"]) == 0 {
		my.requestHeaders["Accept"] = values
	} else {
		my.requestHeaders["Accept"] = append(my.requestHeaders["Accept"], values...)
	}
}

// SetTimeoutSecond 设置超时：秒
func (my *HttpClient) SetTimeoutSecond(timeoutSecond int64) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.timeoutSecond = timeoutSecond

	return my
}

// SetTimeout 设置超时
func (my *HttpClient) SetTimeout(t time.Duration) *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	if t < 0 {
		my.timeoutSecond = 0
	} else {
		my.timeoutSecond = int64(t.Seconds())
	}

	return my
}

// GetTransport 获取自定义传输层
func (my *HttpClient) GetTransport() *http.Transport {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.transport
}

// SetTransport 设置自定义传输层
func (my *HttpClient) SetTransport(transport *http.Transport) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	my.transport = transport

	return my
}

func (my *HttpClient) GetResponse() *http.Response {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getResponse()
}

// getResponse 获取响应对象
func (my *HttpClient) getResponse() *http.Response { return my.response }

// ParseByContentType 根据响应头Content-Type自动解析响应体
func (my *HttpClient) ParseByContentType(target any) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	switch ContentType(my.getResponse().Header.Get("Content-Type")) {
	case ContentTypeJson:
		my.getResponseJsonBody(target)
	case ContentTypeXml:
		my.getResponseXmlBody(target)
	}

	return my
}

// GetResponseRawBody 获取原始响应体
func (my *HttpClient) GetResponseRawBody() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.responseBody
}

func (my *HttpClient) GetResponseJsonBody(target any, keys ...any) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getResponseJsonBody(target, keys...)
}

// getResponseJsonBody 获取json格式响应体
func (my *HttpClient) getResponseJsonBody(target any, keys ...any) *HttpClient {
	if my.responseBody == nil {
		return my
	}

	if len(my.responseBody) == 0 {
		return my
	}

	if len(keys) > 0 {
		jsonIter.Get(my.responseBody, keys...).ToVal(&target)
		return my
	} else {
		if e := json.Unmarshal(my.responseBody, &target); e != nil {
			my.Err = UnmarshalJsonErr.Wrap(e)
		}
	}

	return my
}

func (my *HttpClient) GetResponseXmlBody(target any) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getResponseXmlBody(target)
}

// getResponseXmlBody 获取xml格式响应体
func (my *HttpClient) getResponseXmlBody(target any) *HttpClient {
	if len(my.responseBody) == 0 {
		return my
	}

	if e := xml.Unmarshal(my.responseBody, &target); e != nil {
		my.Err = UnmarshalXmlErr.Wrap(e)
	}

	return my
}

// SaveResponseSteamFile 保存二进制到文件
//
//go:fix 建议使用Download方法
func (my *HttpClient) SaveResponseSteamFile(filename string) *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	if len(my.responseBody) == 0 {
		return my
	}

	// 创建一个新的文件
	file, err := os.Create(filename)
	if err != nil {
		my.Err = err
		return my
	}

	// 将二进制数据写入文件
	_, err = file.Write(my.responseBody)
	if err != nil {
		my.Err = err
		return my
	}

	my.Err = file.Close()

	return my
}

// GetRequest 获取请求
func (my *HttpClient) GetRequest() *http.Request {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.request
}

func (my *HttpClient) GenerateRequest() *HttpClient {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.generateRequest()
}

// generateRequest 生成请求对象
func (my *HttpClient) generateRequest() *HttpClient {
	var e error

	// 设置url参数
	my.setQueries()

	my.request, e = http.NewRequest(my.requestMethod, my.requestUrl, bytes.NewReader(my.requestBody))
	if e != nil {
		my.Err = GenerateRequestErr.Wrap(e)
		return my
	}

	// 设置请求头
	my.addHeaders()

	// 检查请求对象
	if my.Err = my.check(); my.Err != nil {
		return my
	}

	// 创建一个新的证书池，并将证书添加到池中
	if len(my.cert) > 0 {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(my.cert) {
			my.Err = GenerateCertErr.Panic()
			return my
		}

		// 创建一个新的Transport
		my.transport.TLSClientConfig = &tls.Config{RootCAs: certPool}
	}

	my.isReady = true

	return my
}

// beforeSend 发送请求前置动作
func (my *HttpClient) beforeSend() *http.Client {
	if !my.isReady {
		if my.generateRequest(); my.Err != nil {
			return nil
		}
	}

	client := &http.Client{}

	// 发送新的请求
	client.Transport = my.transport

	// 设置超时
	if my.timeoutSecond > 0 {
		client.Timeout = time.Duration(my.timeoutSecond) * time.Second
	}

	return client
}

// Download 使用下载器下载文件
func (my *HttpClient) Download(filename string) *HttpClientDownload {
	my.lock.Lock()
	defer my.lock.Unlock()

	return HttpClientDownloadApp.New(my, filename)
}

// Send 发送请求
func (my *HttpClient) Send() *HttpClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.responseBodyBuffer.Reset() // 重置响应体缓存
	my.responseBody = []byte{}    // 重置响应体

	client := my.beforeSend()
	if my.Err != nil {
		return my
	}

	my.request.Header.Set("Content-Length", fmt.Sprintf("%d", len(my.requestBody)))

	my.response, my.Err = client.Do(my.request)
	if my.Err != nil {
		return my
	}
	defer my.response.Body.Close()

	// 读取新的响应的主体
	if my.response.ContentLength > 1*1024*1024 { // 1MB
		if _, my.Err = io.Copy(my.responseBodyBuffer, my.response.Body); my.Err != nil {
			my.Err = ReadResponseErr.Wrap(my.Err)
			return my
		}
		my.responseBody = my.responseBodyBuffer.Bytes()
	} else {
		my.responseBody, my.Err = io.ReadAll(my.response.Body)
		if my.Err != nil {
			my.Err = ReadResponseErr.Wrap(my.Err)
			return my
		}
	}

	my.response.Body = io.NopCloser(bytes.NewBuffer(my.responseBody)) // 还原响应体

	my.isReady = false

	return my
}

// 检查条件是否满足
func (my *HttpClient) check() error {
	if my.requestUrl == "" {
		return UrlEmptyErr.Panic()
	}

	if my.requestMethod == "" {
		my.requestMethod = http.MethodGet
	}

	return nil
}

// 设置url参数
func (my *HttpClient) setQueries() {
	if len(my.requestQueries) > 0 {
		queries := url.Values{}
		for k, v := range my.requestQueries {
			queries.Add(k, v)
		}

		if len(queries) > 0 {
			my.requestUrl = my.requestUrl + "?" + queries.Encode()
		}
	}
}

// 设置请求头
func (my *HttpClient) addHeaders() {
	for k, v := range my.requestHeaders {
		my.request.Header[k] = append(my.request.Header[k], v...)
	}
}
