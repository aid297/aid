package httpClientV4

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// HTTPClientAttributer 接口 - 简化接口减少逃逸
type HTTPClientAttributer interface {
	Apply(client *HTTPClient)
}

// ============ URL 属性 ============

type AttrURL struct {
	url string
}

func URL(parts ...string) *AttrURL {
	if len(parts) == 0 {
		return &AttrURL{url: ""}
	}
	if len(parts) == 1 {
		return &AttrURL{url: parts[0]}
	}

	// 使用 strings.Builder 避免多次字符串拼接
	var sb strings.Builder
	totalLen := 0
	for _, p := range parts {
		totalLen += len(p)
	}
	sb.Grow(totalLen)

	for _, p := range parts {
		sb.WriteString(p)
	}

	return &AttrURL{url: sb.String()}
}

func (a *AttrURL) Apply(client *HTTPClient) {
	client.url = a.url
}

// ============ Queries 属性 ============

type AttrQueries struct {
	queries map[string]string
}

func Queries(queries map[string]string) *AttrQueries {
	if queries == nil {
		return &AttrQueries{queries: make(map[string]string)}
	}
	return &AttrQueries{queries: queries}
}

func (a *AttrQueries) Append(queries map[string]string) *AttrQueries {
	if len(queries) > 0 {
		for k, v := range queries {
			a.queries[k] = v
		}
	}
	return a
}

func (a *AttrQueries) AppendOne(key, value string) *AttrQueries {
	a.queries[key] = value
	return a
}

func (a *AttrQueries) Apply(client *HTTPClient) {
	for k, v := range a.queries {
		client.queries[k] = v
	}
}

// ============ Method 属性 ============

type AttrMethod struct {
	method string
}

func Method(method string) *AttrMethod {
	return &AttrMethod{method: method}
}

func (a *AttrMethod) Apply(client *HTTPClient) {
	client.method = a.method
}

// ============ Header 属性 ============

type AttrHeader struct {
	header http.Header
	append bool
}

func AppendHeader(key string, values ...string) *AttrHeader {
	h := make(http.Header)
	h[key] = values
	return &AttrHeader{header: h, append: true}
}

func SetHeader(key string, values ...string) *AttrHeader {
	h := make(http.Header)
	h[key] = values
	return &AttrHeader{header: h, append: false}
}

func AppendHeaders(headers http.Header) *AttrHeader {
	return &AttrHeader{header: headers, append: true}
}

func SetHeaders(headers http.Header) *AttrHeader {
	return &AttrHeader{header: headers, append: false}
}

func (a *AttrHeader) Apply(client *HTTPClient) {
	if a.append {
		for k, v := range a.header {
			client.headers[k] = append(client.headers[k], v...)
		}
	} else {
		for k, v := range a.header {
			client.headers[k] = v
		}
	}
}

// ============ ContentType 和 Accept ============

type AttrContentType struct {
	contentType ContentType
}

func ContentType_(ct ContentType) *AttrContentType {
	return &AttrContentType{contentType: ct}
}

func (a *AttrContentType) Apply(client *HTTPClient) {
	client.headers.Set("Content-Type", string(a.contentType))
}

type AttrAccept struct {
	accept Accept
}

func Accept_(accept Accept) *AttrAccept {
	return &AttrAccept{accept: accept}
}

func (a *AttrAccept) Apply(client *HTTPClient) {
	client.headers.Set("Accept", string(a.accept))
}

// ============ Authorization ============

type AttrAuthorization struct {
	username string
	password string
	title    string
}

func Authorization(username, password, title string) *AttrAuthorization {
	return &AttrAuthorization{
		username: username,
		password: password,
		title:    title,
	}
}

func (a *AttrAuthorization) Apply(client *HTTPClient) {
	credentials := fmt.Sprintf("%s:%s", a.username, a.password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	client.headers.Set("Authorization", a.title+" "+encoded)
}

// ============ Body 属性 ============

type AttrBody struct {
	err         error
	body        []byte
	contentType ContentType
}

func JSON(body any) *AttrBody {
	ins := &AttrBody{}
	ins.body, ins.err = json.Marshal(body)
	ins.contentType = ContentTypeJSON
	return ins
}

func XML(body any) *AttrBody {
	ins := &AttrBody{}
	ins.body, ins.err = xml.Marshal(body)
	ins.contentType = ContentTypeXML
	return ins
}

func Form(body map[string]string) *AttrBody {
	ins := &AttrBody{}
	params := url.Values{}
	for k, v := range body {
		params.Add(k, v)
	}
	ins.body = []byte(params.Encode())
	ins.contentType = ContentTypeXWwwFormURLencoded
	return ins
}

func FormData(fields, files map[string]string) *AttrBody {
	ins := &AttrBody{}
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	writer := multipart.NewWriter(buf)

	// 写入字段
	if len(fields) > 0 {
		for k, v := range fields {
			if err := writer.WriteField(k, v); err != nil {
				ins.err = err
				return ins
			}
		}
	}

	// 写入文件
	if len(files) > 0 {
		for k, v := range files {
			fileWriter, _ := writer.CreateFormFile("file", k)
			file, err := os.Open(v)
			if err != nil {
				ins.err = err
				return ins
			}
			_, err = io.Copy(fileWriter, file)
			_ = file.Close()
			if err != nil {
				ins.err = err
				return ins
			}
		}
	}

	_ = writer.Close()

	// 复制数据(因为 buffer 会被归还)
	ins.body = make([]byte, buf.Len())
	copy(ins.body, buf.Bytes())
	ins.contentType = ContentType(writer.FormDataContentType())

	return ins
}

func Plain(body string) *AttrBody {
	return &AttrBody{
		body:        []byte(body),
		contentType: ContentTypePlain,
	}
}

func HTML(body string) *AttrBody {
	return &AttrBody{
		body:        []byte(body),
		contentType: ContentTypeHTML,
	}
}

func CSS(body string) *AttrBody {
	return &AttrBody{
		body:        []byte(body),
		contentType: ContentTypeCSS,
	}
}

func Javascript(body string) *AttrBody {
	return &AttrBody{
		body:        []byte(body),
		contentType: ContentTypeJavascript,
	}
}

func Bytes(body []byte) *AttrBody {
	return &AttrBody{body: body}
}

func Reader(body io.ReadCloser) *AttrBody {
	ins := &AttrBody{}
	if body == nil {
		ins.err = errors.New("设置steam流失败：不能为空")
		return ins
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if _, ins.err = io.Copy(buf, body); ins.err != nil {
		return ins
	}

	// 复制数据
	ins.body = make([]byte, buf.Len())
	copy(ins.body, buf.Bytes())

	return ins
}

func File(filename string) *AttrBody {
	ins := &AttrBody{}

	file, err := os.Open(filename)
	if err != nil {
		ins.err = err
		return ins
	}
	defer file.Close()

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	if size > 1*1024*1024 { // 1MB
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)

		if _, ins.err = io.Copy(buf, file); ins.err != nil {
			return ins
		}
		ins.body = make([]byte, buf.Len())
		copy(ins.body, buf.Bytes())
	} else {
		if ins.body, ins.err = io.ReadAll(file); ins.err != nil {
			return ins
		}
	}

	return ins
}

func (a *AttrBody) Apply(client *HTTPClient) {
	client.requestBody = a.body
	if a.contentType != "" {
		client.headers.Set("Content-Type", string(a.contentType))
	}
	client.err = a.err
}

// ============ Timeout 属性 ============

type AttrTimeout struct {
	timeout time.Duration
}

func Timeout(timeout time.Duration) *AttrTimeout {
	if timeout <= 0 {
		timeout = 0
	}
	return &AttrTimeout{timeout: timeout}
}

func (a *AttrTimeout) Apply(client *HTTPClient) {
	client.timeout = a.timeout
}

// ============ Transport 属性 ============

type AttrTransport struct {
	transport *http.Transport
}

func Transport(transport *http.Transport) *AttrTransport {
	return &AttrTransport{transport: transport}
}

func (a *AttrTransport) Apply(client *HTTPClient) {
	client.transport = a.transport
}

func TransportDefault() *AttrTransport {
	return &AttrTransport{transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}}
}

// ============ Cert 属性 ============

type AttrCert struct {
	cert []byte
}

func Cert(cert []byte) *AttrCert {
	return &AttrCert{cert: cert}
}

func (a *AttrCert) Apply(client *HTTPClient) {
	client.cert = a.cert
}

// ============ AutoCopy 属性 ============

type AttrAutoCopy struct {
	autoCopy bool
}

func AutoCopy(autoCopy bool) *AttrAutoCopy {
	return &AttrAutoCopy{autoCopy: autoCopy}
}

func (a *AttrAutoCopy) Apply(client *HTTPClient) {
	client.autoCopy = a.autoCopy
}
