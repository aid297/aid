package mysqlPool

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// BinaryUUID 以 binary(16) 存储的 UUID 类型，底层为 [16]byte。
// 实现 driver.Valuer/sql.Scanner 用于数据库二进制存储，
// 实现 json.Marshaler/json.Unmarshaler 用于 JSON 字符串序列化。
type BinaryUUID uuid.UUID

// ToUUID 转换为 uuid.UUID
func (b BinaryUUID) ToUUID() uuid.UUID { return uuid.UUID(b) }

// String 返回标准 UUID 字符串格式（36字符）
func (b BinaryUUID) String() string {
	if b.IsNil() {
		return ""
	}

	return uuid.UUID(b).String()
}

// IsNil 判断是否为零值
func (b BinaryUUID) IsNil() bool { return uuid.UUID(b) == uuid.Nil }

func (b BinaryUUID) NotNil() bool { return !b.IsNil() }

// Value 实现 driver.Valuer，返回 16 字节原始二进制用于 binary(16) 存储
// 注意：必须使用值接收者，否则 GORM 反射提取字段时可能无法调用
func (b BinaryUUID) Value() (driver.Value, error) {
	return b[:], nil
}

func (b BinaryUUID) Equal(something BinaryUUID) bool { return bytes.Equal(b[:], something[:]) }

func (b BinaryUUID) NotEqual(something BinaryUUID) bool { return !b.Equal(something) }

// Scan 实现 sql.Scanner，从数据库读取 16 字节二进制或字符串
func (b *BinaryUUID) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		if len(v) != 16 {
			return fmt.Errorf("BinaryUUID.Scan: expected 16 bytes, got %d", len(v))
		}
		copy(b[:], v)
		return nil
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("BinaryUUID.Scan: %w", err)
		}
		*b = BinaryUUID(parsed)
		return nil
	case nil:
		*b = BinaryUUID{}
		return nil
	default:
		return fmt.Errorf("BinaryUUID.Scan: unsupported type %T", src)
	}
}

// MarshalJSON 序列化为 JSON 字符串（标准 UUID 格式）
func (b BinaryUUID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, b.String())), nil
}

// UnmarshalJSON 从 JSON 字符串反序列化
func (b *BinaryUUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*b = BinaryUUID{}
		return nil
	}
	// 去除首尾引号
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("BinaryUUID.UnmarshalJSON: %w", err)
	}
	*b = BinaryUUID(parsed)
	return nil
}

// BinaryFromUUID 将 uuid.UUID 转换为 BinaryUUID
func BinaryFromUUID(u uuid.UUID) BinaryUUID { return BinaryUUID(u) }

// BinaryFromString 从字符串解析为 BinaryUUID
func BinaryFromString(s string) (BinaryUUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return BinaryUUID{}, err
	}
	return BinaryUUID(u), nil
}

// MustBinaryFromString 从字符串解析为 BinaryUUID，失败时 panic
func MustBinaryFromString(s string) BinaryUUID { return BinaryFromUUID(uuid.MustParse(s)) }

// BinaryFromBytes 从 []uint8（16字节）创建 BinaryUUID
func BinaryFromBytes(data []uint8) (BinaryUUID, error) {
	u, err := uuid.FromBytes(data)
	if err != nil {
		return BinaryUUID{}, fmt.Errorf("BinaryFromBytes: %w", err)
	}
	return BinaryUUID(u), nil
}

// MustBinaryFromBytes 从 []uint8（16字节）创建 BinaryUUID，失败时 panic
func MustBinaryFromBytes(data []uint8) BinaryUUID {
	b, err := BinaryFromBytes(data)
	if err != nil {
		panic(err)
	}
	return b
}

// NewBinaryUUIDV7 生成新的 UUIDv7 并返回 BinaryUUID
func NewBinaryUUIDV7() BinaryUUID { return BinaryFromUUID(uuid.Must(uuid.NewV7())) }
