package httpClients

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aid297/aid/v2/compressions"
	"github.com/aid297/aid/v2/secrets"
	"github.com/aid297/aid/v2/str"
)

type HTTPClientAttr func(hc HTTPClient) (err error)

func URL(urls ...any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setURL(urls...) }
}

func Queries(queries map[string]any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setQueries(queries) }
}

func QueriesNotEmpty(queries map[string]any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setQueriesNotEmpty(queries) }
}

func Query(key string, value any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setQuery(key, value) }
}

func Method(method string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setMethod(method) }
}

func Headers(headers map[string][]any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setHeaders(headers) }
}

func Header(key string, value any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setHeader(key, value) }
}

func AppendHeaders(headers map[string][]any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.appendHeaders(headers) }
}

func AppendHeader(key string, value any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.appendHeader(key, value) }
}

func ContentType(ct HTTPContentType) HTTPClientAttr {
	return AppendHeader("Content-Type", HTTPContentTypes[ct])
}

func AppendContentType(ct HTTPContentType) HTTPClientAttr {
	return Header("Content-Type", HTTPContentTypes[ct])
}

func Accept(accept HTTPAccept) HTTPClientAttr {
	return AppendHeader("Accept", HTTPAccepts[accept])
}

func AppendAccept(accept HTTPAccept) HTTPClientAttr {
	return AppendHeader("Accept", HTTPAccepts[accept])
}

func Authorization(username, password, title string) HTTPClientAttr {
	return Header("Authorization", str.APP.Buffer.NewString(title, " ", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, password))).String())
}

func JSON(body any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyJSON(body); return }
}

func XML(body any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyXML(body); return }
}

func Form(body map[string]any) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyForm(body); return }
}

func FormData(fields, files map[string]string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyFormData(fields, files); return }
}

func Plain(body string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyPlain(body); return }
}

func HTML(body string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyHTML(body); return }
}

func CSS(body string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyCSS(body); return }
}

func Javascript(body string) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyJavascript(body); return }
}

func Bytes(body []byte) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyBytes(body); return }
}

func ReadCloser(body io.ReadCloser) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setBodyReadCloser(body); return }
}

func File(filename string, goroutineCount uint64) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { return hc.setBodyFile(filename, goroutineCount) }
}

func RateLimit(rate uint64) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setRateLimit(rate); return }
}

func Timeout(timeout time.Duration) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setTimeout(timeout); return }
}

func Transport(transport *http.Transport) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setTransport(transport); return }
}

func TransportDefault() HTTPClientAttr {
	return Transport(&http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	})
}

func Cert(cert []byte) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setCert(cert); return }
}

func AutoCopy(autoCopy bool) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setAutoCopy(autoCopy); return }
}

func Encrypt(symmetricEncryptor secrets.Symmetric) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setEncryptor(symmetricEncryptor); return }
}

func Compressor(compressor compressions.Compressor) HTTPClientAttr {
	return func(hc HTTPClient) (err error) { hc.setCompressor(compressor); return }
}
