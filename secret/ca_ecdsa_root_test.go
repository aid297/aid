package secret_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aid297/aid/v2/filesystem"
	"github.com/aid297/aid/v2/secret"
	myECDSA "github.com/aid297/aid/v2/secret/asymmetric/ecdsa"
)

var testPath = filepath.Join("test.data", "ecdsa")

// 生成 ECDSA P-256 CA 根证书与私钥。
func newRootCrt(t *testing.T) (*x509.Certificate, secret.Semen) {
	t.Helper()

	var (
		err                              error
		caSem                            secret.Semen
		caPriv                           secret.SemenPriKey
		caPub                            secret.SemenPubKey
		caPrivPEM, caPubPEM, caPEM       []byte
		ok                               bool
		sn                               *big.Int
		now                              time.Time
		caTpl                            *x509.Certificate
		caDER                            []byte
		caPrivFile, caPubFile, caCrtFile filesystem.IFilesystem
		pemBlock                         *pem.Block
		rootCA                           *x509.Certificate
	)

	if caSem, err = myECDSA.NewSem(); err != nil {
		t.Fatalf("生成 CA 种子失败：%v", err)
	}

	if caPriv, ok = caSem.GetPriKey().(*ecdsa.PrivateKey); !ok {
		t.Fatalf("CA 私钥应为 *ecdsa.PrivateKey")
	}

	if caPrivPEM, err = caSem.GetPriKeyPEM(); err != nil {
		t.Fatalf("获取 CA 私钥 PEM 失败：%v", err)
	}

	caPrivFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.key"))
	if err = caPrivFile.Write(caPrivPEM, filesystem.Mode(0600), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 CA 私钥 PEM 失败：%v", err)
	}

	if caPub, ok = caSem.GetPubKey().(*ecdsa.PublicKey); !ok {
		t.Fatalf("CA 公钥应为 *ecdsa.PublicKey")
	}

	if caPubPEM, err = caSem.GetPubKeyPEM(); err != nil {
		t.Fatalf("获取 CA 公钥 PEM 失败：%v", err)
	}

	caPubFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.pub"))
	if err = caPubFile.Write(caPubPEM, filesystem.Mode(0644), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 CA 公钥 PEM 失败：%v", err)
	}

	if sn, err = newSN(); err != nil {
		t.Fatalf("生成序列号失败：%v", err)
	}

	// 获取 CA 证书
	now = time.Now()
	caTpl = &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"root-ca"},
			CommonName:   "ECDSA P-256 Root CA",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}
	if caDER, err = x509.CreateCertificate(rand.Reader, caTpl, caTpl, caPub, caPriv); err != nil {
		t.Fatalf("CA 证书模板 转 DER 格式失败：%v", err)
	}
	if rootCA, err = x509.ParseCertificate(caDER); err != nil {
		t.Fatalf("解析 CA 证书模板失败：%v", err)
	}
	if err = rootCA.CheckSignatureFrom(rootCA); err != nil {
		t.Fatalf("CA 证书自签名校验失败：%v", err)
	}

	// 将 DER 格式转为 PEM
	if caPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}); len(caPEM) == 0 {
		t.Fatal("CA 证书 PEM 编码结果为空")
	}
	if pemBlock, _ = pem.Decode(caPEM); pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		t.Fatal("PEM 解码失败或类型不是 CERTIFICATE")
	}
	if !bytes.Equal(pemBlock.Bytes, caDER) {
		t.Fatal("PEM 内 DER 与原始 DER 不一致")
	}
	caCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.crt"))
	if err = caCrtFile.Write(caPEM, filesystem.Mode(0644), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 CA 证书 PEM 失败：%v", err)
	}

	return rootCA, caSem
}

// 生成 ECDSA P-256 Server CA 与私钥。
func newServerCrt(t *testing.T) (*x509.Certificate, secret.Semen) {
	t.Helper()

	var (
		err                                                                 error
		serverSem, caSem                                                    secret.Semen
		caPriv                                                              secret.SemenPriKey
		serverPub                                                           secret.SemenPubKey
		serverPrivPEM, serverPubPEM, serverPEM                              []byte
		serverPrivFile, serverPubFile, caPrivFile, caCrtFile, serverCrtFile filesystem.IFilesystem
		ok                                                                  bool
		serverDER, caCrtPEM, caPrivPEM, caPrivDER                           []byte
		caCrtBlock, caPrivBlock                                             *pem.Block
		caCrt, serverTpl, serverCrt                                         *x509.Certificate
		sn                                                                  *big.Int
		now                                                                 time.Time
		pemBlock                                                            *pem.Block
	)

	// 获取 CA 根证书与私钥
	if caPrivFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.key")); !caPrivFile.GetExist() {
		t.Fatalf("[CA 私钥]文件不存在：%v", err)
	}
	if caPrivPEM, err = caPrivFile.Read(filesystem.Mode(0600)); err != nil {
		t.Fatalf("读取 [CA 私钥] PEM 失败：%v", err)
	}
	if caSem, err = myECDSA.NewSem(myECDSA.PriKeyBytes(caPrivPEM)); err != nil {
		t.Fatalf("生成 [CA 种子] 失败：%v", err)
	}
	if caPriv, ok = caSem.GetPriKey().(*ecdsa.PrivateKey); !ok {
		t.Fatalf("CA 私钥应为 *ecdsa.PrivateKey")
	}
	if caPrivBlock, _ = pem.Decode(caPrivPEM); caPrivBlock == nil || caPrivBlock.Type != "PRIVATE KEY" {
		t.Fatal("[CA 私钥] PEM 解码失败或类型不是 PRIVATE KEY")
	}
	if caPrivDER, err = caSem.GetPriKeyBytes(); err != nil {
		t.Fatalf("获取 [CA 私钥] DER 失败：%v", err)
	}
	if !bytes.Equal(caPrivBlock.Bytes, caPrivDER) {
		t.Fatal("[CA 私钥] PEM 内 DER 与原始 DER 不一致")
	}
	if caCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.crt")); !caCrtFile.GetExist() {
		t.Fatalf("[CA 证书]文件不存在：%v", err)
	}
	if caCrtPEM, err = caCrtFile.Read(filesystem.Mode(0644)); err != nil {
		t.Fatalf("读取 [CA 证书] PEM 失败：%v", err)
	}
	if caCrtBlock, _ = pem.Decode(caCrtPEM); caCrtBlock == nil || caCrtBlock.Type != "CERTIFICATE" {
		t.Fatal("[CA 证书] PEM 解码失败或类型不是 CERTIFICATE")
	}
	if caCrt, err = x509.ParseCertificate(caCrtBlock.Bytes); err != nil {
		t.Fatalf("解析 [CA 证书] 失败：%v", err)
	}

	// 生成服务器公私钥
	if serverSem, err = myECDSA.NewSem(); err != nil {
		t.Fatalf("生成 [服务器种子] 失败：%v", err)
	}
	if serverPrivPEM, err = serverSem.GetPriKeyPEM(); err != nil {
		t.Fatalf("获取 [服务器私钥] PEM 失败：%v", err)
	}
	serverPrivFile = filesystem.NewFile(filesystem.Rel(testPath, "server.key"))
	if err = serverPrivFile.Write(serverPrivPEM, filesystem.Mode(0600), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 [服务器私钥] PEM 失败：%v", err)
	}
	if serverPub, ok = serverSem.GetPubKey().(*ecdsa.PublicKey); !ok {
		t.Fatalf("[服务器公钥]应为 *ecdsa.PublicKey")
	}
	if serverPubPEM, err = serverSem.GetPubKeyPEM(); err != nil {
		t.Fatalf("获取 [服务器公钥] PEM 失败：%v", err)
	}
	serverPubFile = filesystem.NewFile(filesystem.Rel(testPath, "server.pub"))
	if err = serverPubFile.Write(serverPubPEM, filesystem.Mode(0644), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 [服务器公钥] PEM 失败：%v", err)
	}

	// 生成服务器证书模板
	if sn, err = newSN(); err != nil {
		t.Fatalf("生成序列号失败：%v", err)
	}
	now = time.Now()
	serverTpl = &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"server-ca"},
			CommonName:   "ECDSA P-256 Server CA",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	if serverDER, err = x509.CreateCertificate(rand.Reader, serverTpl, caCrt, serverPub, caPriv); err != nil {
		t.Fatalf("生成 [服务器证书] 失败：%v", err)
	}
	if serverCrt, err = x509.ParseCertificate(serverDER); err != nil {
		t.Fatalf("解析 [服务器证书] 失败：%v", err)
	}
	if err = serverCrt.CheckSignatureFrom(caCrt); err != nil {
		t.Fatalf("[服务器证书]校验失败：%v", err)
	}

	// 将 DER 格式转为 PEM
	if serverPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER}); len(caCrtPEM) == 0 {
		t.Fatal("[服务器证书] PEM 编码结果为空")
	}
	if pemBlock, _ = pem.Decode(serverPEM); pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		t.Fatal("[服务器证书] PEM 解码失败或类型不是 CERTIFICATE!!")
	}
	serverCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "server.crt"))
	if err = serverCrtFile.Write(serverPEM, filesystem.Mode(0644), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 [服务器证书] PEM 失败：%v", err)
	}

	return serverCrt, serverSem
}

// 生成 ECDSA P-256 Client CA 与私钥。
func newClientCrt(t *testing.T) (*x509.Certificate, secret.Semen) {
	t.Helper()

	var (
		err                         error
		ok                          bool
		clientSem, caSem            secret.Semen
		clientPub                   secret.SemenPubKey
		clientPrivPEM, clientPubPEM []byte
		caPrivFile, clientPrivFile, clientPubFile,
		caCrtFile, serverCrtFile, clientCrtFile filesystem.IFilesystem
		caCrtPEM, serverCrtPEM, caPrivPEM, clientDER []byte
		caCrtBlock, serverCrtBlock                   *pem.Block
		caCrt, serverCrt, clientTpl, clientCrt       *x509.Certificate
		caPriv, clientPriv                           secret.SemenPriKey
		sn                                           *big.Int
		now                                          time.Time
		clientCSR                                    *x509.CertificateRequest
		clientCSRDER                                 []byte
	)

	// 生成客户端：公私钥
	if clientSem, err = myECDSA.NewSem(); err != nil {
		t.Fatalf("生成 [客户端种子] 失败：%v", err)
	}
	if clientPriv, ok = clientSem.GetPriKey().(*ecdsa.PrivateKey); !ok {
		t.Fatalf("[客户端私钥]应为 *ecdsa.PrivateKey")
	}
	if clientPrivPEM, err = clientSem.GetPriKeyPEM(); err != nil {
		t.Fatalf("获取 [客户端私钥] PEM 失败：%v", err)
	}
	clientPrivFile = filesystem.NewFile(filesystem.Rel(testPath, "client.key"))
	if err = clientPrivFile.Write(clientPrivPEM, filesystem.Mode(0600), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 [客户端私钥] PEM 失败：%v", err)
	}
	if clientPub, ok = clientSem.GetPubKey().(*ecdsa.PublicKey); !ok {
		t.Fatalf("[客户端公钥]应为 *ecdsa.PublicKey")
	}
	if clientPubPEM, err = clientSem.GetPubKeyPEM(); err != nil {
		t.Fatalf("获取 [客户端公钥] PEM 失败：%v", err)
	}
	clientPubFile = filesystem.NewFile(filesystem.Rel(testPath, "client.pub"))
	if err = clientPubFile.Write(clientPubPEM, filesystem.Mode(0644), filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY)).GetError(); err != nil {
		t.Fatalf("写入 [客户端公钥] PEM 失败：%v", err)
	}

	// 获取 CA 证书
	if caCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.crt")); !caCrtFile.GetExist() {
		t.Fatalf("[CA 证书]文件不存在：%v", err)
	}
	if caCrtPEM, err = caCrtFile.Read(filesystem.Mode(0644)); err != nil {
		t.Fatalf("读取 [CA 证书] PEM 失败：%v", err)
	}
	if caCrtBlock, _ = pem.Decode(caCrtPEM); caCrtBlock == nil || caCrtBlock.Type != "CERTIFICATE" {
		t.Fatal("[CA 证书] PEM 解码失败或类型不是 CERTIFICATE")
	}
	if caCrt, err = x509.ParseCertificate(caCrtBlock.Bytes); err != nil {
		t.Fatalf("解析 [CA 证书] 失败：%v", err)
	}

	// 获取 CA 私钥
	if caPrivFile = filesystem.NewFile(filesystem.Rel(testPath, "ca.key")); !caPrivFile.GetExist() {
		t.Fatalf("[CA 私钥]文件不存在：%v", err)
	}
	if caPrivPEM, err = caPrivFile.Read(filesystem.Mode(0600)); err != nil {
		t.Fatalf("读取 [CA 私钥] PEM 失败：%v", err)
	}
	if caSem, err = myECDSA.NewSem(myECDSA.PriKeyBytes(caPrivPEM)); err != nil {
		t.Fatalf("生成 [CA 种子] 失败：%v", err)
	}
	if caPriv, ok = caSem.GetPriKey().(*ecdsa.PrivateKey); !ok {
		t.Fatalf("CA 私钥应为 *ecdsa.PrivateKey")
	}

	// 生成 客户端 CSR（模拟客户端）
	clientCSR = &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{"client-ca"},
			CommonName:   "ECDSA P-256 Client CA",
		},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		PublicKey:          clientPub,
	}
	if clientCSRDER, err = x509.CreateCertificateRequest(rand.Reader, clientCSR, clientPriv); err != nil {
		t.Fatalf("生成 [客户端 CSR] 失败：%v", err)
	}

	// 解析 客户端 CSR（模拟服务器）
	if clientCSR, err = x509.ParseCertificateRequest(clientCSRDER); err != nil {
		t.Fatalf("解析 [客户端 CSR] 失败：%v", err)
	}

	// 生成客户端证书并保存
	if sn, err = newSN(); err != nil {
		t.Fatalf("生成序列号失败：%v", err)
	}
	now = time.Now()
	clientTpl = &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"client-ca"},
			CommonName:   "ECDSA P-256 Client CA",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}
	if clientDER, err = x509.CreateCertificate(rand.Reader, clientTpl, caCrt, clientCSR.PublicKey.(*ecdsa.PublicKey), caPriv); err != nil {
		t.Fatalf("生成 [客户端证书] 失败：%v", err)
	}
	if clientCrt, err = x509.ParseCertificate(clientDER); err != nil {
		t.Fatalf("解析 [客户端证书] 失败：%v", err)
	}
	if err = clientCrt.CheckSignatureFrom(caCrt); err != nil {
		t.Fatalf("[客户端证书]校验失败：%v", err)
	}
	clientCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "client.crt"))
	if err = clientCrtFile.Write(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER}),
		filesystem.Mode(0644),
		filesystem.Flag(os.O_CREATE|os.O_TRUNC|os.O_WRONLY),
	).GetError(); err != nil {
		t.Fatalf("写入 [客户端证书] PEM 失败：%v", err)
	}

	if serverCrtFile = filesystem.NewFile(filesystem.Rel(testPath, "server.crt")); !serverCrtFile.GetExist() {
		t.Fatalf("[服务器证书]文件不存在：%v", err)
	}
	if serverCrtPEM, err = serverCrtFile.Read(filesystem.Mode(0644)); err != nil {
		t.Fatalf("读取 [服务器证书] PEM 失败：%v", err)
	}
	if serverCrtBlock, _ = pem.Decode(serverCrtPEM); serverCrtBlock == nil || serverCrtBlock.Type != "CERTIFICATE" {
		t.Fatal("[服务器证书] PEM 解码失败或类型不是 CERTIFICATE")
	}
	if serverCrt, err = x509.ParseCertificate(serverCrtBlock.Bytes); err != nil {
		t.Fatalf("解析 [服务器证书] 失败：%v", err)
	}

	if err = serverCrt.CheckSignatureFrom(caCrt); err != nil {
		t.Fatalf("[服务器证书]校验失败：%v", err)
	}

	return clientCrt, clientSem
}

func TestCA(t *testing.T) {
	newRootCrt(t)
	newServerCrt(t)
	newClientCrt(t)
}

func newSN() (*big.Int, error) { return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)) }
