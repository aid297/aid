package gormPools_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aid297/aid/v2/gormPools"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type condInRow struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (condInRow) TableName() string { return "cond_in_rows" }

func newCondMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 sqlmock 失败：%v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开 gorm 失败：%v", err)
	}
	return db, mock
}

func condInCondition(operator string, values ...any) *gormPools.FinderCondition {
	return &gormPools.FinderCondition{
		Queries: []gormPools.Query{{
			Option:     "and",
			Conditions: []gormPools.Condition{{Key: "name", Operator: operator, Values: values}},
		}},
	}
}

func condInArgs(values []any) []driver.Value {
	args := make([]driver.Value, len(values))
	for idx, value := range values {
		args[idx] = value
	}
	return args
}

// TestFinderConditionInBindsAllValues in/not in 多值条件必须绑定全部值，
// 而不是仅取 values[0]（回归：备份集「已过期」页签 IN [EXPIRED, CLEANED] 只返回 EXPIRED）。
func TestFinderConditionInBindsAllValues(t *testing.T) {
	cases := []struct {
		name     string
		operator string
		values   []any
		wantSQL  string
	}{
		{"in 双值", "in", []any{"a", "b"}, "name in (?,?)"},
		{"not in 双值", "not in", []any{"a", "b"}, "name not in (?,?)"},
		{"in 单值", "in", []any{"a"}, "name in (?)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newCondMockDB(t)
			mock.ExpectQuery(regexp.QuoteMeta(tc.wantSQL)).
				WithArgs(condInArgs(tc.values)...).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"))

			var rows []condInRow
			finder := gormPools.NewFinder(db.Model(&condInRow{}))
			finder.FindOnlyCondition(condInCondition(tc.operator, tc.values...), &rows)

			if err := finder.Error(); err != nil {
				t.Fatalf("查询失败：%v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("期望 2 行，实际 %d 行", len(rows))
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("SQL 期望未满足（多值未全部绑定）：%v", err)
			}
		})
	}
}

// TestFinderGConditionInBindsAllValues 泛型 Finder 的 in/not in 多值条件同样必须绑定全部值。
func TestFinderGConditionInBindsAllValues(t *testing.T) {
	cases := []struct {
		name     string
		operator string
		values   []any
		wantSQL  string
	}{
		{"in 双值", "in", []any{"a", "b"}, "name in (?,?)"},
		{"not in 双值", "not in", []any{"a", "b"}, "name not in (?,?)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newCondMockDB(t)
			mock.ExpectQuery(regexp.QuoteMeta(tc.wantSQL)).
				WithArgs(condInArgs(tc.values)...).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"))

			rows, _, err := gormPools.ToFinderG[condInRow](db).FindOnlyCondition(context.Background(), condInCondition(tc.operator, tc.values...))
			if err != nil {
				t.Fatalf("查询失败：%v", err)
			}
			if len(rows) != 2 {
				t.Fatalf("期望 2 行，实际 %d 行", len(rows))
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("SQL 期望未满足（多值未全部绑定）：%v", err)
			}
		})
	}
}