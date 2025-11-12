package httpClientV2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cast"

	"github.com/aid297/aid/str"
)

type (
	HttpClientAttributer interface {
		Register(req *HttpClient)
		Error() error
		ImplHttpClientAttributer()
	}

	AttrURL                struct{ url string }
	AttrQueries            struct{ queries map[string]any }
	AttrMethod             struct{ method string }
	AttrAppendHeaderValues struct{ headers map[string][]string }
	AttrAppendHeaderValue  struct{ headers map[string]string }
	AttrSetHeaderValues    struct{ headers map[string][]string }
	AttrSetHeaderValue     struct{ headers map[string]string }
	AttrBody               struct {
		err         error
		body        []byte
		contentType ContentType
	}
	AttrTimeout          struct{ timeout time.Duration }
	AttrTransport        struct{ transport *http.Transport }
	AttrTransportDefault struct{ transport *http.Transport }
	AttrCert             struct{ cert []byte }
	AttrAutoCopyResBody  struct{ autoCopy bool }
)

func URL(urls ...string) HttpClientAttributer {
	ins := &AttrURL{url: ""}
	switch {

	}
	if len(urls) == 0 {
	} else if len(urls) == 1 {
		ins.url = urls[0]
	} else {
		ins.url = str.APP.Buffer.JoinString(urls...)
	}

	return ins
}

func (my *AttrURL) Register(req *HttpClient) { req.url = my.url }

func (my *AttrURL) Error() error { return nil }

func (*AttrURL) ImplHttpClientAttributer() {}

func Queries(queries map[string]any) *AttrQueries {
	if queries == nil {
		queries = map[string]any{}
	}

	return &AttrQueries{queries: queries}
}

func (my *AttrQueries) Append(queries map[string]any) *AttrQueries {
	if len(queries) > 0 {
		maps.Copy(my.queries, queries)
	}

	return my
}

func (my *AttrQueries) AppendOne(key string, value any) *AttrQueries {
	my.queries[key] = value
	return my
}

func (my *AttrQueries) Register(req *HttpClient) { req.queries = my.queries }

func (my *AttrQueries) Error() error { return nil }

func (*AttrQueries) ImplHttpClientAttributer() {}

func Method(method string) *AttrMethod {
	return &AttrMethod{method: method}
}

func (my *AttrMethod) Register(req *HttpClient) { req.method = my.method }

func (my *AttrMethod) Error() error { return nil }

func (*AttrMethod) ImplHttpClientAttributer() {}

func AppendHeaderValue(headers map[string]string) *AttrAppendHeaderValue {
	if headers == nil {
		headers = map[string]string{}
	}

	return &AttrAppendHeaderValue{headers: headers}
}

func (my *AttrAppendHeaderValue) Append(headers map[string]string) *AttrAppendHeaderValue {
	if len(headers) > 0 {
		maps.Copy(my.headers, headers)
	}

	return my
}

func (my *AttrAppendHeaderValue) AppendOne(key string, value string) *AttrAppendHeaderValue {
	my.headers[key] = value
	return my
}

func (my *AttrAppendHeaderValue) ContentType(contentType ContentType) *AttrAppendHeaderValue {
	my.headers["Content-Type"] = ContentTypes[contentType]
	return my
}

func (my *AttrAppendHeaderValue) Accept(accept Accept) *AttrAppendHeaderValue {
	my.headers["Accept"] = Accepts[accept]
	return my
}

func (my *AttrAppendHeaderValue) Authorization(username, password, title string) *AttrAppendHeaderValue {
	my.headers["Authorization"] = str.BufferApp.NewString(title, " ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String()
	return my
}

func (my *AttrAppendHeaderValue) Register(req *HttpClient) {
	if req.headers == nil {
		req.headers = map[string][]string{}
	} else {
		for key, values := range my.headers {
			if _, exists := req.headers[key]; !exists {
				req.headers[key] = []string{values}
			} else {
				req.headers[key] = append(req.headers[key], []string{values}...)
			}
		}
	}
}

func (my *AttrAppendHeaderValue) Error() error { return nil }

func (*AttrAppendHeaderValue) ImplHttpClientAttributer() {}

func AppendHeaderValues(headers map[string][]string) *AttrAppendHeaderValues {
	if headers == nil {
		headers = map[string][]string{}
	}

	return &AttrAppendHeaderValues{headers: headers}
}

func (my *AttrAppendHeaderValues) Append(headers map[string][]string) *AttrAppendHeaderValues {
	if len(headers) > 0 {
		maps.Copy(my.headers, headers)
	}

	return my
}

func (my *AttrAppendHeaderValues) AppendOne(key string, values ...string) *AttrAppendHeaderValues {
	my.headers[key] = values
	return my
}

func (my *AttrAppendHeaderValues) ContentType(contentType ContentType) *AttrAppendHeaderValues {
	my.headers["Content-Type"] = []string{ContentTypes[contentType]}
	return my
}

func (my *AttrAppendHeaderValues) Accept(accept Accept) *AttrAppendHeaderValues {
	my.headers["Accept"] = []string{Accepts[accept]}
	return my
}

func (my *AttrAppendHeaderValues) Authorization(username, password, title string) *AttrAppendHeaderValues {
	my.headers["Authorization"] = []string{str.BufferApp.NewString(title, " ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String()}
	return my
}

func (my *AttrAppendHeaderValues) Register(req *HttpClient) {
	if req.headers == nil {
		req.headers = my.headers
	} else {
		for key, values := range my.headers {
			if _, exists := req.headers[key]; !exists {
				req.headers[key] = values
			} else {
				req.headers[key] = append(req.headers[key], values...)
			}
		}
	}
}

func (my *AttrAppendHeaderValues) Error() error { return nil }

func (*AttrAppendHeaderValues) ImplHttpClientAttributer() {}

func SetHeaderValue(headers map[string]string) *AttrSetHeaderValue {
	if headers == nil {
		headers = map[string]string{}
	}

	return &AttrSetHeaderValue{headers: headers}
}

func (my *AttrSetHeaderValue) ContentType(contentType ContentType) *AttrSetHeaderValue {
	my.headers["Content-Type"] = ContentTypes[contentType]
	return my
}

func (my *AttrSetHeaderValue) Accept(accept Accept) *AttrSetHeaderValue {
	my.headers["Accept"] = Accepts[accept]
	return my
}

func (my *AttrSetHeaderValue) Authorization(username, password, title string) *AttrSetHeaderValue {
	my.headers["Authorization"] = str.BufferApp.NewString(title, " ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String()
	return my
}

func (my *AttrSetHeaderValue) Register(req *HttpClient) {
	if req.headers == nil {
		req.headers = map[string][]string{}
	} else {
		for idx := range my.headers {
			req.headers[idx] = []string{my.headers[idx]}
		}
	}
}

func (my *AttrSetHeaderValue) Error() error { return nil }

func (*AttrSetHeaderValue) ImplHttpClientAttributer() {}

func SetHeaderValues(headers map[string][]string) *AttrSetHeaderValues {
	if headers == nil {
		headers = map[string][]string{}
	}

	return &AttrSetHeaderValues{headers: headers}
}

func (my *AttrSetHeaderValues) ContentType(contentType ContentType) *AttrSetHeaderValues {
	my.headers["Content-Type"] = []string{ContentTypes[contentType]}
	return my
}

func (my *AttrSetHeaderValues) Accept(accept Accept) *AttrSetHeaderValues {
	my.headers["Accept"] = []string{Accepts[accept]}
	return my
}

func (my *AttrSetHeaderValues) Authorization(username, password, title string) *AttrSetHeaderValues {
	my.headers["Authorization"] = []string{str.BufferApp.NewString(title, " ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String()}
	return my
}

func (my *AttrSetHeaderValues) Register(req *HttpClient) {
	if req.headers == nil {
		req.headers = my.headers
	} else {
		maps.Copy(req.headers, my.headers)
	}
}

func (my *AttrSetHeaderValues) Error() error { return nil }

func (*AttrSetHeaderValues) ImplHttpClientAttributer() {}

func JSON(body any) *AttrBody {
	ins := &AttrBody{}
	ins.body, ins.err = json.Marshal(body)
	ins.contentType = ContentTypeJson

	return ins
}

func XML(body any) HttpClientAttributer {
	ins := &AttrBody{}
	ins.body, ins.err = xml.Marshal(body)
	ins.contentType = ContentTypeXml

	return ins
}

func Form(body map[string]any) *AttrBody {
	ins := &AttrBody{}
	params := url.Values{}
	for k, v := range body {
		params.Add(k, cast.ToString(v))
	}
	ins.body = []byte(params.Encode())
	ins.contentType = ContentTypeXWwwFormUrlencoded

	return ins
}

func FormData(fields, files map[string]string) *AttrBody {
	var (
		e      error
		buffer bytes.Buffer
		ins    = &AttrBody{}
	)

	writer := multipart.NewWriter(&buffer)
	if len(fields) > 0 {
		for k, v := range fields {
			if e = writer.WriteField(k, v); e != nil {
				ins.err = e
				return ins
			}
		}
	}

	if len(files) > 0 {
		for k, v := range files {
			fileWriter, _ := writer.CreateFormFile("file", k)
			file, e := os.Open(v)
			if e != nil {
				ins.err = e
				return ins
			}
			_, e = io.Copy(fileWriter, file)
			if e != nil {
				ins.err = e
				return ins
			}

			_ = file.Close()
		}
	}

	ins.body = []byte(writer.FormDataContentType())
	ins.contentType = ContentTypeFormData

	return ins
}

func Plain(body string) *AttrBody {
	ins := &AttrBody{}
	ins.body = []byte(body)
	ins.contentType = ContentTypePlain

	return ins
}

func HTML(body string) *AttrBody {
	ins := &AttrBody{}
	ins.body = []byte(body)
	ins.contentType = ContentTypeXml

	return ins
}

func CSS(body string) *AttrBody {
	ins := &AttrBody{}
	ins.body = []byte(body)
	ins.contentType = ContentTypeCss

	return ins
}

func Javascript(body string) *AttrBody {
	ins := &AttrBody{}
	ins.body = []byte(body)
	ins.contentType = ContentTypeJavascript

	return ins
}

func Bytes(body []byte) *AttrBody {
	ins := &AttrBody{body: body}

	return ins
}

func Reader(body io.ReadCloser) *AttrBody {
	var (
		ins    = &AttrBody{}
		buffer = bytes.NewBuffer([]byte{})
	)
	if body == nil {
		ins.err = errors.New("设置steam流失败：不能为空")
		return ins
	}

	if _, ins.err = io.Copy(buffer, body); ins.err != nil {
		return ins
	}
	ins.body = buffer.Bytes()

	return ins
}

func File(filename string) *AttrBody {
	var (
		ins    = &AttrBody{}
		file   *os.File
		buffer = bytes.NewBuffer([]byte{})
	)

	if file, ins.err = os.Open(filename); ins.err != nil {
		return ins
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
		if _, ins.err = io.Copy(buffer, file); ins.err != nil {
			return ins
		}
		ins.body = buffer.Bytes()
	} else {
		if ins.body, ins.err = io.ReadAll(file); ins.err != nil {
			return ins
		}
	}

	return ins
}

func (my *AttrBody) Register(req *HttpClient) {
	req.requestBody = my.body
	if my.contentType != "" {
		req.headers["Content-Type"] = []string{ContentTypes[my.contentType]}
	}
	req.err = my.err
}

func (my *AttrBody) Error() error { return my.err }

func (*AttrBody) ImplHttpClientAttributer() {}

func Timeout(timeout time.Duration) *AttrTimeout {
	if timeout <= 0 {
		timeout = 0
	}

	return &AttrTimeout{timeout: timeout}
}

func (my *AttrTimeout) Register(req *HttpClient) { req.timeout = my.timeout }

func (*AttrTimeout) Error() error { return nil }

func (*AttrTimeout) ImplHttpClientAttributer() {}

func Transport(transport *http.Transport) *AttrTransport {
	return &AttrTransport{transport: transport}
}

func (my *AttrTransport) Register(req *HttpClient) { req.transport = my.transport }

func (my *AttrTransport) Error() error { return nil }

func (*AttrTransport) ImplHttpClientAttributer() {}

func TransportDefault() *AttrTransportDefault {
	return &AttrTransportDefault{transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}}
}

func (my *AttrTransportDefault) Register(req *HttpClient) { req.transport = my.transport }

func (my *AttrTransportDefault) Error() error { return nil }

func (*AttrTransportDefault) ImplHttpClientAttributer() {}

func Cert(cert []byte) *AttrCert { return &AttrCert{cert: cert} }

func (my *AttrCert) Register(req *HttpClient) { req.cert = my.cert }

func (my *AttrCert) Error() error { return nil }

func (*AttrCert) ImplHttpClientAttributer() {}

func AutoCopy(autoCopy bool) *AttrAutoCopyResBody { return &AttrAutoCopyResBody{autoCopy: autoCopy} }

func (my *AttrAutoCopyResBody) Register(req *HttpClient) { req.autoCopy = my.autoCopy }

func (*AttrAutoCopyResBody) Error() error { return nil }

func (*AttrAutoCopyResBody) ImplHttpClientAttributer() {}
