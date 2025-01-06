package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// 加载证书和私钥
func loadCertificate(crt, key string) tls.Certificate {
	// 如果crt与key都为空字符串, 返回一个空的证书
	if crt == "" && key == "" {
		return tls.Certificate{}
	}
	cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return tls.Certificate{}
	}
	return cert
}

// 加载根证书颁发机构
func loadRootCAs(ca string) *x509.CertPool {
	// 尝试加载系统默认的CA证书
	systemPool, err := x509.SystemCertPool()
	if err != nil || systemPool == nil {
		// 如果无法加载系统证书池，则创建一个新的空证书池
		systemPool = x509.NewCertPool()
	}
	if ca != "" {
		// 如果提供了额外的CA证书，则尝试将其添加到证书池中
		if ok := systemPool.AppendCertsFromPEM([]byte(ca)); !ok {
			fmt.Println("Failed to parse root CA certificate")
		}
	}
	return systemPool
}
