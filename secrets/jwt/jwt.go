// Package jwt 提供 JWT 的生成与校验，通过 secrets.Asymmetric 接入各类签名算法。
// 本包位于 secrets/jwt，作为 secrets 下的工具集，与 asymmetric 下的具体算法实现解耦。
package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aid297/aid/v2/secrets"
)

// Header JWT 头部
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Claims JWT 标准声明
type Claims struct {
	Iss   string         `json:"iss,omitempty"` // 签发者
	Sub   string         `json:"sub,omitempty"` // 主题
	Aud   string         `json:"aud,omitempty"` // 受众
	Exp   int64          `json:"exp,omitempty"` // 过期时间
	Nbf   int64          `json:"nbf,omitempty"` // 生效时间
	Iat   int64          `json:"iat,omitempty"` // 签发时间
	Jti   string         `json:"jti,omitempty"` // JWT ID
	Extra map[string]any `json:"-"`             // 自定义声明
}

// JWT JWT 实现，支持 RS256/ES256/EdDSA/SM2 等算法
type JWT struct {
	alg   string // 算法标识：RS256/ES256/EdDSA/SM2
	asymm secrets.Asymmetric
}

// Alg JWT 算法枚举
type Alg string

const (
	AlgRS256 Alg = "RS256" // RSA + SHA-256 + PKCS1v15
	AlgRS384 Alg = "RS384" // RSA + SHA-384 + PKCS1v15
	AlgRS512 Alg = "RS512" // RSA + SHA-512 + PKCS1v15
	AlgES256 Alg = "ES256" // ECDSA P-256 + SHA-256
	AlgES384 Alg = "ES384" // ECDSA P-384 + SHA-384
	AlgES512 Alg = "ES512" // ECDSA P-521 + SHA-512
	AlgEdDSA Alg = "EdDSA" // Ed25519
	AlgSM2   Alg = "SM2"   // 国密 SM2
)

// New 创建 JWT 实例，使用 secrets.Asymmetric 接口
// 算法标识 alg 可选，不设置时 header 中不包含 alg 字段
func New(asymm secrets.Asymmetric) *JWT {
	return &JWT{asymm: asymm}
}

// NewWithAlg 创建 JWT 实例并指定算法标识
// alg 支持：RS256/RS384/RS512/ES256/ES384/ES512/EdDSA/SM2
func NewWithAlg(alg Alg, asymm secrets.Asymmetric) *JWT {
	return &JWT{alg: string(alg), asymm: asymm}
}

// GetAlg 获取当前算法标识
func (my *JWT) GetAlg() string { return my.alg }

// Generate 生成 JWT token
func (my *JWT) Generate(claims *Claims) (string, error) {
	var (
		err           error
		header        = Header{Typ: "JWT"}
		headerBytes   []byte
		headerBase64  string
		payloadBytes  []byte
		payloadBase64 string
		signInput     string
		signature     string
		token         string
	)
	if claims == nil {
		return "", errors.New("claims 不能为空")
	}
	if my.asymm == nil {
		return "", errors.New("Asymmetric 不能为空")
	}

	// 设置 iat 默认值
	if claims.Iat == 0 {
		claims.Iat = time.Now().Unix()
	}

	// 编码 header（包含算法标识）
	if my.alg != "" {
		header.Alg = my.alg
	}
	if headerBytes, err = json.Marshal(header); err != nil {
		return "", err
	}
	headerBase64 = base64.RawURLEncoding.EncodeToString(headerBytes)

	// 编码 payload
	if payloadBytes, err = json.Marshal(claims); err != nil {
		return "", err
	}
	payloadBase64 = base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 签名
	signInput = fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	if signature, err = my.asymm.Sign([]byte(signInput)); err != nil {
		return "", err
	}

	// 自验证
	token = fmt.Sprintf("%s.%s", signInput, signature)
	if _, err = my.Verify(token); err != nil {
		return "", err
	}

	return token, nil
}

// Verify 验证 JWT token
func (my *JWT) Verify(tokenString string) (*Claims, error) {
	if my.asymm == nil {
		return nil, errors.New("Asymmetric 不能为空")
	}

	var (
		err          error
		parts        []string
		signInput    string
		valid        bool
		payloadBytes []byte
		claims       Claims
		now          int64
	)

	if parts = strings.Split(tokenString, "."); len(parts) != 3 {
		return nil, errors.New("令牌格式错误")
	}

	headerBase64, payloadBase64, signatureBase64 := parts[0], parts[1], parts[2]

	// 验证签名
	signInput = headerBase64 + "." + payloadBase64
	if valid, err = my.asymm.Verify([]byte(signInput), signatureBase64); err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("签名不合法")
	}

	// 解析 payload
	if payloadBytes, err = base64.RawURLEncoding.DecodeString(payloadBase64); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	// 验证时间声明
	now = time.Now().Unix()

	// 验证 exp (过期时间)
	if claims.Exp != 0 && now > claims.Exp {
		return nil, errors.New("令牌已过期")
	}

	// 验证 nbf (生效时间)
	if claims.Nbf != 0 && now < claims.Nbf {
		return nil, errors.New("令牌未生效")
	}

	return &claims, nil
}
