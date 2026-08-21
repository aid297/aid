package html_test

import (
	"testing"

	"github.com/aid297/aid/v3/texts/html"
)

func TestHTML1(t *testing.T) {
	t.Run("HTML", func(t *testing.T) {
		t.Logf("HTML：%s", html.HTMLWriterImpl{}.New(
			html.P("这里是标题"),
			html.P(html.HTMLWriterImpl{}.New(html.Normal("我想要跳转到"), html.A("百度", "https://www.baidu.com")).End()),
			html.A("Google", "https://www.google.com", html.HTMLProperty{Key: "class", Value: "my-class"}),
			html.Br(),
			html.Ul("第一项", "第二项"),
		).End())
	})
}

func TestHTML2(t *testing.T) {
	t.Run("HTML Table", func(t *testing.T) {
		t.Logf(
			"%s",
			html.HTMLWriterImpl{}.NewString(
				html.Table(
					[]html.HTMLAttr{
						html.THead(
							[]html.HTMLAttr{
								html.Th("标题1", html.HTMLProperty{Key: "class", Value: "my-th"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
								html.Th("标题2", html.HTMLProperty{Key: "class", Value: "my-th"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
							},
							html.HTMLProperty{Key: "class", Value: "my-thead"}, html.HTMLProperty{Key: "style", Value: "background-color: #f0f0f0;"},
						),
						html.Tr(
							[]html.HTMLAttr{
								html.Td("姓名", html.HTMLProperty{Key: "class", Value: "my-td"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
								html.Td("年龄", html.HTMLProperty{Key: "class", Value: "my-td"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
							},
						),
						html.TBody(
							[]html.HTMLAttr{
								html.Tr(
									[]html.HTMLAttr{
										html.Td("张三", html.HTMLProperty{Key: "class", Value: "my-td"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
										html.Td("18", html.HTMLProperty{Key: "class", Value: "my-td"}, html.HTMLProperty{Key: "style", Value: "border: 1px solid #000; padding: 4px;"}),
									},
								),
							},
							html.HTMLProperty{Key: "class", Value: "my-tbody"},
						),
					},
					html.HTMLProperty{Key: "class", Value: "my-table"},
				),
			),
		)
	})
}
