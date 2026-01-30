package httpClientV3

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/aid297/aid/str"
	"github.com/spf13/cast"
)

type (
	HTTPBuild struct {
		rootURL   string
		transport *http.Transport
		cert      []byte
		timeout   time.Duration
		autoCopy  bool
	}

	HTTPBuilder interface {
		SetRootURL(urls ...string) HTTPBuilder
		SetTransport(transport *http.Transport, cert []byte) (HTTPBuilder, error)
		GetTransport() *http.Transport
		SetTimeout(timeout time.Duration) HTTPBuilder
		GetTimeout() time.Duration
		SetAutoCopy(autoCopy bool) HTTPBuilder
		GetAutoCopy() bool
		GET(urlStr string, params *HTTPClientGetParam) (httpClient *HTTPClient, err error)
	}
)

var (
	httpBuildOnce sync.Once
	httpBuildIns  *HTTPBuild
)

// OnceHTTPBuilder 获取单例对象
func OnceHTTPBuilder() HTTPBuilder {
	httpBuildOnce.Do(func() {
		httpBuildIns = &HTTPBuild{
			timeout: time.Second * 5,
			transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}
	})
	return httpBuildIns
}

// SetRootURL 设置根路由
func (*HTTPBuild) SetRootURL(urls ...string) HTTPBuilder {
	if len(urls) == 0 {
	} else if len(urls) == 1 {
		httpBuildIns.rootURL = cast.ToString(urls[0])
	} else {
		httpBuildIns.rootURL = str.APP.Buffer.JoinString(cast.ToStringSlice(urls)...)
	}

	return httpBuildIns
}

// SetTransport 设置 Transport
func (*HTTPBuild) SetTransport(transport *http.Transport, cert []byte) (HTTPBuilder, error) {
	if transport != nil {
		httpBuildIns.transport = transport
	}

	if len(httpBuildIns.cert) > 0 {
		httpBuildIns.transport.TLSClientConfig = &tls.Config{RootCAs: x509.NewCertPool()}
		if !httpBuildIns.transport.TLSClientConfig.RootCAs.AppendCertsFromPEM(cert) {
			return httpBuildIns, errors.New("生成TLS证书失败")
		}
	}

	return httpBuildIns, nil
}

// GetTransport 获取 Transport
func (*HTTPBuild) GetTransport() *http.Transport { return httpBuildIns.transport }

// SetTimeout 设置超时
func (*HTTPBuild) SetTimeout(timeout time.Duration) HTTPBuilder {
	httpBuildIns.timeout = timeout
	return httpBuildIns
}

// GetTimeout 获取超时
func (*HTTPBuild) GetTimeout() time.Duration { return httpBuildIns.timeout }

// SetAutoCopy 设置是否自动备份响应体
func (*HTTPBuild) SetAutoCopy(autoCopy bool) HTTPBuilder {
	httpBuildIns.autoCopy = autoCopy
	return httpBuildIns
}

// GetAutoCopy 获取是否自动备份响应体
func (*HTTPBuild) GetAutoCopy() bool { return httpBuildIns.autoCopy }

// GET 获取 GET 请求
func (*HTTPBuild) GET(urlStr string, params *HTTPClientGetParam) (httpClient *HTTPClient, err error) {
	urlBuffer := str.APP.Buffer.NewString(httpBuildIns.rootURL).S(urlStr)

	if params != nil && params.Queries != nil {
		queries := url.Values{}
		if params.Queries.Length() > 0 {
			params.Queries.Each(func(key, value string) { queries.Add(key, value) })
			urlBuffer = urlBuffer.S("?").S(queries.Encode())
		}
	}

	httpClient, err = NewHTTPClientGET(httpBuildIns, urlBuffer.String(), params.Headers)
	return
}
