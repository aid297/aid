package secret_test

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

	myECDSA "github.com/aid297/aid/v2/secret/asymmetric/ecdsa"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// issueGRPCServerECDSALeaf 由 CA 签发用于 gRPC 服务端 TLS 的 ECDSA 叶子证书（含 localhost / 127.0.0.1 SAN，ServerAuth）。
func issueGRPCServerECDSALeaf(t *testing.T, ca *x509.Certificate, caPriv *ecdsa.PrivateKey) tls.Certificate {
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
			Organization: []string{"aid-secret-test"},
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
	serverRootCert, serverRootPrivKey := newServerECDSACARoot(t)
	serverTLSCert := issueGRPCServerECDSALeaf(t, serverRootCert, serverRootPrivKey)

	serverCreds := credentials.NewServerTLSFromCert(&serverTLSCert)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen：%v", err)
	}
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer := grpc.NewServer(grpc.Creds(serverCreds))
	hs := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(func() { grpcServer.Stop() })

	pool := x509.NewCertPool()
	pool.AddCert(serverRootCert)
	clientTLS := &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}
	clientCreds := credentials.NewTLS(clientTLS)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(clientCreds),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient：%v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	hc := healthpb.NewHealthClient(conn)
	resp, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Health.Check：%v", err)
	}
	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("Health 状态：%v，期望 SERVING", resp.GetStatus())
	}
}
