package str

import (
	"log"
	"testing"
)

func TestMarkdown1(t *testing.T) {
	t.Run("markdown", func(t *testing.T) {
		log.Print(APP.Markdown.New(
			MarkdownNormal("这里是普通文本"),
			MarkdownBr(),
			MarkdownA("百度", "https://www.baidu.com"),
			MarkdownBr(),
			MarkdownUl("第一项", "第二项", "第三项"),
		).End())
	})
}
