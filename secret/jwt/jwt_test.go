package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/aid297/aid/v2/secret"
	"github.com/aid297/aid/v2/secret/asymmetric/ecdsa"
	"github.com/aid297/aid/v2/secret/asymmetric/ed25519"
	"github.com/aid297/aid/v2/secret/asymmetric/rsa"
	"github.com/aid297/aid/v2/secret/asymmetric/sm2"
)

func TestGenerateAndVerify(t *testing.T) {
	// 生成 RSA 密钥对
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	// 创建 JWT 实例（使用 secret.Asymmetric）
	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := New(asymm)

	// 构建声明
	claims := &Claims{
		Iss: "test-issuer",
		Sub: "test-subject",
		Aud: "test-audience",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour).Unix(),
		Nbf: time.Now().Unix() - 60,
		Jti: "unique-token-id",
		Extra: map[string]any{
			"name":  "张三",
			"role":  "admin",
			"score": 95.5,
		},
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	if token == "" {
		t.Fatal("生成的 token 为空")
	}

	// 验证 token
	verifiedClaims, err := jwtInstance.Verify(token)
	if err != nil {
		t.Fatalf("验证 JWT 失败: %v", err)
	}

	// 验证声明
	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配: 期望 %s, 实际 %s", claims.Iss, verifiedClaims.Iss)
	}
	if verifiedClaims.Sub != claims.Sub {
		t.Errorf("Sub 不匹配: 期望 %s, 实际 %s", claims.Sub, verifiedClaims.Sub)
	}
	if verifiedClaims.Aud != claims.Aud {
		t.Errorf("Aud 不匹配: 期望 %s, 实际 %s", claims.Aud, verifiedClaims.Aud)
	}
	if verifiedClaims.Jti != claims.Jti {
		t.Errorf("Jti 不匹配: 期望 %s, 实际 %s", claims.Jti, verifiedClaims.Jti)
	}
}

func TestVerifyExpiredToken(t *testing.T) {
	// 生成 RSA 密钥对
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := New(asymm)

	// 构建已过期的声明
	claims := &Claims{
		Iss: "test-issuer",
		Iat: time.Now().Add(-2 * time.Hour).Unix(),
		Exp: time.Now().Add(-1 * time.Hour).Unix(), // 已过期 1 小时
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	// 验证应失败
	_, err = jwtInstance.Verify(token)
	if err == nil {
		t.Fatal("应该返回过期错误")
	}
	if err.Error() != "token expired" {
		t.Errorf("期望过期错误，实际: %v", err)
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	// 生成两组密钥对
	rsaSem1, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}
	rsaSem2, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	// 用第一组密钥生成 token
	var asymm1 secret.Asymmetric = rsa.New(rsaSem1)
	claims := &Claims{
		Iss: "test-issuer",
		Exp: time.Now().Add(time.Hour).Unix(),
	}
	token, err := New(asymm1).Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	// 用第二组密钥验证（应该失败）
	var asymm2 secret.Asymmetric = rsa.New(rsaSem2)
	_, err = New(asymm2).Verify(token)
	if err == nil {
		t.Fatal("应该返回签名错误")
	}
}

func TestVerifyNotYetValid(t *testing.T) {
	// 生成 RSA 密钥对
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := New(asymm)

	// 构建尚未生效的声明
	claims := &Claims{
		Iss: "test-issuer",
		Iat: time.Now().Unix(),
		Nbf: time.Now().Add(1 * time.Hour).Unix(), // 1 小时后才生效
		Exp: time.Now().Add(2 * time.Hour).Unix(),
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	// 验证应失败（尚未生效）
	_, err = jwtInstance.Verify(token)
	if err == nil {
		t.Fatal("应该返回尚未生效错误")
	}
	if err.Error() != "token not yet valid" {
		t.Errorf("期望尚未生效错误，实际: %v", err)
	}
}

func TestVerifyInvalidFormat(t *testing.T) {
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := New(asymm)

	// 测试各种无效格式
	invalidTokens := []string{
		"",
		"no-dot",
		"only.two.parts",
		"too.many.parts.here",
	}

	for _, token := range invalidTokens {
		_, err := jwtInstance.Verify(token)
		if err == nil {
			t.Errorf("Token %q 应该返回格式错误", token)
		}
	}
}

func TestWithExistingKeyPair(t *testing.T) {
	// 生成 RSA 密钥对
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	// 获取公私钥
	pubKeyBytes, _ := rsaSem.GetPubKeyBytes()
	priKeyBytes, _ := rsaSem.GetPriKeyBytes()

	// 使用公钥创建 JWT（仅验证）
	rsaSemPub, _ := rsa.NewSem(rsa.PubKeyBytes(pubKeyBytes))
	// 使用私钥创建 JWT（仅签名）
	rsaSemPri, _ := rsa.NewSem(rsa.PriKeyBytes(priKeyBytes))

	// 签名
	var signerAsymm secret.Asymmetric = rsa.New(rsaSemPri)
	signer := New(signerAsymm)
	claims := &Claims{
		Iss: "test",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := signer.Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	// 验证
	var verifierAsymm secret.Asymmetric = rsa.New(rsaSemPub)
	verifiedClaims, err := New(verifierAsymm).Verify(token)
	if err != nil {
		t.Fatalf("验证 JWT 失败: %v", err)
	}

	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配")
	}
}

// TestES256GenerateAndVerify 测试 ES256（ECDSA P-256）签名和验证
func TestES256GenerateAndVerify(t *testing.T) {
	// 生成 ECDSA 密钥对
	ecdsaSem, err := ecdsa.NewSem()
	if err != nil {
		t.Fatalf("生成 ECDSA 密钥对失败: %v", err)
	}

	// 创建 JWT 实例（使用 ES256 算法）
	var asymm secret.Asymmetric = ecdsa.New(ecdsaSem)
	jwtInstance := NewWithAlg(AlgES256, asymm)

	// 构建声明
	claims := &Claims{
		Iss: "es256-issuer",
		Sub: "es256-subject",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 ES256 JWT 失败: %v", err)
	}

	// 验证 token
	verifiedClaims, err := jwtInstance.Verify(token)
	if err != nil {
		t.Fatalf("验证 ES256 JWT 失败: %v", err)
	}

	// 验证声明
	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配: 期望 %s, 实际 %s", claims.Iss, verifiedClaims.Iss)
	}

	// 验证头部包含 ES256 算法标识
	parts := strings.Split(token, ".")
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header Header
	json.Unmarshal(headerBytes, &header)
	if header.Alg != "ES256" {
		t.Errorf("Header alg 不匹配: 期望 ES256, 实际 %s", header.Alg)
	}
}

// TestEdDSAGenerateAndVerify 测试 EdDSA（Ed25519）签名和验证
func TestEdDSAGenerateAndVerify(t *testing.T) {
	// 生成 Ed25519 密钥对
	edSem, err := ed25519.NewSem()
	if err != nil {
		t.Fatalf("生成 Ed25519 密钥对失败: %v", err)
	}

	// 创建 JWT 实例（使用 EdDSA 算法）
	var asymm secret.Asymmetric = ed25519.New(edSem)
	jwtInstance := NewWithAlg(AlgEdDSA, asymm)

	// 构建声明
	claims := &Claims{
		Iss: "eddsa-issuer",
		Sub: "eddsa-subject",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 EdDSA JWT 失败: %v", err)
	}

	// 验证 token
	verifiedClaims, err := jwtInstance.Verify(token)
	if err != nil {
		t.Fatalf("验证 EdDSA JWT 失败: %v", err)
	}

	// 验证声明
	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配: 期望 %s, 实际 %s", claims.Iss, verifiedClaims.Iss)
	}

	// 验证头部包含 EdDSA 算法标识
	parts := strings.Split(token, ".")
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header Header
	json.Unmarshal(headerBytes, &header)
	if header.Alg != "EdDSA" {
		t.Errorf("Header alg 不匹配: 期望 EdDSA, 实际 %s", header.Alg)
	}
}

// TestSM2GenerateAndVerify 测试 SM2 国密签名和验证
func TestSM2GenerateAndVerify(t *testing.T) {
	// 生成 SM2 密钥对
	sm2Sem, err := sm2.NewSem()
	if err != nil {
		t.Fatalf("生成 SM2 密钥对失败: %v", err)
	}

	// 创建 JWT 实例（使用 SM2 算法）
	var asymm secret.Asymmetric = sm2.New(sm2Sem)
	jwtInstance := NewWithAlg(AlgSM2, asymm)

	// 构建声明
	claims := &Claims{
		Iss: "sm2-issuer",
		Sub: "sm2-subject",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	// 生成 token
	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 SM2 JWT 失败: %v", err)
	}

	// 验证 token
	verifiedClaims, err := jwtInstance.Verify(token)
	if err != nil {
		t.Fatalf("验证 SM2 JWT 失败: %v", err)
	}

	// 验证声明
	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配: 期望 %s, 实际 %s", claims.Iss, verifiedClaims.Iss)
	}

	// 验证头部包含 SM2 算法标识
	parts := strings.Split(token, ".")
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header Header
	json.Unmarshal(headerBytes, &header)
	if header.Alg != "SM2" {
		t.Errorf("Header alg 不匹配: 期望 SM2, 实际 %s", header.Alg)
	}
}

// TestRS256WithAlg 测试 RS256 带算法标识的 JWT
func TestRS256WithAlg(t *testing.T) {
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := NewWithAlg(AlgRS256, asymm)

	claims := &Claims{
		Iss: "rs256-issuer",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 RS256 JWT 失败: %v", err)
	}

	// 验证头部包含 RS256 算法标识
	parts := strings.Split(token, ".")
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header Header
	json.Unmarshal(headerBytes, &header)
	if header.Alg != "RS256" {
		t.Errorf("Header alg 不匹配: 期望 RS256, 实际 %s", header.Alg)
	}
	if header.Typ != "JWT" {
		t.Errorf("Header typ 不匹配: 期望 JWT, 实际 %s", header.Typ)
	}
}

// TestJWTWithoutAlg 测试不带算法标识的 JWT
func TestJWTWithoutAlg(t *testing.T) {
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)
	jwtInstance := New(asymm) // 不设置算法

	claims := &Claims{
		Iss: "no-alg-issuer",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := jwtInstance.Generate(claims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	// 验证头部不包含 alg 字段
	parts := strings.Split(token, ".")
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var header Header
	json.Unmarshal(headerBytes, &header)
	if header.Alg != "" {
		t.Errorf("Header alg 应该为空，实际: %s", header.Alg)
	}
	if header.Typ != "JWT" {
		t.Errorf("Header typ 不匹配: 期望 JWT, 实际 %s", header.Typ)
	}
}

// TestGetAlg 测试获取算法标识
func TestGetAlg(t *testing.T) {
	rsaSem, err := rsa.NewSem()
	if err != nil {
		t.Fatalf("生成密钥对失败: %v", err)
	}

	var asymm secret.Asymmetric = rsa.New(rsaSem)

	// 不带算法
	jwt1 := New(asymm)
	if jwt1.GetAlg() != "" {
		t.Errorf("GetAlg 应该为空")
	}

	// 带算法
	jwt2 := NewWithAlg(AlgES256, asymm)
	if jwt2.GetAlg() != "ES256" {
		t.Errorf("GetAlg 应该返回 ES256")
	}
}

// TestES256WithExistingKeyPair 测试 ES256 使用已有密钥对
func TestES256WithExistingKeyPair(t *testing.T) {
	// 生成 ECDSA 密钥对
	ecdsaSem, err := ecdsa.NewSem()
	if err != nil {
		t.Fatalf("生成 ECDSA 密钥对失败: %v", err)
	}

	// 获取公私钥
	pubKeyBytes, _ := ecdsaSem.GetPubKeyBytes()
	priKeyBytes, _ := ecdsaSem.GetPriKeyBytes()

	// 使用公钥创建 JWT（仅验证）
	ecdsaSemPub, _ := ecdsa.NewSem(ecdsa.PubKeyBytes(pubKeyBytes))
	// 使用私钥创建 JWT（仅签名）
	ecdsaSemPri, _ := ecdsa.NewSem(ecdsa.PriKeyBytes(priKeyBytes))

	// 签名
	var signerAsymm secret.Asymmetric = ecdsa.New(ecdsaSemPri)
	signer := NewWithAlg(AlgES256, signerAsymm)
	claims := &Claims{
		Iss: "test",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := signer.Generate(claims)
	if err != nil {
		t.Fatalf("生成 ES256 JWT 失败: %v", err)
	}

	// 验证
	var verifierAsymm secret.Asymmetric = ecdsa.New(ecdsaSemPub)
	verifiedClaims, err := NewWithAlg(AlgES256, verifierAsymm).Verify(token)
	if err != nil {
		t.Fatalf("验证 ES256 JWT 失败: %v", err)
	}

	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配")
	}
}

// TestEdDSAWithExistingKeyPair 测试 EdDSA 使用已有密钥对
func TestEdDSAWithExistingKeyPair(t *testing.T) {
	// 生成 Ed25519 密钥对
	edSem, err := ed25519.NewSem()
	if err != nil {
		t.Fatalf("生成 Ed25519 密钥对失败: %v", err)
	}

	// 获取公私钥
	pubKeyBytes, _ := edSem.GetPubKeyBytes()
	priKeyBytes, _ := edSem.GetPriKeyBytes()

	// 使用公钥创建 JWT（仅验证）
	edSemPub, _ := ed25519.NewSem(ed25519.PubKeyBytes(pubKeyBytes))
	// 使用私钥创建 JWT（仅签名）
	edSemPri, _ := ed25519.NewSem(ed25519.PriKeyBytes(priKeyBytes))

	// 签名
	var signerAsymm secret.Asymmetric = ed25519.New(edSemPri)
	signer := NewWithAlg(AlgEdDSA, signerAsymm)
	claims := &Claims{
		Iss: "test",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := signer.Generate(claims)
	if err != nil {
		t.Fatalf("生成 EdDSA JWT 失败: %v", err)
	}

	// 验证
	var verifierAsymm secret.Asymmetric = ed25519.New(edSemPub)
	verifiedClaims, err := NewWithAlg(AlgEdDSA, verifierAsymm).Verify(token)
	if err != nil {
		t.Fatalf("验证 EdDSA JWT 失败: %v", err)
	}

	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配")
	}
}

// TestSM2WithExistingKeyPair 测试 SM2 使用已有密钥对
func TestSM2WithExistingKeyPair(t *testing.T) {
	// 生成 SM2 密钥对
	sm2Sem, err := sm2.NewSem()
	if err != nil {
		t.Fatalf("生成 SM2 密钥对失败: %v", err)
	}

	// 获取公私钥
	pubKeyBytes, _ := sm2Sem.GetPubKeyBytes()
	priKeyBytes, _ := sm2Sem.GetPriKeyBytes()

	// 使用公钥创建 JWT（仅验证）
	sm2SemPub, _ := sm2.NewSem(sm2.PubKeyBytes(pubKeyBytes))
	// 使用私钥创建 JWT（仅签名）
	sm2SemPri, _ := sm2.NewSem(sm2.PriKeyBytes(priKeyBytes))

	// 签名
	var signerAsymm secret.Asymmetric = sm2.New(sm2SemPri)
	signer := NewWithAlg(AlgSM2, signerAsymm)
	claims := &Claims{
		Iss: "test",
		Exp: time.Now().Add(time.Hour).Unix(),
	}

	token, err := signer.Generate(claims)
	if err != nil {
		t.Fatalf("生成 SM2 JWT 失败: %v", err)
	}

	// 验证
	var verifierAsymm secret.Asymmetric = sm2.New(sm2SemPub)
	verifiedClaims, err := NewWithAlg(AlgSM2, verifierAsymm).Verify(token)
	if err != nil {
		t.Fatalf("验证 SM2 JWT 失败: %v", err)
	}

	if verifiedClaims.Iss != claims.Iss {
		t.Errorf("Iss 不匹配")
	}
}
