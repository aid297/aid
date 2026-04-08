package validatorV3

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
