package transport

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aid297/aid/simpleDB/driver"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type TokenManager struct {
	mu       sync.RWMutex
	database string
	secret   []byte
	ttl      time.Duration
	revoked  map[string]int64
	store    TokenStore
}

type TokenClaims struct {
	TokenID     string   `json:"jti"`
	Subject     string   `json:"sub"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName,omitempty"`
	Status      string   `json:"status,omitempty"`
	IsAdmin     bool     `json:"isAdmin"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IssuedAt    int64    `json:"iat"`
	ExpiresAt   int64    `json:"exp"`
}

type IssuedToken struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
}

func NewTokenManager(database, secret string, ttl time.Duration) *TokenManager {
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}
	database = strings.TrimSpace(database)
	secret = strings.TrimSpace(secret)
	if secret == "" {
		secret = "simpledb.transport." + database + ".secret"
	}
	return &TokenManager{
		database: database,
		secret:   []byte(secret),
		ttl:      ttl,
		revoked:  make(map[string]int64),
	}
}

func (m *TokenManager) WithStore(store TokenStore) *TokenManager {
	m.store = store
	if m.store != nil && m.database != "" {
		now := time.Now().UTC().Unix()
		// 启动时清理过期 Token
		_ = m.store.ClearExpired(m.database, now)

		// 加载已撤销 Token 列表
		if loaded, err := m.store.LoadRevoked(m.database, now); err == nil {
			m.mu.Lock()
			for k, v := range loaded {
				m.revoked[k] = v
			}
			m.mu.Unlock()
		}
	}
	return m
}

func (m *TokenManager) Issue(user *driver.AuthenticatedUser) (*IssuedToken, error) {
	if user == nil {
		return nil, ErrInvalidToken
	}
	tokenID, err := newTokenID()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl).Unix()

	// 记录到可用 Token 列表
	if m.store != nil && m.database != "" {
		if err := m.store.SaveActive(m.database, tokenID, expiresAt); err != nil {
			return nil, err
		}
	}

	claims := TokenClaims{
		TokenID:     tokenID,
		Subject:     user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		IsAdmin:     user.IsAdmin,
		Roles:       append([]string(nil), user.Roles...),
		Permissions: append([]string(nil), user.Permissions...),
		IssuedAt:    now.Unix(),
		ExpiresAt:   expiresAt,
	}
	header := map[string]string{"alg": "HS256", "typ": "SDBT"}
	token, err := m.encode(header, claims)
	if err != nil {
		return nil, err
	}
	return &IssuedToken{AccessToken: token, TokenType: "Bearer", ExpiresAt: claims.ExpiresAt}, nil
}

func (m *TokenManager) Parse(token string) (*TokenClaims, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}
	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(m.sign(signingInput))) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	var claims TokenClaims
	if err = json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.TokenID == "" || claims.Subject == "" || claims.Username == "" {
		return nil, ErrInvalidToken
	}

	now := time.Now().UTC().Unix()

	// 1. 以数据库可用列表为准进行有效期检查（Source of Truth）
	if m.store != nil && m.database != "" {
		dbExpiresAt, err := m.store.GetActiveExpiresAt(m.database, claims.TokenID)
		if err != nil || dbExpiresAt <= now {
			return nil, ErrInvalidToken // 不在可用列表中或已过期
		}

		// 2. 自动续期检查：如果数据库中的过期时间距离现在少于 2 小时 (7200 秒)
		if dbExpiresAt-now <= 7200 {
			newExpiresAt := now + int64(m.ttl.Seconds())
			// 原地更新数据库中的过期时间，不生成新 Token 字符串
			_ = m.store.SaveActive(m.database, claims.TokenID, newExpiresAt)
		}
	} else {
		// 如果没有配置存储，则退化为检查 JWT 本身的过期时间
		if now >= claims.ExpiresAt {
			return nil, ErrExpiredToken
		}
	}

	// 3. 检查是否在已撤销列表中
	if m.isRevoked(claims.TokenID) {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (m *TokenManager) Revoke(token string) error {
	claims, err := m.Parse(token)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Unix()

	m.mu.Lock()
	m.revoked[claims.TokenID] = claims.ExpiresAt
	m.compactRevokedLocked(now)
	m.mu.Unlock()

	if m.store != nil && m.database != "" {
		_ = m.store.SaveRevoked(m.database, claims.TokenID, claims.ExpiresAt)
		_ = m.store.ClearExpired(m.database, now)
	}
	return nil
}

func (m *TokenManager) Refresh(token string) (*IssuedToken, *TokenClaims, error) {
	claims, err := m.Parse(token)
	if err != nil {
		return nil, nil, err
	}
	// 刷新意味着撤销旧的，颁发新的
	if err = m.Revoke(token); err != nil {
		return nil, nil, err
	}
	issued, err := m.Issue(&driver.AuthenticatedUser{
		ID:          claims.Subject,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		Status:      claims.Status,
		IsAdmin:     claims.IsAdmin,
		Roles:       append([]string(nil), claims.Roles...),
		Permissions: append([]string(nil), claims.Permissions...),
	})
	if err != nil {
		return nil, nil, err
	}
	newClaims, err := m.Parse(issued.AccessToken)
	if err != nil {
		return nil, nil, err
	}
	return issued, newClaims, nil
}

func (m *TokenManager) encode(header any, claims TokenClaims) (string, error) {
	headerRaw, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadRaw, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signingInput := base64.RawURLEncoding.EncodeToString(headerRaw) + "." + base64.RawURLEncoding.EncodeToString(payloadRaw)
	return fmt.Sprintf("%s.%s", signingInput, m.sign(signingInput)), nil
}

func (m *TokenManager) sign(signingInput string) string {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (m *TokenManager) isRevoked(tokenID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.revoked[tokenID]
	return ok
}

func (m *TokenManager) compactRevokedLocked(now int64) {
	for tokenID, expiresAt := range m.revoked {
		if now >= expiresAt {
			delete(m.revoked, tokenID)
		}
	}
}

func newTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
