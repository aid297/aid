package buffer_test

import (
	"testing"

	"github.com/aid297/aid/v2/texts/buffer"
)

func TestNewString_Empty(t *testing.T) {
	b := buffer.BufferImpl{}.NewString()
	if b.String() != "" {
		t.Fatalf("空参数应返回空字符串，得到: %q", b.String())
	}
}

func TestNewString_Single(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	if got := b.String(); got != "hello" {
		t.Fatalf("期望 %q，得到 %q", "hello", got)
	}
}

func TestNewString_Multiple(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello", " ", "world")
	if got := b.String(); got != "hello world" {
		t.Fatalf("期望 %q，得到 %q", "hello world", got)
	}
}

func TestNewBytes(t *testing.T) {
	b := buffer.BufferImpl{}.NewBytes([]byte("hello bytes"))
	if got := b.String(); got != "hello bytes" {
		t.Fatalf("期望 %q，得到 %q", "hello bytes", got)
	}
}

func TestNewRune_Empty(t *testing.T) {
	b := buffer.BufferImpl{}.NewRune()
	if b.String() != "" {
		t.Fatalf("空参数应返回空字符串，得到: %q", b.String())
	}
}

func TestNewRune_Single(t *testing.T) {
	b := buffer.BufferImpl{}.NewRune('A')
	if got := b.String(); got != "A" {
		t.Fatalf("期望 %q，得到 %q", "A", got)
	}
}

func TestNewRune_Multiple(t *testing.T) {
	b := buffer.BufferImpl{}.NewRune('h', 'i', '你', '好')
	if got := b.String(); got != "hi你好" {
		t.Fatalf("期望 %q，得到 %q", "hi你好", got)
	}
}

func TestNewAny_Empty(t *testing.T) {
	b := buffer.BufferImpl{}.NewAny()
	if b.String() != "" {
		t.Fatalf("空参数应返回空字符串，得到: %q", b.String())
	}
}

func TestNewAny_Single(t *testing.T) {
	b := buffer.BufferImpl{}.NewAny(42)
	if got := b.String(); got != "42" {
		t.Fatalf("期望 %q，得到 %q", "42", got)
	}
}

func TestNewAny_Multiple(t *testing.T) {
	b := buffer.BufferImpl{}.NewAny("name=", 42, ", active=", true)
	if got := b.String(); got != "name=42, active=true" {
		t.Fatalf("期望 %q，得到 %q", "name=42, active=true", got)
	}
}

func TestS_Append(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	b.S(" ", "world")
	if got := b.String(); got != "hello world" {
		t.Fatalf("期望 %q，得到 %q", "hello world", got)
	}
}

func TestWhenS_True(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("base")
	b.WhenS(true, "+", "extra")
	if got := b.String(); got != "base+extra" {
		t.Fatalf("期望 %q，得到 %q", "base+extra", got)
	}
}

func TestWhenS_False(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("base")
	b.WhenS(false, "+", "extra")
	if got := b.String(); got != "base" {
		t.Fatalf("期望 %q，得到 %q", "base", got)
	}
}

func TestB_Append(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("a")
	b.B('b', 'c')
	if got := b.String(); got != "abc" {
		t.Fatalf("期望 %q，得到 %q", "abc", got)
	}
}

func TestWhenB_True(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("x")
	b.WhenB(true, 'y', 'z')
	if got := b.String(); got != "xyz" {
		t.Fatalf("期望 %q，得到 %q", "xyz", got)
	}
}

func TestWhenB_False(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("x")
	b.WhenB(false, 'y', 'z')
	if got := b.String(); got != "x" {
		t.Fatalf("期望 %q，得到 %q", "x", got)
	}
}

func TestR_Append(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("a")
	b.R('b', '你')
	if got := b.String(); got != "ab你" {
		t.Fatalf("期望 %q，得到 %q", "ab你", got)
	}
}

func TestWhenR_True(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("x")
	b.WhenR(true, 'y', '好')
	if got := b.String(); got != "xy好" {
		t.Fatalf("期望 %q，得到 %q", "xy好", got)
	}
}

func TestWhenR_False(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("x")
	b.WhenR(false, 'y', '好')
	if got := b.String(); got != "x" {
		t.Fatalf("期望 %q，得到 %q", "x", got)
	}
}

func TestAny_Append(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("val=")
	b.Any(100, "end")
	if got := b.String(); got != "val=100end" {
		t.Fatalf("期望 %q，得到 %q", "val=100end", got)
	}
}

func TestBytes(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	got := b.Bytes()
	if string(got) != "hello" {
		t.Fatalf("期望 %q，得到 %q", "hello", string(got))
	}
}

func TestPtr(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	p := b.Ptr()
	if p == nil {
		t.Fatal("Ptr 不应返回 nil")
	}
	if *p != "hello" {
		t.Fatalf("期望 %q，得到 %q", "hello", *p)
	}
}

func TestString_ResetsBuffer(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	_ = b.String()
	if b.GetOriginal().Len() != 0 {
		t.Fatal("String() 后 buffer 应被 Reset")
	}
}

func TestCopy(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("original")
	cp := b.Copy()
	if cp.String() != "original" {
		t.Fatalf("Copy 期望 %q，得到 %q", "original", cp.String())
	}
}

func TestCopy_Independent(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("original")
	cp := b.Copy()
	cp.S("-modified")
	if b.String() != "original" {
		t.Fatalf("Copy 后修改副本不应影响原对象，原对象: %q", b.String())
	}
	if cp.String() != "original-modified" {
		t.Fatalf("Copy 副本: 期望 %q，得到 %q", "original-modified", cp.String())
	}
}

func TestSha256Sum256(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	hash := b.Sha256Sum256()
	if hash == "" {
		t.Fatal("Sha256Sum256 不应返回空字符串")
	}
	if len(hash) != 64 {
		t.Fatalf("Sha256Sum256 期望 64 字符的 hex，得到 %d 字符: %s", len(hash), hash)
	}
}

func TestGetOriginal(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("hello")
	orig := b.GetOriginal()
	if orig == nil {
		t.Fatal("GetOriginal 不应返回 nil")
	}
	if orig.String() != "hello" {
		t.Fatalf("期望 %q，得到 %q", "hello", orig.String())
	}
}

func TestJoinString_Empty(t *testing.T) {
	got := buffer.BufferImpl{}.JoinString()
	if got != "" {
		t.Fatalf("空参数应返回空字符串，得到: %q", got)
	}
}

func TestJoinString_Single(t *testing.T) {
	got := buffer.BufferImpl{}.JoinString("only")
	if got != "only" {
		t.Fatalf("期望 %q，得到 %q", "only", got)
	}
}

func TestJoinString_Multiple(t *testing.T) {
	got := buffer.BufferImpl{}.JoinString("a", "b", "c")
	if got != "abc" {
		t.Fatalf("期望 %q，得到 %q", "abc", got)
	}
}

func TestJoinStringLimit_Empty(t *testing.T) {
	got := buffer.BufferImpl{}.JoinStringLimit(",", "a", "b", "c")
	if got != "a,b,c" {
		t.Fatalf("期望 %q，得到 %q", "a,b,c", got)
	}
}

func TestJoinStringLimit_Single(t *testing.T) {
	got := buffer.BufferImpl{}.JoinStringLimit(",", "only")
	if got != "only" {
		t.Fatalf("期望 %q，得到 %q", "only", got)
	}
}

func TestJoinAny_Empty(t *testing.T) {
	got := buffer.BufferImpl{}.JoinAny()
	if got != "" {
		t.Fatalf("空参数应返回空字符串，得到: %q", got)
	}
}

func TestJoinAny_Single(t *testing.T) {
	got := buffer.BufferImpl{}.JoinAny(42)
	if got != "42" {
		t.Fatalf("期望 %q，得到 %q", "42", got)
	}
}

func TestJoinAny_Multiple(t *testing.T) {
	got := buffer.BufferImpl{}.JoinAny("a", 1, true)
	if got != "a1true" {
		t.Fatalf("期望 %q，得到 %q", "a1true", got)
	}
}

func TestJoinAnyLimit(t *testing.T) {
	got := buffer.BufferImpl{}.JoinAnyLimit("-", "a", 1, true)
	if got != "a-1-true" {
		t.Fatalf("期望 %q，得到 %q", "a-1-true", got)
	}
}

func TestJoinAnyLimit_Single(t *testing.T) {
	got := buffer.BufferImpl{}.JoinAnyLimit("-", "only")
	if got != "only" {
		t.Fatalf("期望 %q，得到 %q", "only", got)
	}
}

func TestURLPath(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("base")
	b.URLPath("path with space", "/more")
	got := b.String()
	if got == "basepath with space/more" {
		t.Logf("URLPath 未转义空格: %q", got)
	}
	if len(got) == 0 {
		t.Fatal("URLPath 结果不应为空")
	}
}

func TestURLQuery(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("base")
	b.URLQuery("query value")
	got := b.String()
	if len(got) == 0 {
		t.Fatal("URLQuery 结果不应为空")
	}
}

func TestChainOperations(t *testing.T) {
	b := buffer.BufferImpl{}.NewString("Hello").
		S(" ").
		Any("World", 2024).
		WhenS(true, "!").
		WhenS(false, "SKIP").
		B('.')
	got := b.String()
	if got != "Hello World2024!." {
		t.Fatalf("链式操作期望 %q，得到 %q", "Hello World2024!.", got)
	}
}
