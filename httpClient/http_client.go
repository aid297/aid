package httpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/spf13/cast"

	"github.com/aid297/aid/v2/anyMap"
	"github.com/aid297/aid/v2/compressions"
	"github.com/aid297/aid/v2/consts/volumeInfo"
	"github.com/aid297/aid/v2/debugLogger"
	"github.com/aid297/aid/v2/operation"
	"github.com/aid297/aid/v2/secret"
	"github.com/aid297/aid/v2/str"
)

var (
	_             HTTPClient = (*HTTPClientImpl)(nil)
	fileSplitSize int64      = 10 * volumeInfo.MB
)

type (
	// chunkedFileReader 切块文件读取器，用于大文件并行上传
	chunkedFileReader struct {
		file          *os.File
		fileSize      int64
		chunkSize     int64
		numChunks     int
		offset        int64
		lock          sync.Mutex
		encryptor     secret.Symmetric
		compressor    compressions.Compressor
		maxGoroutines uint64
	}

	HTTPClient interface {
		init(method string, attrs ...HTTPClientAttr) (HTTPClient, error)
		setAttrs(attrs ...HTTPClientAttr) (err error)
		SetAttrs(attrs ...HTTPClientAttr) (HTTPClient, error)
		setURL(urls ...any) (err error)
		setQueries(queries map[string]any) (err error)
		setQueriesNotEmpty(queries map[string]any) (err error)
		setQuery(key string, value any) (err error)
		setMethod(method string) (err error)
		setHeaders(headers map[string][]any) (err error)
		setHeader(key string, value any) (err error)
		appendHeaders(headers map[string][]any) (err error)
		appendHeader(key string, value any) (err error)
		setBody(body *bytes.Buffer) (err error)
		setBodyJSON(body any) (err error)
		setBodyXML(body any) (err error)
		setBodyForm(body map[string]any) (err error)
		setBodyFormData(fields, files map[string]string) (err error)
		setBodyPlain(body string) (err error)
		setBodyHTML(body string) (err error)
		setBodyCSS(body string) (err error)
		setBodyJavascript(body string) (err error)
		setBodyBytes(body []byte) (err error)
		setBodyReadCloser(body io.ReadCloser) (err error)
		setBodyFile(filename string, goroutineCount uint64) (err error)
		setRateLimit(rate uint64) (err error)
		setTimeout(timeout time.Duration) (err error)
		setTransport(transport *http.Transport) (err error)
		setCert(cert []byte) (err error)
		setAutoCopy(autoCopy bool) (err error)
		setEncryptor(symmetricEncryptor secret.Symmetric)
		setCompressor(compressor compressions.Compressor)

		GetURL() string
		getURL() string
		GetQueries() map[string]any
		getQueries() map[string]any
		GetMethod() string
		getMethod() string
		GetHeaders() map[string][]any
		getHeaders() map[string][]any
		GetBody() []byte
		getBody() []byte
		GetTimeout() time.Duration
		getTimeout() time.Duration
		GetTransport() *http.Transport
		getTransport() *http.Transport
		GetCert() []byte
		getCert() []byte
		GetRawRequest() *http.Request
		getRawRequest() *http.Request
		GetRawResponse() *http.Response
		getRawResponse() *http.Response
		GetClient() *http.Client
		getClient() *http.Client
		send() HTTPClient
		OK() error
		isNeedRetry(condition func(statusCode int, err error) bool) (needRetry bool)
		SendWithRetry(count uint, interval time.Duration, condition func(statusCode int, err error) bool) (HTTPClient, []error)
		Send() HTTPClient
		parseBody()
		ToJSON(target any, keys ...any) HTTPClient
		ToXML(target any) HTTPClient
		ToBytes() []byte
		ToWriter(writer http.ResponseWriter) HTTPClient
		Error() error
		GetStatusCode() int
		GetStatus() string
	}

	HTTPClientImpl struct {
		err                error
		url                string
		queries            map[string]any
		method             string
		headers            map[string][]any
		requestBody        io.Reader
		requestBodyBuffer  *bytes.Buffer
		rateLimit          uint64
		responseBody       []byte
		timeout            time.Duration
		transport          *http.Transport
		cert               []byte
		rawRequest         *http.Request
		rawResponse        *http.Response
		client             *http.Client
		autoCopy           bool
		lock               sync.RWMutex
		symmetricEncryptor secret.Symmetric
		compressor         compressions.Compressor
		filePath           string
		fileSize           int64
		isChunked          bool
		chunkedReader      *chunkedFileReader
	}
)

func (*HTTPClientImpl) init(method string, attrs ...HTTPClientAttr) (HTTPClient, error) {
	buffer := bytes.NewBuffer([]byte{})
	ins := &HTTPClientImpl{method: method, requestBody: buffer, requestBodyBuffer: buffer, headers: map[string][]any{}}
	if err := ins.setAttrs(TransportDefault()); err != nil {
		return nil, err
	}
	if err := ins.setAttrs(attrs...); err != nil {
		return nil, err
	}
	return ins, nil
}

// SetDefaultFileSplitSize 设置默认文件切块大小
func SetDefaultFileSplitSize(chunkSize int64) { fileSplitSize = chunkSize }

func New(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodGet, attrs...)
}

func GET(attrs ...HTTPClientAttr) (HTTPClient, error) { return New(attrs...) }

func POST(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodPost, attrs...)
}

func PUT(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodPut, attrs...)
}

func PATCH(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodPatch, attrs...)
}

func DELETE(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodDelete, attrs...)
}

func HEAD(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodHead, attrs...)
}

func OPTIONS(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodOptions, attrs...)
}

func TRACE(attrs ...HTTPClientAttr) (HTTPClient, error) {
	return new(HTTPClientImpl).init(http.MethodTrace, attrs...)
}

func (my *HTTPClientImpl) setAttrs(attrs ...HTTPClientAttr) (err error) {
	if len(attrs) > 0 {
		for _, option := range attrs {
			if err = option(my); err != nil {
				return err
			}
		}
	}
	return
}

func (my *HTTPClientImpl) SetAttrs(attrs ...HTTPClientAttr) (HTTPClient, error) {
	my.lock.Lock()
	defer my.lock.Unlock()

	err := my.setAttrs(attrs...)

	return my, err
}

func (my *HTTPClientImpl) setURL(urls ...any) (err error) {
	if len(urls) == 0 {
	} else if len(urls) == 1 {
		my.url = cast.ToString(urls[0])
	} else {
		my.url = str.APP.Buffer.JoinString(cast.ToStringSlice(urls)...)
	}
	return
}

func (my *HTTPClientImpl) setQueries(queries map[string]any) (err error) {
	if queries == nil {
		queries = map[string]any{}
	}

	maps.Copy(my.queries, queries)
	return
}

func (my *HTTPClientImpl) setQueriesNotEmpty(queries map[string]any) (err error) {
	if queries == nil {
		queries = map[string]any{}
	}

	maps.Copy(my.queries, anyMap.New(anyMap.Map(queries)).RemoveEmpty().ToMap())
	return
}

func (my *HTTPClientImpl) setQuery(key string, value any) (err error) {
	my.queries[key] = value
	return
}

func (my *HTTPClientImpl) setMethod(method string) (err error) { my.method = method; return }

func (my *HTTPClientImpl) setHeaders(headers map[string][]any) (err error) {
	my.headers = headers
	return
}

func (my *HTTPClientImpl) setHeader(key string, value any) (err error) {
	my.headers[key] = []any{value}
	return
}

func (my *HTTPClientImpl) appendHeaders(headers map[string][]any) (err error) {
	if len(headers) > 0 {
		for key := range headers {
			if _, exists := my.headers[key]; !exists {
				my.headers[key] = headers[key]
			} else {
				my.headers[key] = append(my.headers[key], headers[key]...)
			}
		}
	}

	return
}

func (my *HTTPClientImpl) appendHeader(key string, value any) (err error) {
	if _, exists := my.headers[key]; !exists {
		my.headers[key] = []any{value}
	} else {
		my.headers[key] = append(my.headers[key], value)
	}

	return
}

func (my *HTTPClientImpl) setBody(body *bytes.Buffer) (err error) {

	my.requestBody = bytes.NewBuffer(nil)
	my.requestBodyBuffer = bytes.NewBuffer(nil)

	if body != nil {
		my.requestBody = body
		my.requestBodyBuffer = body
	}

	return
}

func (my *HTTPClientImpl) setBodyJSON(body any) (err error) {
	var bodies []byte
	if bodies, err = json.Marshal(body); err != nil {
		return
	}
	my.setBody(bytes.NewBuffer(bodies))
	my.setHeader("Content-Type", ContentTypeJSON)

	return
}

func (my *HTTPClientImpl) setBodyXML(body any) (err error) {
	var bodies []byte
	if bodies, err = xml.Marshal(body); err != nil {
		return
	}
	my.setBody(bytes.NewBuffer(bodies))
	my.setHeader("Content-Type", ContentTypeXML)

	return
}

func (my *HTTPClientImpl) setBodyForm(body map[string]any) (err error) {
	var bodies []byte
	params := url.Values{}
	for k, v := range body {
		params.Add(k, cast.ToString(v))
	}
	bodies = []byte(params.Encode())
	my.setBody(bytes.NewBuffer(bodies))
	my.setHeader("Content-Type", ContentTypeXWWWFormURLEncoded)

	return
}

func (my *HTTPClientImpl) setBodyFormData(fields, files map[string]string) (err error) {
	var buffer bytes.Buffer

	writer := multipart.NewWriter(&buffer)
	if len(fields) > 0 {
		for k, v := range fields {
			if err = writer.WriteField(k, v); err != nil {
				return
			}
		}
	}

	if len(files) > 0 {
		for k, v := range files {
			var file *os.File
			fileWriter, _ := writer.CreateFormFile("file", k)
			if file, err = os.Open(v); err != nil {
				return
			}
			if _, err = io.Copy(fileWriter, file); err != nil {
				return
			}

			_ = file.Close()
		}
	}

	writer.Close()
	my.setBody(&buffer)
	my.headers["Content-Type"] = []any{writer.FormDataContentType()}

	return
}

func (my *HTTPClientImpl) setBodyPlain(body string) (err error) {
	my.setBody(bytes.NewBuffer([]byte(body)))
	my.setHeader("Content-Type", ContentTypePlain)

	return
}

func (my *HTTPClientImpl) setBodyHTML(body string) (err error) {
	my.setBody(bytes.NewBuffer([]byte(body)))
	my.setHeader("Content-Type", ContentTypeHTML)

	return
}

func (my *HTTPClientImpl) setBodyCSS(body string) (err error) {
	my.setBody(bytes.NewBuffer([]byte(body)))
	my.setHeader("Content-Type", ContentTypeCSS)

	return
}

func (my *HTTPClientImpl) setBodyJavascript(body string) (err error) {
	my.setBody(bytes.NewBuffer([]byte(body)))

	return
}

func (my *HTTPClientImpl) setBodyBytes(body []byte) (err error) {
	my.setBody(bytes.NewBuffer(body))
	return
}

func (my *HTTPClientImpl) setBodyReadCloser(body io.ReadCloser) (err error) {
	var a []byte

	if _, err = body.Read(a); err == nil {
		my.setBody(bytes.NewBuffer(a))
		return
	}

	return
}

func (my *HTTPClientImpl) setBodyFile(filename string, goroutineCount uint64) (err error) {
	var file *os.File

	if file, err = os.Open(filename); err != nil {
		return
	}

	// 获取文件大小
	stat, _ := file.Stat()
	size := stat.Size()

	if goroutineCount > 1 && size > fileSplitSize {
		// 大文件：设置切块读取
		my.filePath = filename
		my.fileSize = size
		my.isChunked = true

		chunkSize := max(size/int64(goroutineCount), 1*volumeInfo.MB) // 每个chunk至少1MB

		my.chunkedReader = &chunkedFileReader{
			file:          file,
			fileSize:      size,
			chunkSize:     chunkSize,
			numChunks:     int((size + chunkSize - 1) / chunkSize),
			offset:        0,
			encryptor:     my.symmetricEncryptor,
			compressor:    my.compressor,
			maxGoroutines: goroutineCount,
		}
	} else {
		// 小文件或不使用切块：直接读取到buffer
		fileBodis, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			return err
		}
		if err = file.Close(); err != nil {
			debugLogger.Error("[HTTP Client] 关闭文件错误: %v", err)
		}
		my.setBody(bytes.NewBuffer(fileBodis))
	}

	return
}

func (my *HTTPClientImpl) setRateLimit(rate uint64) (err error) { my.rateLimit = max(rate, 0); return }

func (my *HTTPClientImpl) setTimeout(timeout time.Duration) (err error) {
	my.timeout = max(timeout, 0)
	return
}

func (my *HTTPClientImpl) setTransport(transport *http.Transport) (err error) {
	my.transport = transport
	return
}

func (my *HTTPClientImpl) setCert(cert []byte) (err error) { my.cert = cert; return }

func (my *HTTPClientImpl) setAutoCopy(autoCopy bool) (err error) { my.autoCopy = autoCopy; return }

func (my *HTTPClientImpl) setEncryptor(symmetricEncryptor secret.Symmetric) {
	my.symmetricEncryptor = symmetricEncryptor
}

func (my *HTTPClientImpl) setCompressor(compressor compressions.Compressor) {
	my.compressor = compressor
}

func (my *HTTPClientImpl) GetURL() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getURL()
}

func (my *HTTPClientImpl) getURL() string {
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

func (my *HTTPClientImpl) GetQueries() map[string]any {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getQueries()
}

func (my *HTTPClientImpl) getQueries() map[string]any { return my.queries }

func (my *HTTPClientImpl) GetMethod() string {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getMethod()
}

func (my *HTTPClientImpl) getMethod() string { return my.method }

func (my *HTTPClientImpl) GetHeaders() map[string][]any {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getHeaders()
}

func (my *HTTPClientImpl) getHeaders() map[string][]any { return my.headers }

func (my *HTTPClientImpl) GetBody() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getBody()
}

func (my *HTTPClientImpl) getBody() []byte {
	if my.requestBodyBuffer != nil {
		return my.requestBodyBuffer.Bytes()
	}
	if my.requestBody != nil {
		b, _ := io.ReadAll(my.requestBody)
		return b
	}
	return nil
}

func (my *HTTPClientImpl) GetTimeout() time.Duration {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTimeout()
}

func (my *HTTPClientImpl) getTimeout() time.Duration { return my.timeout }

func (my *HTTPClientImpl) GetTransport() *http.Transport {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getTransport()
}

func (my *HTTPClientImpl) getTransport() *http.Transport { return my.transport }

func (my *HTTPClientImpl) GetCert() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getCert()
}

func (my *HTTPClientImpl) getCert() []byte { return my.cert }

func (my *HTTPClientImpl) GetRawRequest() *http.Request {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawRequest()
}

func (my *HTTPClientImpl) getRawRequest() *http.Request { return my.rawRequest }

func (my *HTTPClientImpl) GetRawResponse() *http.Response {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getRawResponse()
}

func (my *HTTPClientImpl) getRawResponse() *http.Response { return my.rawResponse }

func (my *HTTPClientImpl) GetClient() *http.Client {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.getClient()
}

func (my *HTTPClientImpl) getClient() *http.Client { return my.client }

func (my *HTTPClientImpl) send() HTTPClient {
	var (
		cipherText  []byte
		compressed  []byte
		processData []byte
	)

	if my.err != nil {
		return my
	}

	// 处理切块大文件上传
	if my.isChunked && my.chunkedReader != nil {
		return my.sendChunked()
	}

	// 获取待处理数据
	if my.requestBodyBuffer != nil {
		processData = my.requestBodyBuffer.Bytes()
	} else if my.requestBody != nil {
		processData, my.err = io.ReadAll(my.requestBody)
		if my.err != nil {
			return my
		}
	}

	// 压缩处理（只要设置了 compressor 就自动启用）
	if my.compressor != nil && len(processData) > 0 {
		my.compressor.SetData(processData)
		if compressed, my.err = my.compressor.Encode(); my.err != nil {
			return my
		}
		if len(compressed) > 0 {
			processData = compressed
		}
	}

	// 加密处理
	if my.symmetricEncryptor != nil && len(processData) > 0 {
		if cipherText, my.err = my.symmetricEncryptor.Encrypt(processData); my.err != nil {
			return my
		}
	}

	if cipherText != nil {
		my.setBody(bytes.NewBuffer(cipherText))
	} else if compressed != nil {
		my.setBody(bytes.NewBuffer(compressed))
	}

	bodyReader := my.requestBody
	if bodyReader != nil && my.rateLimit > 0 {
		bodyReader = NewUploadRateReader(bodyReader, my.rateLimit*volumeInfo.KB)
	}

	if my.rawRequest, my.err = http.NewRequest(my.method, my.getURL(), bodyReader); my.err != nil {
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

// sendChunked 切块发送大文件（并行协程）
func (my *HTTPClientImpl) sendChunked() HTTPClient {
	if my.chunkedReader == nil {
		my.err = errors.New("切块读取器未初始化")
		return my
	}

	// 使用管道进行并行写入
	reader, writer := io.Pipe()

	// 启动多个goroutine进行并行读取和发送
	var wg sync.WaitGroup
	errChan := make(chan error, my.chunkedReader.numChunks)
	chunkSize := my.chunkedReader.chunkSize
	numChunks := my.chunkedReader.numChunks

	for i := range numChunks {
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()

			offset := int64(chunkIndex) * chunkSize
			remainSize := my.chunkedReader.fileSize - offset
			currentChunkSize := min(remainSize, chunkSize)

			// 读取文件块
			chunkData := make([]byte, currentChunkSize)
			my.chunkedReader.lock.Lock()
			_, err := my.chunkedReader.file.ReadAt(chunkData, offset)
			my.chunkedReader.lock.Unlock()
			if err != nil && err != io.EOF {
				errChan <- fmt.Errorf("读取文件块 %d 失败: %w", chunkIndex, err)
				return
			}

			// 压缩处理（只要设置了 compressor 就自动启用）
			if my.chunkedReader.compressor != nil {
				my.chunkedReader.compressor.SetData(chunkData)
				compressed, err := my.chunkedReader.compressor.Encode()
				if err != nil {
					errChan <- fmt.Errorf("压缩块 %d 失败: %w", chunkIndex, err)
					return
				}
				chunkData = compressed
			}

			// 加密处理
			if my.chunkedReader.encryptor != nil {
				encrypted, err := my.chunkedReader.encryptor.Encrypt(chunkData)
				if err != nil {
					errChan <- fmt.Errorf("加密块 %d 失败: %w", chunkIndex, err)
					return
				}
				chunkData = encrypted
			}

			// 写入管道
			_, err = writer.Write(chunkData)
			if err != nil {
				errChan <- fmt.Errorf("写入块 %d 失败: %w", chunkIndex, err)
				return
			}
		}(i)
	}

	// 关闭writer（所有goroutine完成写入后）
	go func() {
		wg.Wait()
		writer.Close()
	}()

	// 创建HTTP请求，使用管道作为body
	if my.rawRequest, my.err = http.NewRequest(my.method, my.getURL(), reader); my.err != nil {
		return my
	}

	for key, values := range my.headers {
		v := make([]string, 0, len(values))
		for idx := range values {
			v = append(v, cast.ToString(values[idx]))
		}
		my.rawRequest.Header[key] = append(my.rawRequest.Header[key], v...)
	}

	// 设置文件相关header
	my.rawRequest.Header.Set("X-File-Size", cast.ToString(my.chunkedReader.fileSize))
	my.rawRequest.Header.Set("X-Chunk-Count", cast.ToString(numChunks))
	my.rawRequest.Header.Set("X-Chunk-Size", cast.ToString(chunkSize))

	if len(my.cert) > 0 {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(my.cert) {
			my.err = errors.New("生成TLS证书失败")
			return my
		}

		my.transport.TLSClientConfig = &tls.Config{RootCAs: certPool}
	}

	my.client = &http.Client{}

	if my.transport != nil {
		my.client.Transport = my.transport
	}

	if my.timeout > 0 {
		my.client.Timeout = my.timeout
	}

	// 发送请求
	if my.rawResponse, my.err = my.client.Do(my.rawRequest); my.err != nil {
		return my
	}

	// 检查是否有错误
	select {
	case err := <-errChan:
		my.err = err
		return my
	default:
	}

	if my.autoCopy {
		my.parseBody()
		my.rawResponse.Body = io.NopCloser(bytes.NewBuffer(my.responseBody))
	}

	return my
}

// OK 检查响应是否成功，返回布尔值和错误信息
func (my *HTTPClientImpl) OK() error {
	if my.err != nil {
		return my.err
	}

	if my.rawResponse == nil {
		return errors.New("响应体为空")
	}

	if my.rawResponse.StatusCode > 399 {
		return fmt.Errorf("错误：%s", my.GetStatus())
	}

	if my.err != nil {
		return my.err
	}

	return nil
}

func (my *HTTPClientImpl) isNeedRetry(condition func(statusCode int, err error) bool) (needRetry bool) {
	if condition == nil {
		needRetry = my.OK() != nil
	} else {
		condition = func(statusCode int, err error) bool { return statusCode > 399 || err != nil }
		needRetry = condition(my.GetStatusCode(), my.err)
	}

	return
}

func (my *HTTPClientImpl) SendWithRetry(count uint, interval time.Duration, condition func(statusCode int, err error) bool) (HTTPClient, []error) {
	my.lock.Lock()
	defer my.lock.Unlock()

	var (
		wrongs = make([]error, 0, count)
		err    error
	)

	my.send()

	if err = my.OK(); err != nil {
		wrongs = append(wrongs, err) // 记录第一次错误
	}

	if my.isNeedRetry(condition) && count > 0 && interval > 0 {
		for attempt := range count {
			time.Sleep(interval)

			if my.rawResponse != nil && my.rawResponse.Body != nil {
				_ = my.rawResponse.Body.Close()
				my.rawResponse = nil
			}

			my.send()

			if !my.isNeedRetry(condition) {
				break
			}

			wrongs = append(wrongs, my.OK()) // 记录每次错误

			if attempt+1 == count {
				my.err = errors.New("达到最大重试次数，仍然未成功")
				break
			}
		}
	}

	return my, wrongs
}

func (my *HTTPClientImpl) Send() HTTPClient {
	my.lock.Lock()
	defer my.lock.Unlock()

	return my.send()
}

func (my *HTTPClientImpl) parseBody() {
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
	if my.rawResponse.ContentLength > fileSplitSize {
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

func (my *HTTPClientImpl) ToJSON(target any, keys ...any) HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

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
		json.Get(my.responseBody, keys...).ToVal(&target)
	} else {
		my.err = json.Unmarshal(my.responseBody, &target)
	}
	return my
}

func (my *HTTPClientImpl) ToXML(target any) HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

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

func (my *HTTPClientImpl) ToBytes() []byte {
	my.lock.RLock()
	defer my.lock.RUnlock()
	defer func() {
		if my.rawResponse != nil {
			_ = my.rawResponse.Body.Close()
		}
	}()

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

func (my *HTTPClientImpl) ToWriter(writer http.ResponseWriter) HTTPClient {
	my.lock.RLock()
	defer my.lock.RUnlock()
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

func (my *HTTPClientImpl) Error() error {
	var err error
	defer func() { my.err = nil }()

	err = my.err
	return err
}

func (my *HTTPClientImpl) GetStatusCode() int {
	return operation.NewTernary(operation.TrueFn(func() int { return my.GetRawResponse().StatusCode })).GetByValue(my.GetRawResponse() != nil)
}

func (my *HTTPClientImpl) GetStatus() string {
	return operation.NewTernary(operation.TrueFn(func() string { return my.GetRawResponse().Status })).GetByValue(my.GetRawResponse() != nil)
}
