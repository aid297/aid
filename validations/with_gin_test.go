package validations

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWithGin_DefaultsToJSONWhenMissingContentType(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Req struct {
		Name string `json:"name"`
	}

	r.POST("/bind", func(c *gin.Context) {
		form, checker := WithGin[Req](c)
		if !checker.OK() {
			c.String(http.StatusBadRequest, checker.Error().Error())
			return
		}
		c.JSON(http.StatusOK, form)
	})

	req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"name":"alice"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"name":"alice"`) {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestWithGin_UnmarshalTypeErrorChineseMessage(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Req struct {
		Age int `json:"age"`
	}

	r.POST("/bind", func(c *gin.Context) {
		_, checker := WithGin[Req](c)
		if !checker.OK() {
			c.String(http.StatusBadRequest, checker.Error().Error())
			return
		}
		c.Status(http.StatusOK)
	})

	// 传入字符串而非数字,触发 UnmarshalTypeError
	req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{"age":"not_a_number"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	// 验证中文提示包含字段名和类型信息
	if !strings.Contains(body, "参数格式解析错误(age)") {
		t.Errorf("expected field name in error, got: %s", body)
	}
	if !strings.Contains(body, "应为") && !strings.Contains(body, "int") {
		t.Errorf("expected type hint in error, got: %s", body)
	}
}

func TestWithGin_OtherBindErrorsKeepOriginalFormat(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	type Req struct {
		Name string `json:"name" binding:"required"`
	}

	r.POST("/bind", func(c *gin.Context) {
		_, checker := WithGin[Req](c)
		if !checker.OK() {
			c.String(http.StatusBadRequest, checker.Error().Error())
			return
		}
		c.Status(http.StatusOK)
	})

	// 缺少必填字段,触发 validation error 而非 UnmarshalTypeError
	req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	// 非 UnmarshalTypeError 应保持原有格式
	if !strings.Contains(body, "参数格式解析错误") {
		t.Errorf("expected original error format, got: %s", body)
	}
}
