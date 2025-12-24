package httpClientV3

import (
	"net/url"

	"github.com/aid297/aid/dict/anyDictV2"
	"github.com/aid297/aid/str"
	"github.com/spf13/cast"
)

type (
	HTTPRequestorAttributer interface{ RegisterAttr(r *HTTPRequestor) }

	AttrHTTPRequestorURL     struct{ url string }
	AttrHTTPRequestorQueries struct{ queries map[string]any }
	AttrHTTPRequestorHeaders struct {
		headers map[string]any
		mode    string
	}
)

// ******************** URL ********************
func (AttrHTTPRequestorURL) Set(urls ...any) HTTPRequestorAttributer {
	return AttrHTTPRequestorURL{str.APP.Buffer.JoinAnyLimit("/", urls...)}
}
func (my AttrHTTPRequestorURL) RegisterAttr(r *HTTPRequestor) { r.url = my.url }

// ******************** Queries ********************
func (AttrHTTPRequestorQueries) Set(queries map[string]any) HTTPRequestorAttributer {
	return AttrHTTPRequestorQueries{queries: queries}
}
func (my AttrHTTPRequestorQueries) RemoveEmpty() HTTPRequestorAttributer {
	my.queries = anyDictV2.New(anyDictV2.Map(my.queries)).RemoveEmpty().ToMap()
	return my
}
func (my AttrHTTPRequestorQueries) RegisterAttr(r *HTTPRequestor) {
	r.queries = url.Values{}
	if len(my.queries) > 0 {
		for k, v := range my.queries {
			r.queries.Add(k, cast.ToString(v))
		}
	}
}

// ******************** Headers ********************
func (AttrHTTPRequestorHeaders) Set(headers map[string]any) HTTPRequestorAttributer {
	return AttrHTTPRequestorHeaders{headers: headers, mode: "SET"}
}
func (AttrHTTPRequestorHeaders) Append(headers map[string]any) HTTPRequestorAttributer {
	return AttrHTTPRequestorHeaders{headers: headers, mode: "APPEND"}
}
func (my AttrHTTPRequestorHeaders) ContentType(contentType ContentType) HTTPRequestorAttributer {
	my.headers["Content-Type"] = string(contentType)
	return my
}
func (my AttrHTTPRequestorHeaders) Accept(accept Accept) HTTPRequestorAttributer {
	my.headers["Accept"] = string(accept)
	return my
}
func (my AttrHTTPRequestorHeaders) RegisterAttr(r *HTTPRequestor) {
	switch my.mode {
	case "SET":
		for key := range my.headers {
			r.request.Header.Set(key, cast.ToString(my.headers[key]))
		}
	case "APPEND":
		for key := range my.headers {
			r.request.Header.Add(key, cast.ToString(my.headers[key]))
		}
	}
}

// ******************** Body ********************
