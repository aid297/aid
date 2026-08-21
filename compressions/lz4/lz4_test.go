package lz4_test

import (
	"bytes"
	"testing"

	"github.com/aid297/aid/v3/compressions"
	"github.com/aid297/aid/v3/compressions/lz4"
)

func TestEncodeAndDecode(t *testing.T) {
	original := []byte("hello lz4 compression test data, 这是一段测试数据！")

	comp, err := lz4.New(compressions.Data(original))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	compressed, err := comp.Encode()
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("Encode() 返回空数据")
	}

	dec, err := lz4.New(compressions.Data(compressed))
	if err != nil {
		t.Fatalf("New() for decode error: %v", err)
	}

	decompressed, err := dec.Decode()
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Fatalf("解码数据与原始数据不一致:\n原始: %s\n解码: %s", original, decompressed)
	}
}

func TestEncodeEmptyData(t *testing.T) {
	comp, err := lz4.New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	compressed, err := comp.Encode()
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	if compressed != nil {
		t.Fatalf("空数据 Encode() 应返回 nil，得到: %v", compressed)
	}
}

func TestDecodeEmptyData(t *testing.T) {
	comp, err := lz4.New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	decompressed, err := comp.Decode()
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}
	if decompressed != nil {
		t.Fatalf("空数据 Decode() 应返回 nil，得到: %v", decompressed)
	}
}

func TestCompressionLevels(t *testing.T) {
	original := bytes.Repeat([]byte("abcdefghijklmnopqrstuvwxyz"), 100)

	levels := []int{0, 1, 3, 5, 9}
	for _, level := range levels {
		t.Run("level", func(t *testing.T) {
			comp, err := lz4.New(
				compressions.Data(original),
				lz4.Level(level),
			)
			if err != nil {
				t.Fatalf("New(level=%d) error: %v", level, err)
			}

			compressed, err := comp.Encode()
			if err != nil {
				t.Fatalf("Encode(level=%d) error: %v", level, err)
			}

			dec, err := lz4.New(compressions.Data(compressed))
			if err != nil {
				t.Fatalf("New() for decode error: %v", err)
			}

			decompressed, err := dec.Decode()
			if err != nil {
				t.Fatalf("Decode(level=%d) error: %v", level, err)
			}

			if !bytes.Equal(original, decompressed) {
				t.Fatalf("level=%d 解码数据与原始数据不一致", level)
			}
		})
	}
}

func TestInvalidLevel(t *testing.T) {
	_, err := lz4.New(lz4.Level(-1))
	if err == nil {
		t.Fatal("Level(-1) 应返回错误")
	}

	_, err = lz4.New(lz4.Level(10))
	if err == nil {
		t.Fatal("Level(10) 应返回错误")
	}
}

func TestValidLevelAttr(t *testing.T) {
	for level := 0; level <= 9; level++ {
		_, err := lz4.New(lz4.Level(level))
		if err != nil {
			t.Fatalf("Level(%d) 不应返回错误: %v", level, err)
		}
	}
}
