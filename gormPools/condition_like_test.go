package gormPools_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aid297/aid/v2/gormPools"
)

// like 操作符配套 SQL 必须声明 escape '\\'（SQL 文本两个反斜杠，MySQL/PG 解析为一个反斜杠转义符）。
// 回归：曾误写为 escape '\'（单反斜杠），MySQL 默认反斜杠转义下 '\' 的闭合引号被转义 → 字符串未闭合 → Error 1064。
// 期望 SQL 文本为 name like ? escape '\\'，对应 Go 字面量 "name like ? escape '\\\\'"。
const likeEscapeSQL = "name like ? escape '\\\\'"

func likeCases() []struct {
	name     string
	operator string
	value    string
	wantArg  string
} {
	return []struct {
		name     string
		operator string
		value    string
		wantArg  string
	}{
		{"like% 右模糊", "like%", "8.0", "8.0%"},
		{"like 全模糊", "like", "8.0", "%8.0%"},
		{"%like 左模糊", "%like", "8.0", "%8.0"},
		{"like% 通配符转义", "like%", "100%", `100\%%`},
	}
}

// TestFinderConditionLikeEscape FinderImpl 的 like/like%/%like 必须生成正确的 escape '\\' 子句并转义值。
func TestFinderConditionLikeEscape(t *testing.T) {
	for _, tc := range likeCases() {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newCondMockDB(t)
			mock.ExpectQuery(regexp.QuoteMeta(likeEscapeSQL)).
				WithArgs(tc.wantArg).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a"))

			var rows []condInRow
			finder := gormPools.NewFinder(db.Model(&condInRow{}))
			finder.FindOnlyCondition(condInCondition(tc.operator, tc.value), &rows)

			if err := finder.Error(); err != nil {
				t.Fatalf("查询失败（escape 子句转义错误）：%v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("SQL 期望未满足（escape 子句或转义值不符）：%v", err)
			}
		})
	}
}

// TestFinderGConditionLikeEscape 泛型 Finder 的 like/like%/%like 同样必须生成正确的 escape '\\' 子句。
func TestFinderGConditionLikeEscape(t *testing.T) {
	for _, tc := range likeCases() {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newCondMockDB(t)
			mock.ExpectQuery(regexp.QuoteMeta(likeEscapeSQL)).
				WithArgs(tc.wantArg).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a"))

			_, _, err := gormPools.ToFinderG[condInRow](db).FindOnlyCondition(context.Background(), condInCondition(tc.operator, tc.value))
			if err != nil {
				t.Fatalf("查询失败（escape 子句转义错误）：%v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("SQL 期望未满足（escape 子句或转义值不符）：%v", err)
			}
		})
	}
}