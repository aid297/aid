package transport

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aid297/aid/simpleDB/driver"
	"github.com/gin-gonic/gin"
)

func TestHTTPServer_LoginSuccess(t *testing.T) {
	dir := t.TempDir()
	server := New.HTTP(dir)

	body, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	req := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !resp.Success || resp.User == nil {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Token == nil || resp.Token.AccessToken == "" {
		t.Fatalf("expected access token in response: %+v", resp)
	}
	if resp.User.Username != "sdb" {
		t.Fatalf("username = %s, want sdb", resp.User.Username)
	}
	if !resp.User.IsAdmin {
		t.Fatal("expected sdb to be admin")
	}
	if len(resp.User.Roles) == 0 || resp.User.Roles[0] != "super_admin" {
		t.Fatalf("unexpected roles: %+v", resp.User.Roles)
	}
	claims, err := server.ParseToken(resp.Token.AccessToken)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Username != "sdb" {
		t.Fatalf("token username = %s, want sdb", claims.Username)
	}
}

func TestHTTPServer_LoginUnauthorized(t *testing.T) {
	dir := t.TempDir()
	server := New.HTTP(dir)

	body, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "wrong-password"})
	req := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(body))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var resp LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Success {
		t.Fatalf("unexpected success response: %+v", resp)
	}
	if resp.Error == nil || resp.Error.Code != "unauthorized" {
		t.Fatalf("unexpected error response: %+v", resp)
	}
}

func TestHTTPServer_LoginMethodNotAllowed(t *testing.T) {
	server := New.HTTP(t.TempDir())
	req := httptest.NewRequest(http.MethodGet, DefaultLoginPath, nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHTTPServer_EngineAndServe(t *testing.T) {
	server := New.HTTP(t.TempDir())
	if server.Engine() == nil {
		t.Fatal("expected gin engine")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	done := make(chan error, 1)
	go func() {
		done <- server.Serve(listener)
	}()

	body, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	resp, err := http.Post("http://"+listener.Addr().String()+DefaultLoginPath, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	_ = listener.Close()
	serveErr := <-done
	if serveErr != nil && !strings.Contains(serveErr.Error(), "use of closed network connection") {
		t.Fatalf("serve error: %v", serveErr)
	}
}

func TestHTTPServer_AuthMiddleware(t *testing.T) {
	server := New.HTTP(t.TempDir())
	server.Engine().GET("/me", server.AuthMiddleware(), func(ctx *gin.Context) {
		claims, ok := UserFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"username": claims.Username, "isAdmin": claims.IsAdmin})
	})

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "sdb") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHTTPServer_AuthMiddlewareRejectsExpiredToken(t *testing.T) {
	server := New.HTTP(t.TempDir())
	server.tokenManager = NewTokenManager(server.Database, server.TokenSecret, time.Second)
	server.Engine().GET("/me", server.AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	issued, err := server.tokenManager.Issue(&driver.AuthenticatedUser{ID: "1", Username: "sdb", IsAdmin: true})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	if _, err = server.ParseToken(issued.AccessToken); err != nil {
		t.Fatalf("unexpected parse failure before expiry: %v", err)
	}

	server.tokenManager = NewTokenManager(server.Database, server.TokenSecret, 12*time.Hour)
	expiredManager := NewTokenManager(server.Database, server.TokenSecret, time.Second)
	expired, err := expiredManager.Issue(&driver.AuthenticatedUser{ID: "1", Username: "sdb", IsAdmin: true})
	if err != nil {
		t.Fatalf("issue expiring token: %v", err)
	}
	time.Sleep(1100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+expired.AccessToken)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHTTPServer_RefreshAndLogout(t *testing.T) {
	server := New.HTTP(t.TempDir())

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	refreshReq := httptest.NewRequest(http.MethodPost, DefaultRefreshPath, nil)
	refreshReq.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
	refreshRec := httptest.NewRecorder()
	server.ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, want %d", refreshRec.Code, http.StatusOK)
	}

	var refreshResp LoginResponse
	if err := json.Unmarshal(refreshRec.Body.Bytes(), &refreshResp); err != nil {
		t.Fatalf("unmarshal refresh response: %v", err)
	}
	if refreshResp.Token == nil || refreshResp.Token.AccessToken == loginResp.Token.AccessToken {
		t.Fatalf("expected new refreshed token: %+v", refreshResp)
	}
	if _, err := server.ParseToken(loginResp.Token.AccessToken); err == nil {
		t.Fatal("old token should be revoked after refresh")
	}

	logoutReq := httptest.NewRequest(http.MethodPost, DefaultLogoutPath, nil)
	logoutReq.Header.Set("Authorization", "Bearer "+refreshResp.Token.AccessToken)
	logoutRec := httptest.NewRecorder()
	server.ServeHTTP(logoutRec, logoutReq)
	if logoutRec.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want %d", logoutRec.Code, http.StatusOK)
	}
	if _, err := server.ParseToken(refreshResp.Token.AccessToken); err == nil {
		t.Fatal("token should be revoked after logout")
	}
}

func TestHTTPServer_RequireRolesAndPermissions(t *testing.T) {
	server := New.HTTP(t.TempDir())
	server.Engine().GET("/admin", server.AuthMiddleware(), server.RequireRoles("super_admin"), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})
	server.Engine().GET("/reports", server.AuthMiddleware(), server.RequirePermissions("report.read"), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	adminToken, err := server.tokenManager.Issue(&driver.AuthenticatedUser{ID: "1", Username: "sdb", IsAdmin: true, Roles: []string{"super_admin"}})
	if err != nil {
		t.Fatalf("issue admin token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken.AccessToken)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin status = %d, want %d", rec.Code, http.StatusOK)
	}

	reportToken, err := server.tokenManager.Issue(&driver.AuthenticatedUser{ID: "2", Username: "reader", Permissions: []string{"report.read"}})
	if err != nil {
		t.Fatalf("issue permission token: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/reports", nil)
	req.Header.Set("Authorization", "Bearer "+reportToken.AccessToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("reports status = %d, want %d", rec.Code, http.StatusOK)
	}

	deniedToken, err := server.tokenManager.Issue(&driver.AuthenticatedUser{ID: "3", Username: "guest"})
	if err != nil {
		t.Fatalf("issue denied token: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+deniedToken.AccessToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("denied status = %d, want %d", rec.Code, http.StatusForbidden)
	}

	nonSuperAdminToken, err := server.tokenManager.Issue(&driver.AuthenticatedUser{ID: "4", Username: "other_admin", IsAdmin: true})
	if err != nil {
		t.Fatalf("issue non-super-admin token: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+nonSuperAdminToken.AccessToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("non-super-admin should be forbidden, status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestHTTPServer_RegisterInactiveUserCannotLogin(t *testing.T) {
	server := New.HTTP(t.TempDir())

	registerBody, _ := json.Marshal(RegisterRequest{Username: "new_user", Password: "123456", DisplayName: "New User"})
	registerReq := httptest.NewRequest(http.MethodPost, DefaultRegisterPath, bytes.NewReader(registerBody))
	registerRec := httptest.NewRecorder()
	server.ServeHTTP(registerRec, registerReq)

	if registerRec.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerRec.Code, http.StatusCreated)
	}

	var registerResp LoginResponse
	if err := json.Unmarshal(registerRec.Body.Bytes(), &registerResp); err != nil {
		t.Fatalf("unmarshal register response: %v", err)
	}
	if !registerResp.Success || registerResp.User == nil {
		t.Fatalf("unexpected register response: %+v", registerResp)
	}
	if registerResp.User.Status != "inactive" {
		t.Fatalf("status = %s, want inactive", registerResp.User.Status)
	}
	if registerResp.User.IsAdmin {
		t.Fatal("registered user should not be admin")
	}
	if len(registerResp.User.Roles) != 0 || len(registerResp.User.Permissions) != 0 {
		t.Fatalf("new registered user should not have roles/permissions: %+v", registerResp.User)
	}

	loginBody, _ := json.Marshal(LoginRequest{Username: "new_user", Password: "123456"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusUnauthorized {
		t.Fatalf("inactive login status = %d, want %d", loginRec.Code, http.StatusUnauthorized)
	}
}

func TestHTTPServer_AdminActivateAndAssignRoles(t *testing.T) {
	server := New.HTTP(t.TempDir(), WithInitPassword("init-secret"))
	server.Engine().GET("/admin-probe", server.AuthMiddleware(), server.RequireRoles("super_admin"), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true})
	})

	registerBody, _ := json.Marshal(RegisterRequest{Username: "ops_user", Password: "654321", DisplayName: "Ops User"})
	registerReq := httptest.NewRequest(http.MethodPost, DefaultRegisterPath, bytes.NewReader(registerBody))
	registerRec := httptest.NewRecorder()
	server.ServeHTTP(registerRec, registerReq)
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerRec.Code, http.StatusCreated)
	}

	activateReq := httptest.NewRequest(http.MethodPost, DefaultActivatePath, bytes.NewReader([]byte(`{"username":"ops_user"}`)))
	activateRec := httptest.NewRecorder()
	server.ServeHTTP(activateRec, activateReq)
	if activateRec.Code != http.StatusUnauthorized {
		t.Fatalf("activate without token status = %d, want %d", activateRec.Code, http.StatusUnauthorized)
	}

	adminLoginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	adminLoginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(adminLoginBody))
	adminLoginRec := httptest.NewRecorder()
	server.ServeHTTP(adminLoginRec, adminLoginReq)
	if adminLoginRec.Code != http.StatusOK {
		t.Fatalf("admin login status = %d, want %d", adminLoginRec.Code, http.StatusOK)
	}
	var adminLoginResp LoginResponse
	if err := json.Unmarshal(adminLoginRec.Body.Bytes(), &adminLoginResp); err != nil {
		t.Fatalf("unmarshal admin login response: %v", err)
	}
	if adminLoginResp.Token == nil {
		t.Fatal("expected admin token")
	}

	activateReq = httptest.NewRequest(http.MethodPost, DefaultActivatePath, bytes.NewReader([]byte(`{"username":"ops_user"}`)))
	activateReq.Header.Set("Authorization", "Bearer "+adminLoginResp.Token.AccessToken)
	activateRec = httptest.NewRecorder()
	server.ServeHTTP(activateRec, activateReq)
	if activateRec.Code != http.StatusOK {
		t.Fatalf("activate status = %d, want %d", activateRec.Code, http.StatusOK)
	}

	assignReq := httptest.NewRequest(http.MethodPost, DefaultAssignRolesPath, bytes.NewReader([]byte(`{"username":"ops_user","roles":["super_admin"]}`)))
	assignReq.Header.Set("Authorization", "Bearer "+adminLoginResp.Token.AccessToken)
	assignRec := httptest.NewRecorder()
	server.ServeHTTP(assignRec, assignReq)
	if assignRec.Code != http.StatusUnauthorized {
		t.Fatalf("assign role status = %d, want %d, body=%s", assignRec.Code, http.StatusUnauthorized, assignRec.Body.String())
	}

	assignPermissionReq := httptest.NewRequest(http.MethodPost, DefaultAssignRolePermissionsPath, bytes.NewReader([]byte(`{"roleCode":"super_admin","permissions":["table.demo.select"]}`)))
	assignPermissionReq.Header.Set("Authorization", "Bearer "+adminLoginResp.Token.AccessToken)
	assignPermissionRec := httptest.NewRecorder()
	server.ServeHTTP(assignPermissionRec, assignPermissionReq)
	if assignPermissionRec.Code != http.StatusOK {
		t.Fatalf("assign role permissions status = %d, want %d, body=%s", assignPermissionRec.Code, http.StatusOK, assignPermissionRec.Body.String())
	}

	userLoginBody, _ := json.Marshal(LoginRequest{Username: "ops_user", Password: "654321"})
	userLoginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(userLoginBody))
	userLoginRec := httptest.NewRecorder()
	server.ServeHTTP(userLoginRec, userLoginReq)
	if userLoginRec.Code != http.StatusOK {
		t.Fatalf("user login after activate status = %d, want %d", userLoginRec.Code, http.StatusOK)
	}

	var userLoginResp LoginResponse
	if err := json.Unmarshal(userLoginRec.Body.Bytes(), &userLoginResp); err != nil {
		t.Fatalf("unmarshal user login response: %v", err)
	}
	if userLoginResp.Token == nil {
		t.Fatal("expected user token")
	}

	adminProbeReq := httptest.NewRequest(http.MethodGet, "/admin-probe", nil)
	adminProbeReq.Header.Set("Authorization", "Bearer "+userLoginResp.Token.AccessToken)
	adminProbeRec := httptest.NewRecorder()
	server.ServeHTTP(adminProbeRec, adminProbeReq)
	if adminProbeRec.Code != http.StatusForbidden {
		t.Fatalf("admin probe status = %d, want %d", adminProbeRec.Code, http.StatusForbidden)
	}

	deactivateReq := httptest.NewRequest(http.MethodPost, DefaultDeactivatePath, bytes.NewReader([]byte(`{"username":"ops_user"}`)))
	deactivateReq.Header.Set("Authorization", "Bearer "+adminLoginResp.Token.AccessToken)
	deactivateRec := httptest.NewRecorder()
	server.ServeHTTP(deactivateRec, deactivateReq)
	if deactivateRec.Code != http.StatusOK {
		t.Fatalf("deactivate status = %d, want %d", deactivateRec.Code, http.StatusOK)
	}

	userLoginReq = httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(userLoginBody))
	userLoginRec = httptest.NewRecorder()
	server.ServeHTTP(userLoginRec, userLoginReq)
	if userLoginRec.Code != http.StatusUnauthorized {
		t.Fatalf("user login after deactivate status = %d, want %d", userLoginRec.Code, http.StatusUnauthorized)
	}
}

func TestHTTPServer_InitSDBPassword(t *testing.T) {
	server := New.HTTP(t.TempDir(), WithInitPassword("init-secret"))

	badReq := httptest.NewRequest(http.MethodPost, DefaultInitSDBPasswordPath, bytes.NewReader([]byte(`{"password":"wrong"}`)))
	badRec := httptest.NewRecorder()
	server.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusForbidden {
		t.Fatalf("wrong init password status = %d, want %d", badRec.Code, http.StatusForbidden)
	}

	okReq := httptest.NewRequest(http.MethodPost, DefaultInitSDBPasswordPath, bytes.NewReader([]byte(`{"password":"init-secret"}`)))
	okRec := httptest.NewRecorder()
	server.ServeHTTP(okRec, okReq)
	if okRec.Code != http.StatusOK {
		t.Fatalf("init sdb password status = %d, want %d, body=%s", okRec.Code, http.StatusOK, okRec.Body.String())
	}
	if server.getInitPassword() == "init-secret" {
		t.Fatal("init password should be rotated after successful initialization")
	}

	againReq := httptest.NewRequest(http.MethodPost, DefaultInitSDBPasswordPath, bytes.NewReader([]byte(`{"password":"init-secret"}`)))
	againRec := httptest.NewRecorder()
	server.ServeHTTP(againRec, againReq)
	if againRec.Code != http.StatusForbidden {
		t.Fatalf("old init password should be rejected, status = %d, want %d", againRec.Code, http.StatusForbidden)
	}

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login after init password status = %d, want %d", loginRec.Code, http.StatusOK)
	}
}

func TestHTTPServer_SQLExecute_Unauthorized(t *testing.T) {
	server := New.HTTP(t.TempDir())
	req := httptest.NewRequest(http.MethodPost, DefaultSQLExecutePath, bytes.NewReader([]byte(`{"sql":"SELECT * FROM users"}`)))
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHTTPServer_SQLExecute_CRUDWithParams(t *testing.T) {
	server := New.HTTP(t.TempDir())

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginRec.Code, http.StatusOK)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, DefaultSQLExecutePath, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		return rec
	}

	createRec := request(`{"sql":"CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)"}`)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d, body=%s", createRec.Code, http.StatusOK, createRec.Body.String())
	}

	insertRec := request(`{"sql":"INSERT INTO users (username, age) VALUES (:username, :age)","paramMap":{"username":"alice","age":20}}`)
	if insertRec.Code != http.StatusOK {
		t.Fatalf("insert status = %d, want %d, body=%s", insertRec.Code, http.StatusOK, insertRec.Body.String())
	}

	selectRec := request(`{"sql":"SELECT * FROM users WHERE username=:username","paramMap":{"username":"alice"}}`)
	if selectRec.Code != http.StatusOK {
		t.Fatalf("select status = %d, want %d, body=%s", selectRec.Code, http.StatusOK, selectRec.Body.String())
	}

	var sqlResp SQLExecuteResponse
	if err := json.Unmarshal(selectRec.Body.Bytes(), &sqlResp); err != nil {
		t.Fatalf("unmarshal sql response: %v", err)
	}
	if !sqlResp.Success || sqlResp.Result == nil {
		t.Fatalf("unexpected sql response: %+v", sqlResp)
	}
	if len(sqlResp.Result.Rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(sqlResp.Result.Rows))
	}
	if got := sqlResp.Result.Rows[0]["username"]; got != "alice" {
		t.Fatalf("username = %v, want alice", got)
	}
}

func TestHTTPServer_SQLExecute_BatchInsertAndUpdate(t *testing.T) {
	server := New.HTTP(t.TempDir())

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginRec.Code, http.StatusOK)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, DefaultSQLExecutePath, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		return rec
	}

	createRec := request(`{"sql":"CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)"}`)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d, body=%s", createRec.Code, http.StatusOK, createRec.Body.String())
	}

	insertRec := request(`{"sql":"INSERT INTO users (username, age) VALUES ('alice', 20), ('bob', 21)"}`)
	if insertRec.Code != http.StatusOK {
		t.Fatalf("batch insert status = %d, want %d, body=%s", insertRec.Code, http.StatusOK, insertRec.Body.String())
	}

	var insertResp SQLExecuteResponse
	if err := json.Unmarshal(insertRec.Body.Bytes(), &insertResp); err != nil {
		t.Fatalf("unmarshal insert response: %v", err)
	}
	if !insertResp.Success || insertResp.Result == nil {
		t.Fatalf("unexpected insert response: %+v", insertResp)
	}
	if len(insertResp.Result.InsertedRows) != 2 {
		t.Fatalf("inserted rows = %d, want 2", len(insertResp.Result.InsertedRows))
	}

	updateRec := request(`{"sql":"UPDATE users SET age=:age WHERE id IN :ids","paramMap":{"age":30,"ids":[1,2]}}`)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("batch update status = %d, want %d, body=%s", updateRec.Code, http.StatusOK, updateRec.Body.String())
	}

	var updateResp SQLExecuteResponse
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("unmarshal update response: %v", err)
	}
	if !updateResp.Success || updateResp.Result == nil {
		t.Fatalf("unexpected update response: %+v", updateResp)
	}
	if len(updateResp.Result.UpdatedRows) != 2 {
		t.Fatalf("updated rows = %d, want 2", len(updateResp.Result.UpdatedRows))
	}
}

func TestHTTPServer_SQLExecute_DeleteWithoutWhereRejected(t *testing.T) {
	server := New.HTTP(t.TempDir())

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginRec.Code, http.StatusOK)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, DefaultSQLExecutePath, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		return rec
	}

	createRec := request(`{"sql":"CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED)"}`)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d, body=%s", createRec.Code, http.StatusOK, createRec.Body.String())
	}

	insertRec := request(`{"sql":"INSERT INTO users (username) VALUES ('alice')"}`)
	if insertRec.Code != http.StatusOK {
		t.Fatalf("insert status = %d, want %d, body=%s", insertRec.Code, http.StatusOK, insertRec.Body.String())
	}

	deleteRec := request(`{"sql":"DELETE FROM users"}`)
	if deleteRec.Code != http.StatusBadRequest {
		t.Fatalf("delete without where status = %d, want %d, body=%s", deleteRec.Code, http.StatusBadRequest, deleteRec.Body.String())
	}
}

func TestHTTPServer_SQLExecute_ParamListAndParamMapCanBeMixed(t *testing.T) {
	server := New.HTTP(t.TempDir())

	loginBody, _ := json.Marshal(LoginRequest{Username: "sdb", Password: "simpleDB"})
	loginReq := httptest.NewRequest(http.MethodPost, DefaultLoginPath, bytes.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	server.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginRec.Code, http.StatusOK)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if loginResp.Token == nil {
		t.Fatal("expected token in login response")
	}

	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, DefaultSQLExecutePath, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+loginResp.Token.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		return rec
	}

	createRec := request(`{"sql":"CREATE TABLE users (id int PRIMARY KEY AUTO_INCREMENT, username string UNIQUE REQUIRED, age int DEFAULT 0)"}`)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create status = %d, want %d, body=%s", createRec.Code, http.StatusOK, createRec.Body.String())
	}

	insertRec := request(`{"sql":"INSERT INTO users (username, age) VALUES (?, :age)","paramList":["alice"],"paramMap":{"age":20}}`)
	if insertRec.Code != http.StatusOK {
		t.Fatalf("insert status = %d, want %d, body=%s", insertRec.Code, http.StatusOK, insertRec.Body.String())
	}

	selectRec := request(`{"sql":"SELECT * FROM users WHERE username=? AND age=:age","paramList":["alice"],"paramMap":{"age":20}}`)
	if selectRec.Code != http.StatusOK {
		t.Fatalf("select status = %d, want %d, body=%s", selectRec.Code, http.StatusOK, selectRec.Body.String())
	}

	var sqlResp SQLExecuteResponse
	if err := json.Unmarshal(selectRec.Body.Bytes(), &sqlResp); err != nil {
		t.Fatalf("unmarshal sql response: %v", err)
	}
	if !sqlResp.Success || sqlResp.Result == nil {
		t.Fatalf("unexpected sql response: %+v", sqlResp)
	}
	if len(sqlResp.Result.Rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(sqlResp.Result.Rows))
	}
	if got := sqlResp.Result.Rows[0]["username"]; got != "alice" {
		t.Fatalf("username = %v, want alice", got)
	}
	switch got := sqlResp.Result.Rows[0]["age"].(type) {
	case int:
		if got != 20 {
			t.Fatalf("age = %v, want 20", got)
		}
	case int64:
		if got != 20 {
			t.Fatalf("age = %v, want 20", got)
		}
	case float64:
		if got != 20 {
			t.Fatalf("age = %v, want 20", got)
		}
	default:
		t.Fatalf("unexpected age type: %T (%v)", got, got)
	}

	legacyRec := request(`{"sql":"SELECT * FROM users WHERE username=:username","params":{"username":"alice"}}`)
	if legacyRec.Code != http.StatusOK {
		t.Fatalf("legacy params status = %d, want %d, body=%s", legacyRec.Code, http.StatusOK, legacyRec.Body.String())
	}
}
