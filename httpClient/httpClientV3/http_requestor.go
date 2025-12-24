package httpClientV3

import (
	"net/http"
	"net/url"
)

type HTTPRequestor struct {
	request *http.Request
	url     string
	queries url.Values
}
