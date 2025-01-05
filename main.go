package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("Usage: " + args[0] + " <config_file_path>")
	}
	config(args[1])

	if CONFIG.Proxy.Server == "" {
		panic("Server address is empty")
	}

	ip, port := splitIPv6Port(CONFIG.Proxy.Local)
	if port == 0 {
		ip, port = splitIPv4Port(CONFIG.Proxy.Local)
		if port == 0 {
			if isAnyIPPort(CONFIG.Proxy.Local) {
				ip = ""
				num, _ := strconv.Atoi(CONFIG.Proxy.Local[1:])
				port = uint16(num)
			} else {
				panic("Invalid local address: " + CONFIG.Proxy.Local)
			}
		}
	}
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

	var listener net.Listener
	var err error

	// 根据LOCAL_SSL_CONF.enable判断使用tls还是普通tcp
	if LOCAL_SSL_CONF.Enable {
		// 创建TLS配置
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{loadCertificate(LOCAL_SSL_CONF.crt, LOCAL_SSL_CONF.key)},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			RootCAs:      loadRootCAs(LOCAL_SSL_CONF.ca),
			MinVersion:   tls.VersionTLS13,
		}

		// 使用TLS监听
		listener, err = tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			panic(err)
		}
		fmt.Println("Listening on TLS: ", addr)
	} else {
		// 使用普通TCP监听
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
		}
		fmt.Println("Listening on TCP: ", addr)
	}

	defer listener.Close()
	// 接受新的连接并处理它们...
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn) // 假设你有一个handleConnection函数来处理每个连接
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	serverAddr := CONFIG.Proxy.Server
	if serverAddr == "" {
		panic("Server address not configured")
	}

	var serverConn net.Conn
	var err error

	// 根据SERVER_SSL_CONF.enable判断使用tls还是普通tcp, 建立对server的连接
	if SERVER_SSL_CONF.Enable {
		// 创建TLS配置
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{loadCertificate(SERVER_SSL_CONF.crt, SERVER_SSL_CONF.key)},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			RootCAs:      loadRootCAs(SERVER_SSL_CONF.ca),
			MinVersion:   tls.VersionTLS13,
		}
		// 使用TLS连接到服务器
		serverConn, err = tls.Dial("tcp", serverAddr, tlsConfig)
		if err != nil {
			fmt.Println("Failed to establish TLS connection to server: ", err)
			return
		}
	} else {
		// 使用普通TCP连接到服务器
		serverConn, err = net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Failed to establish TCP connection to server: ", err)
			return
		}
	}
	defer serverConn.Close()

	done := make(chan interface{})

	// 客户端 -> 服务器
	go func() {
		_, err := io.Copy(serverConn, clientConn)
		if err != nil {
			fmt.Println("Error copying data from client to server: ", err)
		}
		close(done)
	}()

	// 服务器 -> 客户端
	go func() {
		_, err := io.Copy(clientConn, serverConn)
		if err != nil {
			fmt.Println("Error copying data from server to client: ", err)
		}
		close(done)
	}()

	// 等待两个goroutine都完成
	<-done
	<-done
}
