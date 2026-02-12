package httpClientV2

import (
	"net/http"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	t.Run("http client request init", func(t *testing.T) {
		hc := APP.HTTPClient.New(
			URL("http://www", ".baidu", ".com"),
			Method(http.MethodGet),
			Queries(map[string]any{"name": "张三", "age": 18}),
			SetHeaderValue(nil).Authorization("username", "password", "Basic").Accept(AcceptJSON).ContentType(ContentTypeJSON),
			SetHeaderValue(nil).ContentType(ContentTypeJSON).Accept(AcceptJSON),
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

func Test2(t *testing.T) {
	hc := APP.HTTPClient.New(
		TransportDefault(),
		URL("http://127.0.0.1:19003/project/list"),
		SetHeaderValue(map[string]any{"User-Info": "eyJ1dWlkIjoiYmIwZjJhYjItYTRiNy00ZjYxLWIzYTQtMWRlNTMzOGNkNmNkIiwibmlja25hbWUiOiLkvZkiLCJ1c2VybmFtZSI6Inl1aml6aG91IiwiZW1haWwiOiIxMzUyMjE3ODA1N0BlbWFpbC5jb20iLCJpc0FkbWluIjpmYWxzZSwidGVhbUlkIjoyLCJvd25lclRlYW1JZHMiOlszM10sImhhc1RlYW1JZHMiOlsyLDMzLDM3XSwidG9rZW4iOiIvZWREdnBCRGhManlaVWJ0TC9iVkdaZktRWlRoajJNdDc1bVBxSVduTWxQdFFGOGdxbWpJQjEzMG5MVDllelF2SmJvN1dGbG9YVzU3SW5JZkQvNkRrMU1ERmtpWVJ5aHdZRENSanZJVnArZzY2bnFwSTd5bDBCcWpDN0FBU2NrT3cyQzFWSmY1emtjaGVqbWxIMDhJNnB5Ylk2NmtzaEwwOWxGNlJMZVIzd0xNQ3l5RGNjSVpmclJQQS9IOUtNM3YvNWdVTFk3UGpKL1BSR0NzSzJlYkMyTHlEdGpTMk02MmJ1N0FRekhRQmhkSFEvN3hsQXlUTk1aT0NPU0tyWlRQWmVWK0V5b2NMMGNqNEozQ0ZzMGRFZHkvdG5iVWV3SFhIa0grU0lvUG0vQ2pPeVN4SFFEV3FCS0plemRiSFdCN2ZUMjZZcWljcjBJVmROOTB3UE1SZGtQSHJvSmVkVHNGQSs2SGhVZ1hsVHc9In0="}),
		Timeout(1*time.Second),
		Method(http.MethodPost),
		JSON(map[string]any{
			"projectUUID": "1f06786e-07bd-6868-8ef5-355bce72ed9b",
		}),
		AutoCopy(true),
	).Send()

	t.Logf("结果：%s\n", hc.ToBytes())
}
