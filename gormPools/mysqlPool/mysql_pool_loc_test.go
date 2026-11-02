package mysqlPool

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aid297/aid/v2/gormPools"
)

// TestLocAttr 测试 Loc Attr 设置时区参数
func TestLocAttr(t *testing.T) {
	tests := []struct {
		name     string
		loc      string
		expected string
	}{
		{
			name:     "设置Asia/Shanghai时区",
			loc:      "Asia/Shanghai",
			expected: "Asia/Shanghai",
		},
		{
			name:     "设置UTC时区",
			loc:      "UTC",
			expected: "UTC",
		},
		{
			name:     "设置Local时区",
			loc:      "Local",
			expected: "Local",
		},
		{
			name:     "空字符串时ternary返回空值",
			loc:      "",
			expected: "", // loc为空时，GetByValue(false) 返回零值 ""
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := &MySQLPoolImpl{
				username: "root",
				password: "root",
				host:     "localhost",
				port:     3306,
				database: "test_db",
				charset:  "utf8mb4",
				loc:      "UTC", // 默认值
			}

			// 应用 Loc Attr
			attr := Loc(tt.loc)
			attr(ins)

			if ins.loc != tt.expected {
				t.Errorf("Loc(%q) → loc = %q, 期望 %q", tt.loc, ins.loc, tt.expected)
			}
		})
	}
}

// TestLocDSNFormat 测试 loc 参数在 DSN 中的格式
func TestLocDSNFormat(t *testing.T) {
	tests := []struct {
		name          string
		loc           string
		expectInDSN   string
		shouldContain string
	}{
		{
			name:          "UTC时区应出现在DSN中",
			loc:           "UTC",
			shouldContain: "loc=UTC",
		},
		{
			name:          "Asia/Shanghai时区应出现在DSN中",
			loc:           "Asia/Shanghai",
			shouldContain: "loc=Asia/Shanghai",
		},
		{
			name:          "Local时区应出现在DSN中",
			loc:           "Local",
			shouldContain: "loc=Local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins := &MySQLPoolImpl{
				username: "root",
				password: "root",
				host:     "127.0.0.1",
				port:     3306,
				database: "test_db",
				charset:  "utf8mb4",
				loc:      tt.loc,
			}

			// 生成 DSN（与 New() 中的逻辑一致）
			dsn := fmt.Sprintf(
				mysqlDSNFormat+ins.costomDSNFormat,
				ins.username, ins.password, ins.host, ins.port, ins.database, ins.charset, ins.loc,
			)

			if !strings.Contains(dsn, tt.shouldContain) {
				t.Errorf("DSN = %q, 应包含 %q", dsn, tt.shouldContain)
			}

			// 验证 DSN 中没有 MISSING 占位符
			if strings.Contains(dsn, "MISSING") {
				t.Errorf("DSN = %q, 包含 MISSING（fmt.Sprintf 参数不足）", dsn)
			}
		})
	}
}

// TestLocAttrWithFullSetup 测试完整 Attr 组合下 loc 的正确设置
func TestLocAttrWithFullSetup(t *testing.T) {
	ins := &MySQLPoolImpl{
		username: "root",
		password: "root",
		host:     "localhost",
		port:     3306,
		database: "test_db",
		charset:  "utf8mb4",
		loc:      "UTC", // 默认值
	}

	// 模拟使用多个 Attr 设置
	ins.SetAttrs(
		Username("admin"),
		Password("secret"),
		Host("192.168.1.100"),
		Port(3307),
		Loc("Asia/Shanghai"),
	)

	// 验证各字段
	if ins.username != "admin" {
		t.Errorf("username = %q, 期望 %q", ins.username, "admin")
	}
	if ins.password != "secret" {
		t.Errorf("password = %q, 期望 %q", ins.password, "secret")
	}
	if ins.host != "192.168.1.100" {
		t.Errorf("host = %q, 期望 %q", ins.host, "192.168.1.100")
	}
	if ins.port != 3307 {
		t.Errorf("port = %d, 期望 %d", ins.port, 3307)
	}
	if ins.loc != "Asia/Shanghai" {
		t.Errorf("loc = %q, 期望 %q", ins.loc, "Asia/Shanghai")
	}

	// 验证最终 DSN
	dsn := fmt.Sprintf(
		mysqlDSNFormat+ins.costomDSNFormat,
		ins.username, ins.password, ins.host, ins.port, ins.database, ins.charset, ins.loc,
	)

	expectedDSN := "admin:secret@tcp(192.168.1.100:3307)/test_db?charset=utf8mb4&parseTime=True&loc=Asia/Shanghai"
	if dsn != expectedDSN {
		t.Errorf("DSN =\n  %q\n期望 =\n  %q", dsn, expectedDSN)
	}
}

// TestDefaultLocIsUTC 验证实例化时默认 loc 为 UTC
func TestDefaultLocIsUTC(t *testing.T) {
	ins := &MySQLPoolImpl{
		loc: "UTC", // New() 中的默认值
	}
	if ins.loc != "UTC" {
		t.Errorf("默认 loc = %q, 期望 %q", ins.loc, "UTC")
	}
}

// TestSetLocDirectly 测试直接调用 SetLoc 方法
func TestSetLocDirectly(t *testing.T) {
	ins := &MySQLPoolImpl{loc: "UTC"}

	ins.SetLoc("America/New_York")
	if ins.loc != "America/New_York" {
		t.Errorf("SetLoc 后 loc = %q, 期望 %q", ins.loc, "America/New_York")
	}
}

// TestSourcesAndReplicasDSN 测试读写分离场景下 sources/replicas 的 DSN 包含 loc
func TestSourcesAndReplicasDSN(t *testing.T) {
	ins := &MySQLPoolImpl{
		username: "root",
		password: "root",
		host:     "localhost",
		port:     3306,
		database: "test_db",
		charset:  "utf8mb4",
		loc:      "Asia/Tokyo",
		sources: map[string]gormPools.MySQLEndpoint{
			"src1": {Username: "root", Password: "pwd", Host: "10.0.0.1", Port: 3306},
		},
		replicas: map[string]gormPools.MySQLEndpoint{
			"rep1": {Username: "root", Password: "pwd", Host: "10.0.0.2", Port: 3306},
		},
	}

	// 验证 sources DSN
	for _, item := range ins.sources {
		dsn := fmt.Sprintf(
			mysqlDSNFormat+ins.costomDSNFormat,
			item.Username, item.Password, item.Host, item.Port, ins.database, ins.charset, ins.loc,
		)
		if !strings.Contains(dsn, "loc=Asia/Tokyo") {
			t.Errorf("source DSN = %q, 应包含 loc=Asia/Tokyo", dsn)
		}
		if strings.Contains(dsn, "MISSING") {
			t.Errorf("source DSN = %q, 包含 MISSING", dsn)
		}
	}

	// 验证 replicas DSN
	for _, item := range ins.replicas {
		dsn := fmt.Sprintf(
			mysqlDSNFormat+ins.costomDSNFormat,
			item.Username, item.Password, item.Host, item.Port, ins.database, ins.charset, ins.loc,
		)
		if !strings.Contains(dsn, "loc=Asia/Tokyo") {
			t.Errorf("replica DSN = %q, 应包含 loc=Asia/Tokyo", dsn)
		}
		if strings.Contains(dsn, "MISSING") {
			t.Errorf("replica DSN = %q, 包含 MISSING", dsn)
		}
	}
}
