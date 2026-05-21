package jwt_test

import (
	"testing"
	"time"

	"github.com/aid297/aid/v2/secret"
	"github.com/aid297/aid/v2/secret/asymmetric/rsa"
	"github.com/aid297/aid/v2/secret/jwt"
)

func TestGenerate(t *testing.T) {
	var (
		err                       error
		rsaSem                    secret.Semen
		rsaEncrypter              secret.Asymmetric
		jwtGenerate, jwtVerify    *jwt.JWT
		srcClaims, verifiedClaims *jwt.Claims
		taskToken                 string
	)
	// 生成 RSA 密钥对
	if rsaSem, err = rsa.NewSem(); err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	// 创建 JWT 实例（使用 secret.Asymmetric）
	rsaEncrypter = rsa.New(rsaSem)
	jwtGenerate = jwt.New(rsaEncrypter)

	// 构建声明
	srcClaims = &jwt.Claims{
		Iss: "签发机构",
		Sub: "归属客户端",                          // 如果是备份任务：ClientA
		Aud: "目标客户端",                          // 如果是备份任务：ClientB
		Iat: time.Now().Unix(),                // 签发时间
		Exp: time.Now().Add(time.Hour).Unix(), // 过期时间
		Nbf: time.Now().Unix() - 60,           // 生效时间
		Jti: "任务ID",
		Extra: map[string]any{
			"ip":          "目标客户端IP",
			"port":        "目标客户端端口",
			"pub_key_pem": "PEM 格式公钥内容", // 目标客户端的公钥
		}, // 额外的参数（非业务参数）
	}

	// 生成 taskToken
	taskToken, err = jwtGenerate.Generate(srcClaims)
	if err != nil {
		t.Fatalf("生成 JWT 失败: %v", err)
	}

	if taskToken == "" {
		t.Fatal("生成的 token 为空")
	}

	// 验证 token
	jwtVerify = jwt.New(rsaEncrypter)
	verifiedClaims, err = jwtVerify.Verify(taskToken)
	if err != nil {
		t.Fatalf("验证 JWT 失败: %v", err)
	}

	// 验证声明
	if verifiedClaims.Iss != srcClaims.Iss {
		t.Errorf("Iss 不匹配: 期望 %s, 实际 %s", srcClaims.Iss, verifiedClaims.Iss)
	}
	if verifiedClaims.Sub != srcClaims.Sub {
		t.Errorf("Sub 不匹配: 期望 %s, 实际 %s", srcClaims.Sub, verifiedClaims.Sub)
	}
	if verifiedClaims.Aud != srcClaims.Aud {
		t.Errorf("Aud 不匹配: 期望 %s, 实际 %s", srcClaims.Aud, verifiedClaims.Aud)
	}
	if verifiedClaims.Jti != srcClaims.Jti {
		t.Errorf("Jti 不匹配: 期望 %s, 实际 %s", srcClaims.Jti, verifiedClaims.Jti)
	}
}
