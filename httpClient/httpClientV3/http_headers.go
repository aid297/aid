package httpClientV3

type HTPPHeaders struct {
	headers map[string][]string
}

func NewHTTPHeaders() *HTPPHeaders { return &HTPPHeaders{headers: map[string][]string{}} }
