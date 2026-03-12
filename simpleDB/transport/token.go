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
	mu      sync.RWMutex
	secret  []byte
	ttl     time.Duration
	revoked map[string]int64
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
	secret = strings.TrimSpace(secret)
	if secret == "" {
		secret = "simpledb.transport." + strings.TrimSpace(database) + ".secret"
	}
	return &TokenManager{secret: []byte(secret), ttl: ttl, revoked: make(map[string]int64)}
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
		ExpiresAt:   now.Add(m.ttl).Unix(),
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
	if time.Now().UTC().Unix() >= claims.ExpiresAt {
		return nil, ErrExpiredToken
	}
	if m.isRevoked(claims.TokenID) {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}

func (m *TokenManager) Revoke(token string) error {
	claims, err := m.Parse(token)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.revoked[claims.TokenID] = claims.ExpiresAt
	m.compactRevokedLocked(time.Now().UTC().Unix())
	return nil
}

func (m *TokenManager) Refresh(token string) (*IssuedToken, *TokenClaims, error) {
	claims, err := m.Parse(token)
	if err != nil {
		return nil, nil, err
	}
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
