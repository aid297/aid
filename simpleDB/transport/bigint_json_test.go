package transport

import (
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestSQLParams_BigIntDoesNotLosePrecision(t *testing.T) {
	t.Parallel()

	const big = "9007199254740993"
	raw := `{"sql":"SELECT * FROM t WHERE id=:id AND ts=?","paramMap":{"id":` + big + `},"paramList":[` + big + `]}`

	var req SQLExecuteRequest
	if err := jsoniter.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	bound, err := bindSQLParams(req.SQL, req.mergedParamMap(), []any(req.ParamList))
	if err != nil {
		t.Fatalf("bindSQLParams: %v", err)
	}

	if !strings.Contains(bound, big) {
		t.Fatalf("bound sql = %q, want contains %q", bound, big)
	}
	if strings.Contains(bound, "9007199254740992") {
		t.Fatalf("bound sql = %q, want not contains rounded value", bound)
	}
}

