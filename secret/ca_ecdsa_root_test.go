package secret_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	myECDSA "github.com/aid297/aid/v2/secret/asymmetric/ecdsa"
)

// caOutDir 返回 secret/test.data/test.ecdsa.genrerate.ca.root.certificate/（与本测试文件同级）。
func caOutDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller 失败")
	}
	return filepath.Join(filepath.Dir(file), "test.data", "test.ecdsa.genrerate.ca.root.certificate")
}

// agentChainOutDir 输出：Agent 与 CA 链式演示（CSR、叶子证书等）。
func agentChainOutDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller 失败")
	}
	return filepath.Join(filepath.Dir(file), "test.data", "test.ecdsa.server.ca.agent.chain")
}

func newSN() (*big.Int, error) { return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)) }

// newServerECDSACARoot 生成服务器侧自签 ECDSA P-256 CA 根证书与私钥。
func newServerECDSACARoot(t *testing.T) (serverRootCert *x509.Certificate, serverRootPrivKey *ecdsa.PrivateKey) {
	t.Helper()
	serverCASem, err := myECDSA.NewSem()
	if err != nil {
		t.Fatalf("生成服务器 CA 种子失败：%v", err)
	}
	serverRootPrivKey, ok := serverCASem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || serverRootPrivKey == nil {
		t.Fatal("服务器 CA 私钥应为 *ecdsa.PrivateKey")
	}
	serverRootPubKey, ok := serverCASem.GetPubKey().(*ecdsa.PublicKey)
	if !ok || serverRootPubKey == nil {
		t.Fatal("服务器 CA 公钥应为 *ecdsa.PublicKey")
	}
	sn, err := newSN()
	if err != nil {
		t.Fatalf("生成序列号失败：%v", err)
	}
	now := time.Now()
	serverRootTemplate := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"aid-secret-test"},
			CommonName:   "ECDSA P-256 Server Root CA",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}
	serverRootDER, err := x509.CreateCertificate(rand.Reader, serverRootTemplate, serverRootTemplate, serverRootPubKey, serverRootPrivKey)
	if err != nil {
		t.Fatalf("CreateCertificate(服务器 CA) 失败：%v", err)
	}
	serverRootCert, err = x509.ParseCertificate(serverRootDER)
	if err != nil {
		t.Fatalf("ParseCertificate(服务器 CA) 失败：%v", err)
	}
	if err = serverRootCert.CheckSignatureFrom(serverRootCert); err != nil {
		t.Fatalf("服务器 CA 自签名校验失败：%v", err)
	}
	return serverRootCert, serverRootPrivKey
}

// TestECDSA_GenerateCARootCertificate 使用 asymmetric/ecdsa 种子生成 P-256 密钥对，
// 自签 CA 根证书，并将证书、私钥、公钥 PEM 写入 secret/test.data/test.ecdsa.genrerate.ca.root.certificate/。
func TestECDSA_GenerateCARootCertificate(t *testing.T) {
	sem, err := myECDSA.NewSem()
	if err != nil {
		t.Fatalf("生成 ECDSA 种子失败：%v", err)
	}

	priv, ok := sem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || priv == nil {
		t.Fatal("私钥应为 *ecdsa.PrivateKey")
	}
	pub, ok := sem.GetPubKey().(*ecdsa.PublicKey)
	if !ok || pub == nil {
		t.Fatal("公钥应为 *ecdsa.PublicKey")
	}

	serialNumber, err := newSN()
	if err != nil {
		t.Fatalf("生成序列号失败：%v", err)
	}

	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"aid-secret-test"},
			CommonName:   "ECDSA P-256 Test Root CA",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, priv)
	if err != nil {
		t.Fatalf("CreateCertificate 失败：%v", err)
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate 失败：%v", err)
	}

	if cert.PublicKeyAlgorithm != x509.ECDSA {
		t.Fatalf("期望公钥算法 ECDSA，实际 %v", cert.PublicKeyAlgorithm)
	}
	if !cert.IsCA {
		t.Fatal("根证书应设置 IsCA")
	}
	if err = cert.CheckSignatureFrom(cert); err != nil {
		t.Fatalf("自签名校验失败：%v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if len(pemBytes) == 0 {
		t.Fatal("PEM 编码结果为空")
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatal("PEM 解码失败或类型不是 CERTIFICATE")
	}
	if !bytes.Equal(block.Bytes, der) {
		t.Fatal("PEM 内 DER 与原始 DER 不一致")
	}

	outDir := caOutDir(t)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("创建输出目录失败：%v", err)
	}

	certPath := filepath.Join(outDir, "ca.crt")
	keyPath := filepath.Join(outDir, "ca.key")
	pubPath := filepath.Join(outDir, "ca.pub")

	if err := os.WriteFile(certPath, pemBytes, 0o644); err != nil {
		t.Fatalf("写入 CA 证书失败：%v", err)
	}

	keyPEM, err := sem.GetPriKeyPEM()
	if err != nil {
		t.Fatalf("GetPriKeyPEM 失败：%v", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0o600); err != nil {
		t.Fatalf("写入 CA 私钥失败：%v", err)
	}

	pubPEM, err := sem.GetPubKeyPEM()
	if err != nil {
		t.Fatalf("GetPubKeyPEM 失败：%v", err)
	}
	if err := os.WriteFile(pubPath, pubPEM, 0o644); err != nil {
		t.Fatalf("写入 CA 公钥失败：%v", err)
	}

	loadedCertPEM, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatalf("回读证书失败：%v", err)
	}
	b, _ := pem.Decode(loadedCertPEM)
	if b == nil {
		t.Fatal("回读证书 PEM 解码失败")
	}
	loadedCert, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		t.Fatalf("回读证书解析失败：%v", err)
	}
	if err = loadedCert.CheckSignatureFrom(loadedCert); err != nil {
		t.Fatalf("磁盘证书自签名校验失败：%v", err)
	}

	t.Logf("CA 根证书主题：%s", cert.Subject.String())
	t.Logf("有效期至：%s", cert.NotAfter.Format(time.RFC3339))
	t.Logf("已写入本地：\n  证书 %s\n  私钥 %s\n  公钥 %s", certPath, keyPath, pubPath)
}

// TestECDSA_ServerCASignsAgentCSR_AgentVerifiesCAAndLeaf 模拟：
// 1）服务器持有 ECDSA CA 根证书与私钥；
// 2）Agent 生成自有密钥，用私钥对 CSR 签名（CSR 内携带公钥）；
// 3）服务器校验 CSR 后，用 CA 私钥为 Agent 签发叶子证书；
// 4）Agent 仅拿到 PEM 形式的 CA 根证书：先校验 CA 自签名有效，再校验叶子证书由该 CA 签发。
func TestECDSA_ServerCASignsAgentCSR_AgentVerifiesCAAndLeaf(t *testing.T) {
	serverRootCert, serverRootPrivKey := newServerECDSACARoot(t)

	agentKeySem, err := myECDSA.NewSem()
	if err != nil {
		t.Fatalf("生成 Agent 密钥种子失败：%v", err)
	}
	agentPrivateKey, ok := agentKeySem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || agentPrivateKey == nil {
		t.Fatal("Agent 私钥应为 *ecdsa.PrivateKey")
	}
	agentPublicKey, ok := agentKeySem.GetPubKey().(*ecdsa.PublicKey)
	if !ok || agentPublicKey == nil {
		t.Fatal("Agent 公钥应为 *ecdsa.PublicKey")
	}

	agentCSRTemplate := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"aid-agent"},
			CommonName:   "ecdsa-test-agent",
		},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}
	agentCSRDER, err := x509.CreateCertificateRequest(rand.Reader, agentCSRTemplate, agentPrivateKey)
	if err != nil {
		t.Fatalf("Agent 创建 CSR 失败：%v", err)
	}

	agentCSR, err := x509.ParseCertificateRequest(agentCSRDER)
	if err != nil {
		t.Fatalf("解析 Agent CSR 失败：%v", err)
	}
	if err = agentCSR.CheckSignature(); err != nil {
		t.Fatalf("Agent CSR 签名校验失败（私钥与 CSR 中公钥不一致）：%v", err)
	}
	agentCSRPublicKey, ok := agentCSR.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Agent CSR 公钥应为 ECDSA")
	}
	if !agentCSRPublicKey.Equal(agentPublicKey) {
		t.Fatal("Agent CSR 中公钥与本地 Agent 公钥不一致")
	}

	serverLeafSerialNumber, err := newSN()
	if err != nil {
		t.Fatalf("服务器生成叶子证书序列号失败：%v", err)
	}
	serverLeafTemplate := &x509.Certificate{
		SerialNumber:          serverLeafSerialNumber,
		Subject:               agentCSR.Subject,
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	serverSignedAgentLeafDER, err := x509.CreateCertificate(rand.Reader, serverLeafTemplate, serverRootCert, agentCSR.PublicKey, serverRootPrivKey)
	if err != nil {
		t.Fatalf("服务器用 CA 签发 Agent 叶子证书失败：%v", err)
	}
	serverSignedAgentLeafCert, err := x509.ParseCertificate(serverSignedAgentLeafDER)
	if err != nil {
		t.Fatalf("解析服务器签发的 Agent 叶子证书失败：%v", err)
	}

	// Agent 侧：仅持有服务器下发的 CA 根证书 PEM，校验信任锚与叶子链
	serverRootCAPEMForAgent := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverRootCert.Raw})
	pemBlock, _ := pem.Decode(serverRootCAPEMForAgent)
	if pemBlock == nil {
		t.Fatal("服务器 CA PEM 解码失败")
	}
	agentTrustedServerRootCA, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		t.Fatalf("Agent 解析服务器 CA 证书失败：%v", err)
	}
	if err = agentTrustedServerRootCA.CheckSignatureFrom(agentTrustedServerRootCA); err != nil {
		t.Fatalf("Agent 校验服务器 CA 根证书（自签名）失败：%v", err)
	}
	if !agentTrustedServerRootCA.IsCA {
		t.Fatal("Agent 信任锚应为 CA")
	}

	if err = serverSignedAgentLeafCert.CheckSignatureFrom(agentTrustedServerRootCA); err != nil {
		t.Fatalf("Agent 校验叶子证书（须由服务器 CA 签发）失败：%v", err)
	}
	agentLeafPublicKey, ok := serverSignedAgentLeafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Agent 叶子证书公钥应为 ECDSA")
	}
	if !agentLeafPublicKey.Equal(agentPublicKey) {
		t.Fatal("叶子证书公钥与 Agent 本地公钥不一致")
	}

	artifactDir := agentChainOutDir(t)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatalf("创建输出目录失败：%v", err)
	}
	serverRootCAPemPath := filepath.Join(artifactDir, "server_ca.crt")
	agentCSRPath := filepath.Join(artifactDir, "agent.csr.pem")
	serverIssuedAgentLeafPemPath := filepath.Join(artifactDir, "agent.crt")
	agentPrivateKeyPemPath := filepath.Join(artifactDir, "agent.key")

	if err := os.WriteFile(serverRootCAPemPath, serverRootCAPEMForAgent, 0o644); err != nil {
		t.Fatalf("写入 server_ca.crt：%v", err)
	}
	agentCSRPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: agentCSRDER})
	if err := os.WriteFile(agentCSRPath, agentCSRPEM, 0o644); err != nil {
		t.Fatalf("写入 agent.csr.pem：%v", err)
	}
	serverSignedAgentLeafPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverSignedAgentLeafDER})
	if err := os.WriteFile(serverIssuedAgentLeafPemPath, serverSignedAgentLeafPEM, 0o644); err != nil {
		t.Fatalf("写入 agent.crt：%v", err)
	}
	agentPrivateKeyPEM, err := agentKeySem.GetPriKeyPEM()
	if err != nil {
		t.Fatalf("GetPriKeyPEM(Agent)：%v", err)
	}
	if err := os.WriteFile(agentPrivateKeyPemPath, agentPrivateKeyPEM, 0o600); err != nil {
		t.Fatalf("写入 agent.key：%v", err)
	}

	t.Logf("服务器 CA：%s", agentTrustedServerRootCA.Subject.String())
	t.Logf("服务器签发之 Agent 叶子：%s", serverSignedAgentLeafCert.Subject.String())
	t.Logf("已写入：%s（server_ca.crt, agent.csr.pem, agent.crt, agent.key）", artifactDir)
}
