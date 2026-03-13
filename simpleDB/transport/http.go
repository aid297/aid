package transport

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aid297/aid/simpleDB/api"
	"github.com/aid297/aid/simpleDB/driver"
	"github.com/aid297/aid/simpleDB/kernal"
	"github.com/gin-gonic/gin"
)

const DefaultLoginPath = "/auth/login"
const DefaultRegisterPath = "/auth/register"
const DefaultRefreshPath = "/auth/refresh"
const DefaultLogoutPath = "/auth/logout"
const DefaultActivatePath = "/auth/activate"
const DefaultDeactivatePath = "/auth/deactivate"
const DefaultAssignRolesPath = "/auth/assign-roles"
const DefaultAssignRolePermissionsPath = "/auth/assign-role-permissions"
const DefaultInitSDBPasswordPath = "/auth/init-sdb-password"
const DefaultSQLExecutePath = "/sql/execute"
const DefaultSQLGrantPath = "/sql/grant"
const DefaultSQLRevokePath = "/sql/revoke"
const SuperAdminRoleCode = "super_admin"

const ContextUserKey = "simpledb.transport.user"

var New app

type app struct{}

type Authenticator interface {
	Authenticate(database, username, password string) (*driver.AuthenticatedUser, error)
	RegisterUser(database, username, password, displayName string) (*driver.AuthenticatedUser, error)
	ActivateUser(database, username string) (*driver.AuthenticatedUser, error)
	DeactivateUser(database, username string) (*driver.AuthenticatedUser, error)
	AssignRoles(database, username string, roleCodes []string) (*driver.AuthenticatedUser, error)
	AssignRolePermissions(database, roleCode string, permissionCodes []string) error
	InitSDBPassword(database string) error
	BindUserDatabase(database string, approver *driver.AuthenticatedUser, username string) error
	RevokeUserDatabase(database string, approver *driver.AuthenticatedUser, username string) error
}

type HTTPServer struct {
	Database                 string
	LoginPath                string
	RegisterPath             string
	RefreshPath              string
	LogoutPath               string
	ActivatePath             string
	DeactivatePath           string
	AssignRolePath           string
	AssignRolePermissionPath string
	InitSDBPasswordPath      string
	SQLExecutePath           string
	SQLGrantPath             string
	SQLRevokePath            string
	SQLAllowedOps            map[string]struct{} // nil = 不限制
	LimitEnabled             bool
	LimitRequests            int
	LimitWindow              time.Duration
	LimitNoTokenPaths        map[string]struct{}
	limitMu                  sync.Mutex
	limitBuckets             map[string]*limitBucket
	InitPassword             string
	initPasswordMu           sync.RWMutex
	initPasswordRotator      func() (string, error)
	TokenTTL                 time.Duration
	TokenSecret              string
	authenticator            Authenticator
	engine                   *gin.Engine
	tokenManager             *TokenManager
	sqlEngine                *api.Engine
	sqlEngineDBAttrs         []kernal.SchemaAttributer
	WebSocketEnabled         bool
	WebSocketRoute           string
	WebSocketHeartbeat       time.Duration
	WebSocketWriteTimeout    time.Duration
	WebSocketReadTimeout     time.Duration
	WebSocketExecutionTimeout time.Duration
}

type limitBucket struct {
	WindowStart time.Time
	Count       int
}

type Option func(*HTTPServer)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type ActivateRequest struct {
	Username string `json:"username"`
}

type AssignRolesRequest struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type AssignRolePermissionsRequest struct {
	RoleCode    string   `json:"roleCode"`
	Permissions []string `json:"permissions"`
}

type InitSDBPasswordRequest struct {
	Password string `json:"password"`
}

type SQLExecuteRequest struct {
	SQL       string         `json:"sql"`
	ParamMap  map[string]any `json:"paramMap"`
	ParamList []any          `json:"paramList"`
	Params    map[string]any `json:"params,omitempty"`
}

func (r SQLExecuteRequest) mergedParamMap() map[string]any {
	if len(r.ParamMap) == 0 && len(r.Params) == 0 {
		return nil
	}
	merged := make(map[string]any, len(r.Params)+len(r.ParamMap))
	for key, value := range r.Params {
		merged[key] = value
	}
	for key, value := range r.ParamMap {
		merged[key] = value
	}
	return merged
}

type SQLExecuteResponse struct {
	Success bool            `json:"success"`
	Result  *api.ExecResult `json:"result,omitempty"`
	Error   *ErrorBody      `json:"error,omitempty"`
}

type LoginResponse struct {
	Success bool                      `json:"success"`
	User    *driver.AuthenticatedUser `json:"user,omitempty"`
	Token   *TokenResponse            `json:"token,omitempty"`
	Error   *ErrorBody                `json:"error,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (*app) HTTP(database string, opts ...Option) *HTTPServer {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.HandleMethodNotAllowed = true

	server := &HTTPServer{
		Database:                 strings.TrimSpace(database),
		LoginPath:                DefaultLoginPath,
		RegisterPath:             DefaultRegisterPath,
		RefreshPath:              DefaultRefreshPath,
		LogoutPath:               DefaultLogoutPath,
		ActivatePath:             DefaultActivatePath,
		DeactivatePath:           DefaultDeactivatePath,
		AssignRolePath:           DefaultAssignRolesPath,
		AssignRolePermissionPath: DefaultAssignRolePermissionsPath,
		InitSDBPasswordPath:      DefaultInitSDBPasswordPath,
		SQLExecutePath:           DefaultSQLExecutePath,
		SQLGrantPath:             DefaultSQLGrantPath,
		SQLRevokePath:            DefaultSQLRevokePath,
		LimitEnabled:             false,
		LimitRequests:            60,
		LimitWindow:              time.Minute,
		TokenTTL:                 12 * time.Hour,
		WebSocketEnabled:         false,
		WebSocketRoute:           "/ws",
		WebSocketHeartbeat:       10 * time.Second,
		WebSocketWriteTimeout:    10 * time.Second,
		WebSocketReadTimeout:     60 * time.Second,
		WebSocketExecutionTimeout: 30 * time.Second,
		authenticator:            &driver.New,
		engine:                   engine,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(server)
		}
	}
	if server.LoginPath == "" {
		server.LoginPath = DefaultLoginPath
	}
	if server.RefreshPath == "" {
		server.RefreshPath = DefaultRefreshPath
	}
	if server.LogoutPath == "" {
		server.LogoutPath = DefaultLogoutPath
	}
	if server.RegisterPath == "" {
		server.RegisterPath = DefaultRegisterPath
	}
	if server.ActivatePath == "" {
		server.ActivatePath = DefaultActivatePath
	}
	if server.DeactivatePath == "" {
		server.DeactivatePath = DefaultDeactivatePath
	}
	if server.AssignRolePath == "" {
		server.AssignRolePath = DefaultAssignRolesPath
	}
	if server.AssignRolePermissionPath == "" {
		server.AssignRolePermissionPath = DefaultAssignRolePermissionsPath
	}
	if server.InitSDBPasswordPath == "" {
		server.InitSDBPasswordPath = DefaultInitSDBPasswordPath
	}
	if server.SQLExecutePath == "" {
		server.SQLExecutePath = DefaultSQLExecutePath
	}
	if server.SQLGrantPath == "" {
		server.SQLGrantPath = DefaultSQLGrantPath
	}
	if server.SQLRevokePath == "" {
		server.SQLRevokePath = DefaultSQLRevokePath
	}
	if server.LimitRequests <= 0 {
		server.LimitRequests = 60
	}
	if server.LimitWindow <= 0 {
		server.LimitWindow = time.Minute
	}
	if server.limitBuckets == nil {
		server.limitBuckets = make(map[string]*limitBucket)
	}
	server.sqlEngine = api.NewEngine(server.Database, api.BackendDriver).WithDBAttrs(server.sqlEngineDBAttrs...)
	server.tokenManager = NewTokenManager(server.Database, server.TokenSecret, server.TokenTTL).
		WithStore(NewDBTokenStore())
	server.registerRoutes()
	return server
}

func WithLoginPath(path string) Option {
	return func(server *HTTPServer) {
		server.LoginPath = normalizePath(path, DefaultLoginPath)
	}
}

func WithRefreshPath(path string) Option {
	return func(server *HTTPServer) {
		server.RefreshPath = normalizePath(path, DefaultRefreshPath)
	}
}

func WithRegisterPath(path string) Option {
	return func(server *HTTPServer) {
		server.RegisterPath = normalizePath(path, DefaultRegisterPath)
	}
}

func WithLogoutPath(path string) Option {
	return func(server *HTTPServer) {
		server.LogoutPath = normalizePath(path, DefaultLogoutPath)
	}
}

func WithActivatePath(path string) Option {
	return func(server *HTTPServer) {
		server.ActivatePath = normalizePath(path, DefaultActivatePath)
	}
}

func WithDeactivatePath(path string) Option {
	return func(server *HTTPServer) {
		server.DeactivatePath = normalizePath(path, DefaultDeactivatePath)
	}
}

func WithAssignRolePath(path string) Option {
	return func(server *HTTPServer) {
		server.AssignRolePath = normalizePath(path, DefaultAssignRolesPath)
	}
}

func WithAssignRolePermissionPath(path string) Option {
	return func(server *HTTPServer) {
		server.AssignRolePermissionPath = normalizePath(path, DefaultAssignRolePermissionsPath)
	}
}

func WithInitSDBPasswordPath(path string) Option {
	return func(server *HTTPServer) {
		server.InitSDBPasswordPath = normalizePath(path, DefaultInitSDBPasswordPath)
	}
}

func WithSQLExecutePath(path string) Option {
	return func(server *HTTPServer) {
		server.SQLExecutePath = normalizePath(path, DefaultSQLExecutePath)
	}
}

func WithSQLGrantPath(path string) Option {
	return func(server *HTTPServer) {
		server.SQLGrantPath = normalizePath(path, DefaultSQLGrantPath)
	}
}

func WithSQLRevokePath(path string) Option {
	return func(server *HTTPServer) {
		server.SQLRevokePath = normalizePath(path, DefaultSQLRevokePath)
	}
}

func WithSQLAllowedOps(ops []string) Option {
	return func(server *HTTPServer) {
		if len(ops) == 0 {
			server.SQLAllowedOps = nil
			return
		}
		m := make(map[string]struct{}, len(ops))
		for _, op := range ops {
			m[strings.ToLower(strings.TrimSpace(op))] = struct{}{}
		}
		server.SQLAllowedOps = m
	}
}

func WithTokenRateLimit(enabled bool, requests int, window time.Duration, noTokenPaths []string) Option {
	return func(server *HTTPServer) {
		server.LimitEnabled = enabled
		if requests > 0 {
			server.LimitRequests = requests
		}
		if window > 0 {
			server.LimitWindow = window
		}
		if len(noTokenPaths) == 0 {
			server.LimitNoTokenPaths = nil
			return
		}
		m := make(map[string]struct{}, len(noTokenPaths))
		for _, path := range noTokenPaths {
			m[normalizePath(path, path)] = struct{}{}
		}
		server.LimitNoTokenPaths = m
	}
}

func WithInitPassword(password string) Option {
	return func(server *HTTPServer) {
		server.initPasswordMu.Lock()
		server.InitPassword = strings.TrimSpace(password)
		server.initPasswordMu.Unlock()
	}
}

func WithInitPasswordRotator(rotator func() (string, error)) Option {
	return func(server *HTTPServer) {
		server.initPasswordRotator = rotator
	}
}

func WithAuthenticator(authenticator Authenticator) Option {
	return func(server *HTTPServer) {
		if authenticator != nil {
			server.authenticator = authenticator
		}
	}
}

func WithTokenSecret(secret string) Option {
	return func(server *HTTPServer) {
		server.TokenSecret = strings.TrimSpace(secret)
	}
}

func WithTokenTTL(ttl time.Duration) Option {
	return func(server *HTTPServer) {
		if ttl > 0 {
			server.TokenTTL = ttl
		}
	}
}

func WithSQLEngineDBAttrs(attrs ...kernal.SchemaAttributer) Option {
	return func(server *HTTPServer) {
		if len(attrs) == 0 {
			server.sqlEngineDBAttrs = nil
			return
		}
		server.sqlEngineDBAttrs = append([]kernal.SchemaAttributer(nil), attrs...)
	}
}

func WithWebSocketConfig(enabled bool, route string, heartbeat, writeTimeout, readTimeout, execTimeout time.Duration) Option {
	return func(s *HTTPServer) {
		s.WebSocketEnabled = enabled
		s.WebSocketRoute = normalizePath(route, "/ws")
		s.WebSocketHeartbeat = heartbeat
		s.WebSocketWriteTimeout = writeTimeout
		s.WebSocketReadTimeout = readTimeout
		s.WebSocketExecutionTimeout = execTimeout
	}
}

func (s *HTTPServer) Handler() http.Handler {
	return s.engine
}

func (s *HTTPServer) Engine() *gin.Engine {
	return s.engine
}

func (s *HTTPServer) Close() error {
	if s.sqlEngine == nil {
		return nil
	}
	return s.sqlEngine.Close()
}

func (s *HTTPServer) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *HTTPServer) Serve(listener net.Listener) error {
	return s.engine.RunListener(listener)
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}

func (s *HTTPServer) ParseToken(token string) (*TokenClaims, error) {
	return s.tokenManager.Parse(token)
}

func (s *HTTPServer) RevokeToken(token string) error {
	return s.tokenManager.Revoke(token)
}

func (s *HTTPServer) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, claims, err := s.readBearerToken(ctx)
		if err != nil {
			writeUnauthorized(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set("simpledb.transport.token", token)
		ctx.Set(ContextUserKey, claims)
		ctx.Next()
	}
}

func (s *HTTPServer) RequireRoles(roles ...string) gin.HandlerFunc {
	required := normalizeValues(roles)
	return func(ctx *gin.Context) {
		claims, ok := UserFromContext(ctx)
		if !ok {
			writeJSON(ctx, http.StatusUnauthorized, LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "missing auth context"}})
			ctx.Abort()
			return
		}
		if hasRole(claims.Roles, SuperAdminRoleCode) || len(required) == 0 || hasAnyValue(claims.Roles, required) {
			ctx.Next()
			return
		}
		writeJSON(ctx, http.StatusForbidden, LoginResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: "missing required role"}})
		ctx.Abort()
	}
}

func (s *HTTPServer) RequirePermissions(permissions ...string) gin.HandlerFunc {
	required := normalizeValues(permissions)
	return func(ctx *gin.Context) {
		claims, ok := UserFromContext(ctx)
		if !ok {
			writeJSON(ctx, http.StatusUnauthorized, LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "missing auth context"}})
			ctx.Abort()
			return
		}
		if hasRole(claims.Roles, SuperAdminRoleCode) || len(required) == 0 || containsAllValues(claims.Permissions, required) {
			ctx.Next()
			return
		}
		writeJSON(ctx, http.StatusForbidden, LoginResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: "missing required permission"}})
		ctx.Abort()
	}
}

func (s *HTTPServer) registerRoutes() {
	if s.LimitEnabled {
		s.engine.Use(s.tokenRateLimitMiddleware())
	}
	s.engine.NoMethod(func(ctx *gin.Context) {
		writeJSON(ctx, http.StatusMethodNotAllowed, LoginResponse{
			Success: false,
			Error:   &ErrorBody{Code: "method_not_allowed", Message: "only POST is allowed"},
		})
	})

	s.engine.POST(s.LoginPath, s.handleLogin)
	s.engine.POST(s.RegisterPath, s.handleRegister)
	s.engine.POST(s.RefreshPath, s.handleRefresh)
	s.engine.POST(s.LogoutPath, s.handleLogout)
	s.engine.POST(s.ActivatePath, s.AuthMiddleware(), s.RequireRoles(SuperAdminRoleCode), s.handleActivate)
	s.engine.POST(s.DeactivatePath, s.AuthMiddleware(), s.RequireRoles(SuperAdminRoleCode), s.handleDeactivate)
	s.engine.POST(s.AssignRolePath, s.AuthMiddleware(), s.RequireRoles(SuperAdminRoleCode), s.handleAssignRoles)
	s.engine.POST(s.AssignRolePermissionPath, s.AuthMiddleware(), s.RequireRoles(SuperAdminRoleCode), s.handleAssignRolePermissions)
	s.engine.POST(s.InitSDBPasswordPath, s.handleInitSDBPassword)
	s.engine.POST(s.SQLExecutePath, s.AuthMiddleware(), s.handleSQLExecute)
	s.engine.POST(s.SQLGrantPath, s.AuthMiddleware(), s.handleSQLGrant)
	s.engine.POST(s.SQLRevokePath, s.AuthMiddleware(), s.handleSQLRevoke)

	if s.WebSocketEnabled && s.WebSocketRoute != "" {
		s.engine.GET(s.WebSocketRoute, s.handleWebSocket)
	}
}

func (s *HTTPServer) tokenRateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !s.LimitEnabled || s.LimitRequests <= 0 || s.LimitWindow <= 0 {
			ctx.Next()
			return
		}

		key, limited := s.resolveLimitKey(ctx)
		if !limited {
			ctx.Next()
			return
		}

		now := time.Now()
		allowed, retryAfter := s.allowLimitKey(key, now)
		if allowed {
			ctx.Next()
			return
		}

		seconds := int(math.Ceil(retryAfter.Seconds()))
		if seconds < 1 {
			seconds = 1
		}
		ctx.Header("Retry-After", strconv.Itoa(seconds))
		ctx.JSON(http.StatusTooManyRequests, LoginResponse{Success: false, Error: &ErrorBody{Code: "too_many_requests", Message: "rate limit exceeded"}})
		ctx.Abort()
	}
}

func (s *HTTPServer) resolveLimitKey(ctx *gin.Context) (string, bool) {
	path := strings.TrimSpace(ctx.FullPath())
	if path == "" {
		path = normalizePath(ctx.Request.URL.Path, ctx.Request.URL.Path)
	}
	if path == "" {
		return "", false
	}

	if _, ok := s.LimitNoTokenPaths[path]; ok {
		return "anon:" + path + ":" + strings.TrimSpace(ctx.ClientIP()), true
	}

	token := readBearerTokenRaw(ctx.GetHeader("Authorization"))
	if token == "" {
		return "", false
	}
	return "token:" + token, true
}

func readBearerTokenRaw(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func (s *HTTPServer) allowLimitKey(key string, now time.Time) (bool, time.Duration) {
	s.limitMu.Lock()
	defer s.limitMu.Unlock()

	bucket, ok := s.limitBuckets[key]
	if !ok {
		s.limitBuckets[key] = &limitBucket{WindowStart: now, Count: 1}
		return true, 0
	}

	if now.Sub(bucket.WindowStart) >= s.LimitWindow {
		bucket.WindowStart = now
		bucket.Count = 1
		return true, 0
	}

	if bucket.Count >= s.LimitRequests {
		retryAfter := s.LimitWindow - now.Sub(bucket.WindowStart)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	bucket.Count++
	return true, 0
}

func (s *HTTPServer) handleRegister(ctx *gin.Context) {
	if s.Database == "" {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{Success: false, Error: &ErrorBody{Code: "database_required", Message: "database is required"}})
		return
	}

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}

	user, err := s.authenticator.RegisterUser(s.Database, req.Username, req.Password, req.DisplayName)
	if err != nil {
		status, code, message := mapDriverError(err, "register failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	writeJSON(ctx, http.StatusCreated, LoginResponse{Success: true, User: user})
}

func (s *HTTPServer) handleLogin(ctx *gin.Context) {
	if s.Database == "" {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{
			Success: false,
			Error:   &ErrorBody{Code: "database_required", Message: "database is required"},
		})
		return
	}

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   &ErrorBody{Code: "bad_request", Message: "invalid request body"},
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   &ErrorBody{Code: "bad_request", Message: "username and password are required"},
		})
		return
	}

	user, err := s.authenticator.Authenticate(s.Database, req.Username, req.Password)
	if err != nil {
		status, code, message := mapDriverError(err, "login failed")
		if status == http.StatusUnauthorized {
			message = "username or password is invalid"
		}
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	token, err := s.tokenManager.Issue(user)
	if err != nil {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{
			Success: false,
			Error:   &ErrorBody{Code: "token_issue_failed", Message: "failed to issue access token"},
		})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{
		Success: true,
		User:    user,
		Token: &TokenResponse{
			AccessToken: token.AccessToken,
			TokenType:   token.TokenType,
			ExpiresAt:   token.ExpiresAt,
		},
	})
}

func (s *HTTPServer) handleRefresh(ctx *gin.Context) {
	token, _, err := s.readBearerToken(ctx)
	if err != nil {
		writeUnauthorized(ctx, err)
		return
	}
	issued, claims, err := s.tokenManager.Refresh(token)
	if err != nil {
		writeUnauthorized(ctx, err)
		return
	}
	writeJSON(ctx, http.StatusOK, LoginResponse{
		Success: true,
		User: &driver.AuthenticatedUser{
			ID:          claims.Subject,
			Username:    claims.Username,
			DisplayName: claims.DisplayName,
			Status:      claims.Status,
			IsAdmin:     claims.IsAdmin,
			Roles:       append([]string(nil), claims.Roles...),
			Permissions: append([]string(nil), claims.Permissions...),
		},
		Token: &TokenResponse{AccessToken: issued.AccessToken, TokenType: issued.TokenType, ExpiresAt: issued.ExpiresAt},
	})
}

func (s *HTTPServer) handleLogout(ctx *gin.Context) {
	token, _, err := s.readBearerToken(ctx)
	if err != nil {
		writeUnauthorized(ctx, err)
		return
	}
	if err = s.tokenManager.Revoke(token); err != nil {
		writeUnauthorized(ctx, err)
		return
	}
	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true})
}

func (s *HTTPServer) handleActivate(ctx *gin.Context) {
	var req ActivateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "username is required"}})
		return
	}

	user, err := s.authenticator.ActivateUser(s.Database, req.Username)
	if err != nil {
		status, code, message := mapDriverError(err, "activate user failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true, User: user})
}

func (s *HTTPServer) handleDeactivate(ctx *gin.Context) {
	var req ActivateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "username is required"}})
		return
	}

	user, err := s.authenticator.DeactivateUser(s.Database, req.Username)
	if err != nil {
		status, code, message := mapDriverError(err, "deactivate user failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true, User: user})
}

func (s *HTTPServer) handleAssignRoles(ctx *gin.Context) {
	var req AssignRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || len(req.Roles) == 0 {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "username and roles are required"}})
		return
	}

	user, err := s.authenticator.AssignRoles(s.Database, req.Username, req.Roles)
	if err != nil {
		status, code, message := mapDriverError(err, "assign roles failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true, User: user})
}

func (s *HTTPServer) handleAssignRolePermissions(ctx *gin.Context) {
	var req AssignRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	req.RoleCode = strings.TrimSpace(req.RoleCode)
	if req.RoleCode == "" || len(req.Permissions) == 0 {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "roleCode and permissions are required"}})
		return
	}

	err := s.authenticator.AssignRolePermissions(s.Database, req.RoleCode, req.Permissions)
	if err != nil {
		status, code, message := mapDriverError(err, "assign role permissions failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true})
}

func (s *HTTPServer) handleInitSDBPassword(ctx *gin.Context) {
	if s.Database == "" {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{Success: false, Error: &ErrorBody{Code: "database_required", Message: "database is required"}})
		return
	}
	currentInitPassword := s.getInitPassword()
	if strings.TrimSpace(currentInitPassword) == "" {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{Success: false, Error: &ErrorBody{Code: "init_password_required", Message: "init password is not configured"}})
		return
	}

	var req InitSDBPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeJSON(ctx, http.StatusBadRequest, LoginResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}

	if req.Password != currentInitPassword {
		writeJSON(ctx, http.StatusForbidden, LoginResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: "init password is invalid"}})
		return
	}

	err := s.authenticator.InitSDBPassword(s.Database)
	if err != nil {
		status, code, message := mapDriverError(err, "init sdb password failed")
		writeJSON(ctx, status, LoginResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	if err = s.rotateInitPassword(); err != nil {
		writeJSON(ctx, http.StatusInternalServerError, LoginResponse{Success: false, Error: &ErrorBody{Code: "init_password_rotate_failed", Message: "failed to rotate init password"}})
		return
	}

	writeJSON(ctx, http.StatusOK, LoginResponse{Success: true})
}

func (s *HTTPServer) handleSQLExecute(ctx *gin.Context) {
	if s.Database == "" {
		ctx.JSON(http.StatusInternalServerError, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "database_required", Message: "database is required"}})
		return
	}

	var req SQLExecuteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	if strings.TrimSpace(req.SQL) == "" {
		ctx.JSON(http.StatusBadRequest, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "sql is required"}})
		return
	}

	boundSQL, err := bindSQLParams(req.SQL, req.mergedParamMap(), req.ParamList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: err.Error()}})
		return
	}

	claims, ok := UserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "missing auth context"}})
		return
	}

	actor := &driver.AuthenticatedUser{
		ID:          claims.Subject,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		Status:      claims.Status,
		IsAdmin:     claims.IsAdmin,
		Roles:       append([]string(nil), claims.Roles...),
		Permissions: append([]string(nil), claims.Permissions...),
	}

	engine := api.NewEngine(s.Database, api.BackendDriver).
		WithSharedCacheFrom(s.sqlEngine).
		WithActor(actor)

	// 白名单校验：解析语句类型，若配置了 SQLAllowedOps 则拒绝不在列表内的操作
	if len(s.SQLAllowedOps) > 0 {
		stmt, parseErr := engine.Parse(boundSQL)
		if parseErr != nil {
			ctx.JSON(http.StatusBadRequest, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: parseErr.Error()}})
			return
		}
		opKey := strings.ToLower(strings.TrimPrefix(string(stmt.Type()), ""))
		if _, allowed := s.SQLAllowedOps[opKey]; !allowed {
			ctx.JSON(http.StatusForbidden, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: "operation not allowed: " + opKey}})
			return
		}
	}

	result, err := engine.Execute(boundSQL)
	if err != nil {
		if errors.Is(err, api.ErrSystemTableAccessDenied) {
			ctx.JSON(http.StatusForbidden, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: err.Error()}})
			return
		}

		status, code, message := mapDriverError(err, "sql execute failed")
		if status == http.StatusInternalServerError && code == "internal_error" {
			status = http.StatusBadRequest
			code = "bad_request"
			message = err.Error()
		}
		ctx.JSON(status, SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}

	ctx.JSON(http.StatusOK, SQLExecuteResponse{Success: true, Result: &result})
}

// SQLGrantRequest is the request body for POST /sql/grant.
type SQLGrantRequest struct {
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	// legacy aliases
	Table   string `json:"table,omitempty"`
	Grantee string `json:"grantee,omitempty"`
}

// SQLGrantResponse is the response body for POST /sql/grant.
type SQLGrantResponse struct {
	Success bool       `json:"success"`
	Error   *ErrorBody `json:"error,omitempty"`
}

func (s *HTTPServer) handleSQLGrant(ctx *gin.Context) {
	if s.Database == "" {
		ctx.JSON(http.StatusInternalServerError, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "database_required", Message: "database is required"}})
		return
	}

	claims, ok := ctx.MustGet(ContextUserKey).(*TokenClaims)
	if !ok || claims == nil {
		ctx.JSON(http.StatusUnauthorized, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "not authenticated"}})
		return
	}
	approver := &driver.AuthenticatedUser{
		ID:          claims.Subject,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}

	var req SQLGrantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	targetDatabase := strings.TrimSpace(req.Database)
	if targetDatabase == "" {
		targetDatabase = strings.TrimSpace(req.Table)
	}
	if targetDatabase == "" {
		targetDatabase = s.Database
	}
	targetUsername := strings.TrimSpace(req.Username)
	if targetUsername == "" {
		targetUsername = strings.TrimSpace(req.Grantee)
	}
	if targetUsername == "" {
		ctx.JSON(http.StatusBadRequest, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "username is required"}})
		return
	}

	err := s.authenticator.BindUserDatabase(targetDatabase, approver, targetUsername)
	if err != nil {
		status, code, message := mapDriverError(err, "database bind failed")
		ctx.JSON(status, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}
	ctx.JSON(http.StatusOK, SQLGrantResponse{Success: true})
}

func (s *HTTPServer) handleSQLRevoke(ctx *gin.Context) {
	if s.Database == "" {
		ctx.JSON(http.StatusInternalServerError, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "database_required", Message: "database is required"}})
		return
	}

	claims, ok := ctx.MustGet(ContextUserKey).(*TokenClaims)
	if !ok || claims == nil {
		ctx.JSON(http.StatusUnauthorized, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "not authenticated"}})
		return
	}
	approver := &driver.AuthenticatedUser{
		ID:          claims.Subject,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}

	var req SQLGrantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid request body"}})
		return
	}
	targetDatabase := strings.TrimSpace(req.Database)
	if targetDatabase == "" {
		targetDatabase = strings.TrimSpace(req.Table)
	}
	if targetDatabase == "" {
		targetDatabase = s.Database
	}
	targetUsername := strings.TrimSpace(req.Username)
	if targetUsername == "" {
		targetUsername = strings.TrimSpace(req.Grantee)
	}
	if targetUsername == "" {
		ctx.JSON(http.StatusBadRequest, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "username is required"}})
		return
	}

	err := s.authenticator.RevokeUserDatabase(targetDatabase, approver, targetUsername)
	if err != nil {
		status, code, message := mapDriverError(err, "database revoke failed")
		ctx.JSON(status, SQLGrantResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}})
		return
	}
	ctx.JSON(http.StatusOK, SQLGrantResponse{Success: true})
}

func bindSQLParams(sql string, paramMap map[string]any, paramList []any) (string, error) {
	if strings.TrimSpace(sql) == "" || (len(paramMap) == 0 && len(paramList) == 0) {
		return sql, nil
	}

	var builder strings.Builder
	builder.Grow(len(sql) + 32)

	inSingleQuote := false
	inDoubleQuote := false
	paramIndex := 0

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch ch {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
			builder.WriteByte(ch)
			continue
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
			builder.WriteByte(ch)
			continue
		}

		if !inSingleQuote && !inDoubleQuote && ch == '?' {
			if paramIndex >= len(paramList) {
				return "", fmt.Errorf("missing sql paramList value at index %d", paramIndex)
			}
			literal, err := toSQLLiteral(paramList[paramIndex])
			if err != nil {
				return "", fmt.Errorf("invalid sql paramList[%d]: %w", paramIndex, err)
			}
			builder.WriteString(literal)
			paramIndex++
			continue
		}

		if !inSingleQuote && !inDoubleQuote && ch == ':' && i+1 < len(sql) && isIdentStart(sql[i+1]) {
			j := i + 2
			for j < len(sql) && isIdentPart(sql[j]) {
				j++
			}
			name := sql[i+1 : j]
			value, ok := paramMap[name]
			if !ok {
				return "", fmt.Errorf("missing sql paramMap key: %s", name)
			}
			literal, err := toSQLLiteral(value)
			if err != nil {
				return "", fmt.Errorf("invalid sql paramMap[%s]: %w", name, err)
			}
			builder.WriteString(literal)
			i = j - 1
			continue
		}

		builder.WriteByte(ch)
	}

	if paramIndex < len(paramList) {
		return "", fmt.Errorf("unused sql paramList values: %d", len(paramList)-paramIndex)
	}

	return builder.String(), nil
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

func toSQLLiteral(value any) (string, error) {
	if value == nil {
		return "null", nil
	}

	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'", nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case []byte:
		return "'" + strings.ReplaceAll(string(v), "'", "''") + "'", nil
	}

	rv := reflect.ValueOf(value)
	if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
		items := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			item, err := toSQLLiteral(rv.Index(i).Interface())
			if err != nil {
				return "", err
			}
			items = append(items, item)
		}
		return "(" + strings.Join(items, ",") + ")", nil
	}

	return "", fmt.Errorf("unsupported param type: %T", value)
}

func (s *HTTPServer) getInitPassword() string {
	s.initPasswordMu.RLock()
	defer s.initPasswordMu.RUnlock()
	return strings.TrimSpace(s.InitPassword)
}

func (s *HTTPServer) setInitPassword(password string) {
	s.initPasswordMu.Lock()
	defer s.initPasswordMu.Unlock()
	s.InitPassword = strings.TrimSpace(password)
}

func (s *HTTPServer) rotateInitPassword() error {
	if s.initPasswordRotator != nil {
		newPassword, err := s.initPasswordRotator()
		if err != nil {
			return err
		}
		s.setInitPassword(newPassword)
		return nil
	}

	generated, err := generateRandomString(24)
	if err != nil {
		return err
	}
	s.setInitPassword(generated)
	return nil
}

func generateRandomString(length int) (string, error) {
	if length <= 0 {
		length = 24
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	text := base64.RawURLEncoding.EncodeToString(buf)
	if len(text) >= length {
		return text[:length], nil
	}
	return text, nil
}

func mapDriverError(err error, fallback string) (int, string, string) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := fallback

	var driverErr *driver.DriverError
	if !errors.As(err, &driverErr) {
		return status, code, message
	}

	switch driverErr.Code {
	case driver.ErrorCodeInvalidArgument:
		return http.StatusBadRequest, "bad_request", driverErr.Err.Error()
	case driver.ErrorCodeUnauthorized:
		return http.StatusUnauthorized, "unauthorized", driverErr.Err.Error()
	case driver.ErrorCodeConflict:
		return http.StatusConflict, "conflict", driverErr.Err.Error()
	case driver.ErrorCodeNotFound:
		return http.StatusNotFound, "not_found", driverErr.Err.Error()
	default:
		if err != nil {
			return status, code, err.Error()
		}
		return status, code, message
	}
}

func writeJSON(ctx *gin.Context, status int, payload LoginResponse) {
	ctx.JSON(status, payload)
}

func writeUnauthorized(ctx *gin.Context, err error) {
	message := "invalid or expired token"
	if errors.Is(err, ErrInvalidToken) {
		message = "invalid token"
	}
	writeJSON(ctx, http.StatusUnauthorized, LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: message}})
}

func UserFromContext(ctx *gin.Context) (*TokenClaims, bool) {
	value, ok := ctx.Get(ContextUserKey)
	if !ok || value == nil {
		return nil, false
	}
	claims, ok := value.(*TokenClaims)
	return claims, ok
}

func (s *HTTPServer) readBearerToken(ctx *gin.Context) (string, *TokenClaims, error) {
	header := strings.TrimSpace(ctx.GetHeader("Authorization"))
	if header == "" {
		return "", nil, ErrInvalidToken
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", nil, ErrInvalidToken
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	claims, err := s.tokenManager.Parse(token)
	if err != nil {
		return token, nil, err
	}
	return token, claims, nil
}

func normalizePath(path, fallback string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fallback
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return trimmed
}

func normalizeValues(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func hasAnyValue(values, required []string) bool {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	for _, value := range required {
		if _, ok := set[value]; ok {
			return true
		}
	}
	return false
}

func hasRole(roles []string, roleCode string) bool {
	for _, role := range roles {
		if role == roleCode {
			return true
		}
	}
	return false
}

func containsAllValues(values, required []string) bool {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	for _, value := range required {
		if _, ok := set[value]; !ok {
			return false
		}
	}
	return true
}
