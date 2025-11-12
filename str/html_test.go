package str

import "testing"

func TestHtml1(t *testing.T) {
	t.Run("Html", func(t *testing.T) {
		t.Logf("HTML：%s", APP.Html.New(
			HtmlP("这里是标题"),
			HtmlP(APP.Html.New(HtmlNormal("我想要跳转到"), HtmlA("百度", "https://www.baidu.com")).End()),
			HtmlA("Google", "https://www.google.com", HtmlProperty{Key: "class", Value: "my-class"}),
			HtmlBr(),
			HtmlUl("第一项", "第二项"),
		).End())
	})
}

func TestHtml2(t *testing.T) {
	t.Run("Html Table", func(t *testing.T) {
		t.Logf(
			"%s",
			APP.Html.NewString(
				HtmlTable(
					HtmlTHead(
						HtmlTh("标题1").AppendProperties(HtmlProperty{Key: "class", Value: "my-th"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
						HtmlTh("标题2").AppendProperties(HtmlProperty{Key: "class", Value: "my-th"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
					).AppendProperties(HtmlProperty{Key: "class", Value: "my-thead"}, HtmlProperty{Key: "style", Value: "background-color: #f0f0f0;"}),
					HtmlTBody(
						HtmlTr(
							HtmlTd("姓名").AppendProperties(HtmlProperty{Key: "class", Value: "my-td"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
							HtmlTd("年龄").AppendProperties(HtmlProperty{Key: "class", Value: "my-td"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
						),
						HtmlTr(
							HtmlTd("张三").AppendProperties(HtmlProperty{Key: "class", Value: "my-td"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
							HtmlTd("18").AppendProperties(HtmlProperty{Key: "class", Value: "my-td"}, HtmlProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
						),
					).AppendProperties(HtmlProperty{Key: "class", Value: "my-tbody"}),
				).AppendProperties(HtmlProperty{Key: "class", Value: "my-table"}),
			),
		)
	})
}
