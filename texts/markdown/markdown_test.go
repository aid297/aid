package markdown_test

import (
	"strings"
	"testing"

	"github.com/aid297/aid/v2/texts/markdown"
)

func TestNormal(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Normal("普通文本"))
	if got != "普通文本" {
		t.Fatalf("期望 %q，得到 %q", "普通文本", got)
	}
}

func TestA(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.A("百度", "https://www.baidu.com"))
	if !strings.Contains(got, "[百度]") {
		t.Fatalf("缺少链接文字: %q", got)
	}
	if !strings.Contains(got, "(https://www.baidu.com)") {
		t.Fatalf("缺少链接地址: %q", got)
	}
}

func TestBr(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Normal("段落1"), markdown.Br(), markdown.Normal("段落2"))
	if !strings.Contains(got, "段落1\n\n段落2") {
		t.Fatalf("换行失败: %q", got)
	}
}

func TestUl(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Ul("第一项", "第二项", "第三项"))
	lines := strings.Split(strings.TrimSpace(got), "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "* ") {
			count++
		}
	}
	if count != 3 {
		t.Fatalf("期望 3 个列表项，得到 %d 个。完整输出: %q", count, got)
	}
}

func TestUl_Empty(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Ul())
	if got != "" {
		t.Fatalf("空列表应返回空字符串，得到: %q", got)
	}
}

func TestOl(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Ol("第一项", "第二项", "第三项"))
	if !strings.Contains(got, "1. 第一项") {
		t.Fatalf("缺少有序项 1: %q", got)
	}
	if !strings.Contains(got, "2. 第二项") {
		t.Fatalf("缺少有序项 2: %q", got)
	}
	if !strings.Contains(got, "3. 第三项") {
		t.Fatalf("缺少有序项 3: %q", got)
	}
}

func TestOl_Empty(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(markdown.Ol())
	if got != "" {
		t.Fatalf("空列表应返回空字符串，得到: %q", got)
	}
}

func TestNew_ChainMultiple(t *testing.T) {
	md := markdown.MarkdownWriterImpl{}.New(
		markdown.Normal("# 标题\n\n"),
		markdown.Normal("这是正文。"),
		markdown.Br(),
		markdown.Ul("项目A", "项目B"),
	)
	got := md.End()
	if !strings.Contains(got, "# 标题") {
		t.Fatalf("缺少标题: %q", got)
	}
	if !strings.Contains(got, "这是正文。") {
		t.Fatalf("缺少正文: %q", got)
	}
	if !strings.Contains(got, "* 项目A") {
		t.Fatalf("缺少列表项A: %q", got)
	}
}

func TestNewString_OneStep(t *testing.T) {
	got := markdown.MarkdownWriterImpl{}.NewString(
		markdown.Normal("## 标题2\n\n"),
		markdown.A("链接", "http://example.com"),
		markdown.Br(),
		markdown.Ol("有序1", "有序2"),
	)
	if !strings.Contains(got, "## 标题2") {
		t.Fatalf("缺少标题2: %q", got)
	}
	if !strings.Contains(got, "[链接]") {
		t.Fatalf("缺少链接: %q", got)
	}
	if !strings.Contains(got, "1. 有序1") {
		t.Fatalf("缺少有序项1: %q", got)
	}
}

func TestSet(t *testing.T) {
	md := markdown.MarkdownWriterImpl{}.New()
	md.Set(markdown.Normal("通过Set写入"))
	got := md.End()
	if !strings.Contains(got, "通过Set写入") {
		t.Fatalf("Set 未写入内容: %q", got)
	}
}

func TestSet_Empty(t *testing.T) {
	md := markdown.MarkdownWriterImpl{}.New()
	md.Set()
	got := md.End()
	if got != "" {
		t.Fatalf("空 Set 应返回空字符串，得到: %q", got)
	}
}

func TestGetBuffer(t *testing.T) {
	md := markdown.MarkdownWriterImpl{}.New(markdown.Normal("test"))
	buf := md.GetBuffer()
	if buf == nil {
		t.Fatal("GetBuffer 不应返回 nil")
	}
}
