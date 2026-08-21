package zstd4_test

import (
	"bytes"
	"testing"

	"github.com/klauspost/compress/zstd"

	"github.com/aid297/aid/v3/compressions"
	"github.com/aid297/aid/v3/compressions/zstd4"
)

func TestEncodeAndDecode(t *testing.T) {
	original := []byte("hello zstd compression test data, 这是一段测试数据！")

	comp, err := zstd4.New(compressions.Data(original))
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

	dec, err := zstd4.New(compressions.Data(compressed))
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
	comp, err := zstd4.New()
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
	comp, err := zstd4.New()
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

	levels := []struct {
		level int
		name  string
	}{
		{int(zstd.SpeedFastest), "SpeedFastest"},
		{int(zstd.SpeedDefault), "SpeedDefault"},
		{int(zstd.SpeedBetterCompression), "SpeedBetterCompression"},
		{int(zstd.SpeedBestCompression), "SpeedBestCompression"},
	}

	for _, tc := range levels {
		t.Run(tc.name, func(t *testing.T) {
			comp, err := zstd4.New(
				compressions.Data(original),
				zstd4.Level(tc.level),
			)
			if err != nil {
				t.Fatalf("New(level=%s) error: %v", tc.name, err)
			}

			compressed, err := comp.Encode()
			if err != nil {
				t.Fatalf("Encode(level=%s) error: %v", tc.name, err)
			}

			dec, err := zstd4.New(compressions.Data(compressed))
			if err != nil {
				t.Fatalf("New() for decode error: %v", err)
			}

			decompressed, err := dec.Decode()
			if err != nil {
				t.Fatalf("Decode(level=%s) error: %v", tc.name, err)
			}

			if !bytes.Equal(original, decompressed) {
				t.Fatalf("level=%s 解码数据与原始数据不一致", tc.name)
			}
		})
	}
}

func TestInvalidLevel(t *testing.T) {
	_, err := zstd4.New(zstd4.Level(0))
	if err == nil {
		t.Fatal("Level(0) 应返回错误")
	}

	_, err = zstd4.New(zstd4.Level(5))
	if err == nil {
		t.Fatal("Level(5) 应返回错误")
	}

	_, err = zstd4.New(zstd4.Level(-1))
	if err == nil {
		t.Fatal("Level(-1) 应返回错误")
	}
}

func TestValidLevelAttr(t *testing.T) {
	validLevels := []int{
		int(zstd.SpeedFastest),
		int(zstd.SpeedDefault),
		int(zstd.SpeedBetterCompression),
		int(zstd.SpeedBestCompression),
	}
	for _, level := range validLevels {
		_, err := zstd4.New(zstd4.Level(level))
		if err != nil {
			t.Fatalf("Level(%d) 不应返回错误: %v", level, err)
		}
	}
}
