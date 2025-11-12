package httpClientV2

import (
	"net/http"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	t.Run("http client request init", func(t *testing.T) {
		hc := APP.HttpClient.New(
			URL("http://www", ".baidu", ".com"),
			Method(http.MethodGet),
			Queries(map[string]any{"name": "张三", "age": 18}),
			AppendHeaderValues(nil).Authorization("username", "password", "Basic").Accept(AcceptJson).ContentType(ContentTypeJson),
			SetHeaderValues(nil).ContentType(ContentTypeJson).Accept(AcceptJson),
			JSON(map[string]any{"李四": 20, "王五": 30, "赵六": 40}),
			Timeout(5*time.Minute),
			Transport(&http.Transport{
				DisableKeepAlives:   true,             // 禁用连接复用
				MaxIdleConns:        100,              // 最大空闲连接数
				IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
				TLSHandshakeTimeout: 10 * time.Second, // TLS握手超时时间
			}),
			Cert(nil),
			AutoCopy(false),
		)

		t.Logf("%+v\n", hc)
		t.Logf("url: %s\n", hc.GetURL())
		t.Logf("method: %s\n", hc.GetMethod())
		t.Logf("queries: %+v\n", hc.GetQueries())
		t.Logf("headers: %+v\n", hc.GetHeaders())
		t.Logf("body: %s\n", string(hc.GetBody()))
		t.Logf("timeout: %s\n", hc.GetTimeout())
		t.Logf("transport: %+v\n", hc.GetTransport())
		t.Logf("error: %+v\n", hc.Error())

		t.Logf("response: %s\n", hc.Send().ToBytes())
	})
}
