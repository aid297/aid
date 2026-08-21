package secrets_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/aid297/aid/v3/secrets"
	myECDSA "github.com/aid297/aid/v3/secrets/asymmetric/ecdsa"
)

// tlsCertFromSem 将 *x509.Certificate 与 Semen 私钥组装为 tls.Certificate，chain 为附加的签发链（如根 CA）。
func tlsCertFromSem(t *testing.T, leaf *x509.Certificate, sem secrets.Semen, chain ...*x509.Certificate) tls.Certificate {
	t.Helper()
	priv, ok := sem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || priv == nil {
		t.Fatal("私钥应为 *ecdsa.PrivateKey")
	}
	raw := make([][]byte, 0, 1+len(chain))
	raw = append(raw, leaf.Raw)
	for _, c := range chain {
		raw = append(raw, c.Raw)
	}
	return tls.Certificate{Certificate: raw, PrivateKey: priv}
}

// certPoolFromRoot 用根 CA 证书构建 TLS 信任池。
func certPoolFromRoot(t *testing.T, root *x509.Certificate) *x509.CertPool {
	t.Helper()
	pool := x509.NewCertPool()
	pool.AddCert(root)
	return pool
}

// TestGRPC 使用 test.data/ecdsa 下的 ca / server / client 证书与种子，建立双向 TLS 的 gRPC 连接并做健康检查。
func TestGRPC(t *testing.T) {
	caCrt, _ := newRootCrt(t)
	serverCrt, serverSem := newServerCrt(t)

	caPool := certPoolFromRoot(t, caCrt)
	serverTLSCert := tlsCertFromSem(t, serverCrt, serverSem, caCrt)

	serverTLS := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
		MinVersion:   tls.VersionTLS12,
	}
	serverCreds := credentials.NewTLS(serverTLS)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听端口失败：%v", err)
	}
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer := grpc.NewServer(grpc.Creds(serverCreds))
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(func() { grpcServer.Stop() })

	clientCrt, clientSem := newClientCrt(t)
	clientTLSCert := tlsCertFromSem(t, clientCrt, clientSem, caCrt)

	clientTLS := &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{clientTLSCert},
		ServerName:   "localhost",
		MinVersion:   tls.VersionTLS12,
	}
	clientCreds := credentials.NewTLS(clientTLS)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(clientCreds))
	if err != nil {
		t.Fatalf("创建 [gRPC 客户端] 失败：%v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	hc := healthpb.NewHealthClient(conn)
	resp, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("[gRPC 客户端] 健康检查失败：%v", err)
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("[gRPC 客户端] 健康状态错误：%v，期望 SERVING", resp.GetStatus())
	}
}

// newServerCert 由 CA 签发用于 gRPC 服务端 TLS 的 ECDSA 叶子证书（含 localhost / 127.0.0.1 SAN，ServerAuth）。
func newServerCert(t *testing.T, ca *x509.Certificate, caPriv *ecdsa.PrivateKey) tls.Certificate {
	t.Helper()
	sem, err := myECDSA.NewSem()
	if err != nil {
		t.Fatalf("生成服务端 TLS 密钥失败：%v", err)
	}
	leafPriv, ok := sem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || leafPriv == nil {
		t.Fatal("服务端叶子私钥应为 *ecdsa.PrivateKey")
	}
	leafPub := &leafPriv.PublicKey

	sn, err := newSN()
	if err != nil {
		t.Fatalf("序列号：%v", err)
	}
	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{"aid-secrets-test"},
			CommonName:   "grpc-ecdsa-server",
		},
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.AddDate(0, 3, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, tpl, ca, leafPub, caPriv)
	if err != nil {
		t.Fatalf("签发服务端 TLS 证书失败：%v", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{leafDER, ca.Raw},
		PrivateKey:  leafPriv,
	}
}

// TestECDSA_gRPCTLS_ClientTrustsCA 模拟客户端仅持有 CA 根证书：建立 gRPC over TLS，调用标准 Health.Check。
func TestECDSA_gRPCTLS_ClientTrustsCA(t *testing.T) {
	var (
		err         error
		caCrt       *x509.Certificate
		caSem       secrets.Semen
		caPrivKey   *ecdsa.PrivateKey
		serverCrt   tls.Certificate
		serverCreds credentials.TransportCredentials
		lis         net.Listener
		grpcServer  *grpc.Server
		hs          *health.Server
	)

	caCrt, caSem = newRootCrt(t)
	caPrivKey, ok := caSem.GetPriKey().(*ecdsa.PrivateKey)
	if !ok || caPrivKey == nil {
		t.Fatal("CA 私钥应为 *ecdsa.PrivateKey")
	}
	serverCrt = newServerCert(t, caCrt, caPrivKey)

	serverCreds = credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{serverCrt}})

	if lis, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
		t.Fatalf("监听端口失败：%v", err)
	}
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer = grpc.NewServer(grpc.Creds(serverCreds))
	hs = health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(func() { grpcServer.Stop() })

	var (
		pool        *x509.CertPool
		clientCreds credentials.TransportCredentials
		conn        *grpc.ClientConn
		hc          healthpb.HealthClient
		resp        *healthpb.HealthCheckResponse
	)

	pool = x509.NewCertPool()
	pool.AddCert(caCrt)
	clientCreds = credentials.NewTLS(&tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if conn, err = grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(clientCreds)); err != nil {
		t.Fatalf("创建 [gRPC 客户端] 失败：%v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	hc = healthpb.NewHealthClient(conn)
	if resp, err = hc.Check(ctx, &healthpb.HealthCheckRequest{}); err != nil {
		t.Fatalf("[gRPC 客户端] 健康检查失败：%v", err)
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("[gRPC 客户端] 健康状态错误：%v，期望 SERVING", resp.GetStatus())
	}
}
