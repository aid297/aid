package httpClientV3

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
	httpBuilder := OnceHTTPBuilder().SetRootURL("http://").SetAutoCopy(true).SetTimeout(5 * time.Second)

	httpClient, err := httpBuilder.GET("baidu.com", NewHTTPClientParam())
	if err != nil {
		t.Fatalf("错误：%v", err)
	}

	res := httpClient.Send().Plain()
	if httpClient.Error != nil {
		t.Fatalf("错误：%v", httpClient.Error)
	}

	t.Logf("结果：%s", string(res))
}
