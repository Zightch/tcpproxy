package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
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

	if CONFIG.Proxy.Local == "" {
		panic("Local address is empty")
	}

	var listener net.Listener
	var err error

	// 根据LOCAL_SSL_CONF.enable判断使用tls还是普通tcp
	if LOCAL_SSL_CONF.Enable {
		// 使用TLS监听
		listener, err = tls.Listen("tcp", CONFIG.Proxy.Local, LOCAL_TLS_CONF)
		if err != nil {
			panic(err)
		}
		fmt.Println("Listening on TLS:", CONFIG.Proxy.Local)
	} else {
		// 使用普通TCP监听
		listener, err = net.Listen("tcp", CONFIG.Proxy.Local)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		fmt.Println("Listening on TCP:", CONFIG.Proxy.Local)
	}

	defer listener.Close()
	// 接受新的连接并处理它们...
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
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

	dialer := net.Dialer{
		Timeout: time.Second * time.Duration(CONFIG.Proxy.ConnectWaitTime), // 设置连接超时时间
	}

	// 根据SERVER_SSL_CONF.enable判断使用tls还是普通tcp, 建立对server的连接
	if SERVER_SSL_CONF.Enable {
		// 使用TLS连接到服务器
		serverConn, err = tls.DialWithDialer(&dialer, "tcp", serverAddr, SERVER_TLS_CONF)
		if err != nil {
			fmt.Println("Failed to establish TLS connection to server:", err)
			return
		}
	} else {
		// 使用普通TCP连接到服务器
		serverConn, err = dialer.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Failed to establish TCP connection to server:", err)
			return
		}
	}
	defer serverConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 -> 服务器
	go func() {
		_, err := io.Copy(serverConn, clientConn)
		if err != nil {
			fmt.Println("Error copying data from client to server:", err)
		}
		serverConn.Close()
		wg.Done()
	}()

	// 服务器 -> 客户端
	go func() {
		_, err := io.Copy(clientConn, serverConn)
		if err != nil {
			fmt.Println("Error copying data from server to client:", err)
		}
		serverConn.Close()
		wg.Done()
	}()

	// 等待两个goroutine都完成
	wg.Wait()
}
