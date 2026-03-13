package transport

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aid297/aid/simpleDB/api"
	"github.com/aid297/aid/simpleDB/driver"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsRequest struct {
	Route string `json:"route"`
	Token string `json:"token,omitempty"` // For login
	// Embed SQLExecuteRequest fields
	SQL       string       `json:"sql"`
	ParamMap  JSONAnyMap   `json:"paramMap"`
	ParamList JSONAnySlice `json:"paramList"`
}

func (s *HTTPServer) handleWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var (
		isAuthenticated bool
		authToken       string
		authClaims      *TokenClaims
	)

	// Use a mutex for concurrent writes
	var writeMu sync.Mutex

	// Helper to write JSON response
	writeResponse := func(v any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		if s.WebSocketWriteTimeout > 0 {
			conn.SetWriteDeadline(time.Now().Add(s.WebSocketWriteTimeout))
		}
		return conn.WriteJSON(v)
	}

	// Helper to send Ping
	sendPing := func() error {
		writeMu.Lock()
		defer writeMu.Unlock()
		if s.WebSocketWriteTimeout > 0 {
			conn.SetWriteDeadline(time.Now().Add(s.WebSocketWriteTimeout))
		}
		return conn.WriteMessage(websocket.PingMessage, nil)
	}

	// Heartbeat loop
	stopHeartbeat := make(chan struct{})
	defer close(stopHeartbeat)
	go func() {
		ticker := time.NewTicker(s.WebSocketHeartbeat)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := sendPing(); err != nil {
					return
				}
			case <-stopHeartbeat:
				return
			}
		}
	}()

	// Read loop
	for {
		if s.WebSocketReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(s.WebSocketReadTimeout))
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var req wsRequest
		if err := json.Unmarshal(message, &req); err != nil {
			writeResponse(SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "invalid json"}})
			continue
		}

		route := normalizePath(req.Route, "")

		// Handle Login
		if route == "/login" {
			claims, err := s.tokenManager.Parse(req.Token)
			if err != nil {
				writeResponse(LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "invalid token"}})
				continue
			}
			isAuthenticated = true
			authToken = req.Token
			authClaims = claims
			writeResponse(LoginResponse{
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
			})
			continue
		}

		// Require Authentication
		if !isAuthenticated {
			writeResponse(LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "login required"}})
			continue
		}

		// Check token expiry
		if _, err := s.tokenManager.Parse(authToken); err != nil {
			writeResponse(LoginResponse{Success: false, Error: &ErrorBody{Code: "unauthorized", Message: "token expired"}})
			break // Disconnect
		}

		// Throttling
		if s.LimitEnabled {
			limitKey := "token:" + authToken
			allowed, _ := s.allowLimitKey(limitKey, time.Now())
			if !allowed {
				writeResponse(LoginResponse{Success: false, Error: &ErrorBody{Code: "too_many_requests", Message: "rate limit exceeded"}})
				continue
			}
		}

		// Handle message
		go func(req wsRequest, claims *TokenClaims) {
			// Routing
			// User said "Socket part only implement execution of DDL and DML".
			// Assume route "/sql/execute" matches SQLExecutePath
			// Normalizing route to match config or fixed string?
			// User said "route unified bind to /ws", and "request body adds a route field".
			// I'll check if req.Route matches "/sql/execute" (or configured path)

			targetRoute := normalizePath(req.Route, "")
			expectedRoute := normalizePath(s.SQLExecutePath, DefaultSQLExecutePath)

			if targetRoute != expectedRoute {
				writeResponse(SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "not_found", Message: "route not found or not supported in websocket"}})
				return
			}

			// Execute SQL
			// Use context for timeout
			execCtx, cancel := context.WithTimeout(context.Background(), s.WebSocketExecutionTimeout)
			defer cancel()

			// Run execution in a channel to support timeout
			resultChan := make(chan SQLExecuteResponse, 1)

			go func() {
				// Reconstruct SQLExecuteRequest
				sqlReq := SQLExecuteRequest{
					SQL:       req.SQL,
					ParamMap:  req.ParamMap,
					ParamList: req.ParamList,
				}

				// Logic copied from handleSQLExecute but adapted
				if strings.TrimSpace(sqlReq.SQL) == "" {
					resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: "sql is required"}}
					return
				}

				boundSQL, err := bindSQLParams(sqlReq.SQL, sqlReq.ParamMap, []any(sqlReq.ParamList))
				if err != nil {
					resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: err.Error()}}
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

				// Whitelist check
				if len(s.SQLAllowedOps) > 0 {
					stmt, parseErr := engine.Parse(boundSQL)
					if parseErr != nil {
						resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "bad_request", Message: parseErr.Error()}}
						return
					}
					opKey := strings.ToLower(strings.TrimPrefix(string(stmt.Type()), ""))
					if _, allowed := s.SQLAllowedOps[opKey]; !allowed {
						resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: "operation not allowed: " + opKey}}
						return
					}
				}

				result, err := engine.Execute(boundSQL)
				if err != nil {
					if errors.Is(err, api.ErrSystemTableAccessDenied) {
						resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "forbidden", Message: err.Error()}}
						return
					}
					status, code, message := mapDriverError(err, "sql execute failed")
					// Map status back to error code if needed, but here we just return body
					if status == http.StatusInternalServerError && code == "internal_error" {
						code = "bad_request"
						message = err.Error()
					}
					resultChan <- SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: code, Message: message}}
					return
				}

				resultChan <- SQLExecuteResponse{Success: true, Result: &result}
			}()

			select {
			case res := <-resultChan:
				writeResponse(res)
			case <-execCtx.Done():
				writeResponse(SQLExecuteResponse{Success: false, Error: &ErrorBody{Code: "timeout", Message: "execution timeout"}})
			}
		}(req, authClaims)
	}
}
