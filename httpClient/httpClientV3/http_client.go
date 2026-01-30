package httpClientV3

import (
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/aid297/aid/dict/anyDictV2"
	"github.com/aid297/aid/str"
	"github.com/spf13/cast"
)

type (
	HttpCourier interface {
		GET() HttpCourier
		POST() HttpCourier
		PUT() HttpCourier
		PATCH() HttpCourier
		DELETE() HttpCourier
		HEAD() HttpCourier
		OPTIONS() HttpCourier
		TRACE() HttpCourier
		SetURL(urls ...string) HttpCourier
	}

	HttpClient struct {
		err          error
		url          string
		queries      map[string]any
		method       string
		headers      map[string][]any
		requestBody  []byte
		responseBody []byte
		timeout      time.Duration
		transport    *http.Transport
		cert         []byte
		rawRequest   *http.Request
		rawResponse  *http.Response
		client       *http.Client
		autoCopy     bool
		lock         sync.Mutex
		OK           *bool
	}

	HttpSend struct {
	}
)

func NewHttpClient() *HttpClient { return &HttpClient{lock: sync.Mutex{}, headers: map[string][]any{}} }

func (my *HttpClient) GET() HttpCourier     { my.method = http.MethodGet; return my }
func (my *HttpClient) POST() HttpCourier    { my.method = http.MethodPost; return my }
func (my *HttpClient) PUT() HttpCourier     { my.method = http.MethodPut; return my }
func (my *HttpClient) PATCH() HttpCourier   { my.method = http.MethodPatch; return my }
func (my *HttpClient) DELETE() HttpCourier  { my.method = http.MethodDelete; return my }
func (my *HttpClient) HEAD() HttpCourier    { my.method = http.MethodHead; return my }
func (my *HttpClient) OPTIONS() HttpCourier { my.method = http.MethodOptions; return my }
func (my *HttpClient) TRACE() HttpCourier   { my.method = http.MethodTrace; return my }

func (my *HttpClient) SetURL(urls ...string) HttpCourier {
	if len(urls) == 0 {
	} else if len(urls) == 1 {
		my.url = cast.ToString(urls[0])
	} else {
		my.url = str.APP.Buffer.JoinString(cast.ToStringSlice(urls)...)
	}

	return my
}

// SetQueries 设置 query 参数
func (my *HttpClient) SetQueries(queries map[string]any) HttpCourier {
	if queries == nil {
		my.queries = map[string]any{}
	}

	return my
}

// AppendQueries 追加 query 参数
func (my *HttpClient) AppendQueries(queries map[string]any) HttpCourier {
	if len(queries) > 0 {
		maps.Copy(my.queries, queries)
	}

	return my
}

// RemoveEmptyQueries 移出 query 参数中的空值
func (my *HttpClient) RemoveEmptyQueries() HttpCourier {
	my.queries = anyDictV2.New(anyDictV2.Map(my.queries)).RemoveEmpty().ToMap()
	return my
}
