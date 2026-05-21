package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func getCAPool() *x509.CertPool {
	caFile := "ca.crt"
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("读取 CA %s: %v", caFile, err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPEM) {
		log.Fatalf("解析 CA PEM 失败: %s", caFile)
	}
	return caPool
}

func main() {
	r := gin.New()

	caPool := getCAPool() // 获取 CA 证书池

	clientATLSConfig := &tls.Config{
		ClientCAs:  caPool,                         // 设置 CA 证书池(所有入站请求都必须在此信任链中)
		ClientAuth: tls.RequireAndVerifyClientCert, // 要求客户端证书(所有入站请求都必须提供有效的客户端证书)
		MinVersion: tls.VersionTLS12,               // 设置最小 TLS 版本
	}

	httpServer := &http.Server{
		Addr:      ":8443",
		Handler:   r,
		TLSConfig: clientATLSConfig,
	}
	httpServer.ListenAndServeTLS("client-a.crt", "client-a.key")

	clientBCrt := "client-b.crt"
	clientBPriv := "client-b.key"
	clientBKeyPair, err := tls.LoadX509KeyPair(clientBCrt, clientBPriv)
	if err != nil {
		log.Fatalf("加载节点证书 %s + %s: %v", clientBCrt, clientBPriv, err)
	}

	clientBTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{clientBKeyPair}, // 设置客户端证书
		RootCAs:      caPool,                            // 设置 CA 证书池(所有出站请求都必须在此信任链中)
		MinVersion:   tls.VersionTLS12,                  // 设置最小 TLS 版本
	}

	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: clientBTLSConfig},
		Timeout:   10 * time.Second,
	}

	httpClient.Get("https://localhost:8443/")
}
