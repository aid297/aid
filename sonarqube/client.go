package sonarqube

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/aid297/aid/httpClient/httpClientV2"
	"github.com/aid297/aid/str"
)

type (
	ClientService struct {
		baseURL, token string
		hc             *httpClientV2.HTTPClient
		log            *zap.Logger
	}

	ClientServiceAttributer  interface{ Register() }
	AttrClientServiceBaseURL struct{ baseURL string }
	AttrClientServiceToken   struct{ token string }
)

var (
	once      sync.Once
	clientIns *ClientService
)

func ClientServiceBaseURL(baseURL string) AttrClientServiceBaseURL {
	return AttrClientServiceBaseURL{baseURL: baseURL}
}
func (my AttrClientServiceBaseURL) Register() {
	clientIns.baseURL = my.baseURL
}

func ClientServiceToken(token string) AttrClientServiceToken {
	return AttrClientServiceToken{token: token}
}
func (my AttrClientServiceToken) Register() {
	clientIns.token = my.token
}

func (ClientService) New(attrs ...ClientServiceAttributer) *ClientService {
	once.Do(func() { clientIns = new(ClientService) })
	return clientIns.SetAttrs(attrs...)
}

func (ClientService) SetAttrs(attrs ...ClientServiceAttributer) *ClientService {
	for idx := range attrs {
		attrs[idx].Register()
	}

	return clientIns
}

func (ClientService) newHC(url string, attrs ...httpClientV2.HTTPClientAttributer) *httpClientV2.HTTPClient {
	clientIns.hc = httpClientV2.APP.HTTPClient.New(
		httpClientV2.URL(clientIns.baseURL, url),
		httpClientV2.SetHeaderValue(map[string]any{"Authorization": str.APP.Buffer.JoinString("Bearer ", clientIns.token)}).Accept(httpClientV2.AcceptJSON),
		httpClientV2.AutoCopy(false),
	).
		SetAttrs(attrs...)
	return clientIns.hc
}

func (ClientService) Find(url string, page, pageSize int, attrs ...httpClientV2.HTTPClientAttributer) *httpClientV2.HTTPClient {
	return APP.Client.newHC(url, httpClientV2.Method(http.MethodGet), httpClientV2.Queries(map[string]any{"p": page, "ps": pageSize})).SetAttrs(attrs...)
}

func (ClientService) First(url string, attrs ...httpClientV2.HTTPClientAttributer) *httpClientV2.HTTPClient {
	return APP.Client.newHC(url, httpClientV2.Method(http.MethodGet)).SetAttrs(attrs...)
}

func (ClientService) POST(url string, attrs ...httpClientV2.HTTPClientAttributer) *httpClientV2.HTTPClient {
	return APP.Client.newHC(url, httpClientV2.Method(http.MethodPost)).SetAttrs(attrs...)
}

func (ClientService) GET(url string, attrs ...httpClientV2.HTTPClientAttributer) *httpClientV2.HTTPClient {
	return APP.Client.newHC(url, httpClientV2.Method(http.MethodGet)).SetAttrs(attrs...)
}

// ProcessHTTPWrong http请求错误处理
func (my ClientService) ProcessHTTPWrong(title string, hc *httpClientV2.HTTPClient) (err error) {

	var (
		statusCode       = hc.GetStatusCode()
		codeQualityWrong CodeQualityWrong
		wrongs           []string
	)

	if statusCode > 399 {
		if err = hc.ToJSON(codeQualityWrong).Error(); err != nil {
			my.log.Error(title, zap.String(fmt.Sprintf("状态码：%d", statusCode), fmt.Sprintf("响应体：%s", string(hc.GetBody()))), zap.Error(hc.Error()))
		}

		if len(codeQualityWrong.Errors) > 0 {
			wrongs = make([]string, 0, len(codeQualityWrong.Errors))
			for idx := range codeQualityWrong.Errors {
				wrongs = append(wrongs, codeQualityWrong.Errors[idx].Msg)
			}
		}

		err = errors.New(strings.Join(wrongs, "；"))
	}
	return
}

func (ClientService) Next(title string, hc *httpClientV2.HTTPClient) (ret map[string]any, err error) {
	if err = APP.Client.ProcessHTTPWrong(title, hc); err != nil {
		return
	}

	if err = hc.ToJSON(&ret).Error(); err != nil {
		err = fmt.Errorf("%s失败：%w", title, err)
		return
	}

	return
}
