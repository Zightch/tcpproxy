package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

type Config struct {
	Proxy struct {
		Local           string `json:"local"`
		Server          string `json:"server"`
		ConnectWaitTime int    `json:"connect_wait_time"`
	} `json:"proxy"`
	LocalSSL struct {
		Enable      bool   `json:"enable"`
		CrtFilePath string `json:"crt_file_path"`
		KeyFilePath string `json:"key_file_path"`
		CAFilePath  string `json:"ca_file_path"`
	} `json:"local_ssl"`
	ServerSSL struct {
		Enable      bool   `json:"enable"`
		CrtFilePath string `json:"crt_file_path"`
		KeyFilePath string `json:"key_file_path"`
		CAFilePath  string `json:"ca_file_path"`
	} `json:"server_ssl"`
}

type SSLConfig struct {
	Enable bool
	crt    string
	key    string
	ca     string
}

var CONFIG = Config{}
var LOCAL_SSL_CONF = SSLConfig{}
var SERVER_SSL_CONF = SSLConfig{}

var LOCAL_TLS_CONF = &tls.Config{
	ClientAuth: tls.VerifyClientCertIfGiven,
	MinVersion: tls.VersionTLS13,
}
var SERVER_TLS_CONF = &tls.Config{
	ClientAuth: tls.VerifyClientCertIfGiven,
	MinVersion: tls.VersionTLS13,
}

func config(configFilePath string) {
	fileData, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileData, &CONFIG)
	if err != nil {
		panic(err)
	}
	if CONFIG.LocalSSL.Enable {
		LOCAL_SSL_CONF.Enable = CONFIG.LocalSSL.Enable
		if CONFIG.LocalSSL.CrtFilePath != "" && CONFIG.LocalSSL.KeyFilePath != "" {
			crt, err := os.ReadFile(CONFIG.LocalSSL.CrtFilePath)
			if err != nil {
				panic(err)
			}
			LOCAL_SSL_CONF.crt = string(crt)
			key, err := os.ReadFile(CONFIG.LocalSSL.KeyFilePath)
			if err != nil {
				panic(err)
			}
			LOCAL_SSL_CONF.key = string(key)
		} else if CONFIG.LocalSSL.CrtFilePath != "" || CONFIG.LocalSSL.KeyFilePath != "" {
			if CONFIG.LocalSSL.CrtFilePath == "" {
				panic("Local crt file path is empty, crt and key must meanwhile set")
			} else if CONFIG.LocalSSL.KeyFilePath == "" {
				panic("Local key file path is empty, crt and key must meanwhile set")
			}
		} else {
			LOCAL_SSL_CONF.crt = ""
			LOCAL_SSL_CONF.key = ""
		}
		if CONFIG.LocalSSL.CAFilePath != "" {
			ca, err := os.ReadFile(CONFIG.LocalSSL.CAFilePath)
			if err != nil {
				panic(err)
			}
			LOCAL_SSL_CONF.ca = string(ca)
		} else {
			LOCAL_SSL_CONF.ca = ""
		}
	} else {
		LOCAL_SSL_CONF.Enable = false
	}
	if CONFIG.ServerSSL.Enable {
		SERVER_SSL_CONF.Enable = CONFIG.ServerSSL.Enable
		if CONFIG.ServerSSL.CrtFilePath != "" && CONFIG.ServerSSL.KeyFilePath != "" {
			crt, err := os.ReadFile(CONFIG.ServerSSL.CrtFilePath)
			if err != nil {
				panic(err)
			}
			SERVER_SSL_CONF.crt = string(crt)
			key, err := os.ReadFile(CONFIG.ServerSSL.KeyFilePath)
			if err != nil {
				panic(err)
			}
			SERVER_SSL_CONF.key = string(key)
		} else if CONFIG.ServerSSL.CrtFilePath != "" || CONFIG.ServerSSL.KeyFilePath != "" {
			if CONFIG.ServerSSL.CrtFilePath == "" {
				panic("Server crt file path is empty, crt and key must meanwhile set")
			} else if CONFIG.ServerSSL.KeyFilePath == "" {
				panic("Server key file path is empty, crt and key must meanwhile set")
			}
		} else {
			SERVER_SSL_CONF.crt = ""
			SERVER_SSL_CONF.key = ""
		}
		if CONFIG.ServerSSL.CAFilePath != "" {
			ca, err := os.ReadFile(CONFIG.ServerSSL.CAFilePath)
			if err != nil {
				panic(err)
			}
			SERVER_SSL_CONF.ca = string(ca)
		} else {
			SERVER_SSL_CONF.ca = ""
		}
	} else {
		SERVER_SSL_CONF.Enable = false
	}
	sslConfig()
}

func sslConfig() {
	if LOCAL_SSL_CONF.Enable {
		if LOCAL_SSL_CONF.crt != "" && LOCAL_SSL_CONF.key != "" {
			cert, err := tls.X509KeyPair([]byte(LOCAL_SSL_CONF.crt), []byte(LOCAL_SSL_CONF.key))
			if err != nil {
				str := fmt.Sprintf("Error loading local certificate: %v", err)
				panic(str)
			}
			LOCAL_TLS_CONF.Certificates = []tls.Certificate{cert}
		}
		if LOCAL_SSL_CONF.ca != "" {
			block, _ := pem.Decode([]byte(LOCAL_SSL_CONF.ca))
			if block == nil || block.Type != "CERTIFICATE" {
				panic("Error loading local certificate: Unable to decode PEM block")
			}
			crt, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				str := fmt.Sprintf("Error loading local certificate: %v", err)
				panic(str)
			}
			systemPool := x509.NewCertPool()
			systemPool.AddCert(crt)
			LOCAL_TLS_CONF.ClientCAs = systemPool
			LOCAL_TLS_CONF.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	if SERVER_SSL_CONF.Enable {
		if SERVER_SSL_CONF.crt != "" && SERVER_SSL_CONF.key != "" {
			cert, err := tls.X509KeyPair([]byte(SERVER_SSL_CONF.crt), []byte(SERVER_SSL_CONF.key))
			if err != nil {
				str := fmt.Sprintf("Error loading server certificate: %v", err)
				panic(str)
			}
			SERVER_TLS_CONF.Certificates = []tls.Certificate{cert}
		}
		// 尝试加载系统默认的CA证书
		systemPool, err := x509.SystemCertPool()
		if err != nil || systemPool == nil {
			// 如果无法加载系统证书池，则创建一个新的空证书池
			systemPool = x509.NewCertPool()
		}
		if SERVER_SSL_CONF.ca != "" {
			block, _ := pem.Decode([]byte(SERVER_SSL_CONF.ca))
			if block == nil || block.Type != "CERTIFICATE" {
				panic("Error loading server certificate: Unable to decode PEM block")
			}
			crt, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				str := fmt.Sprintf("Error loading server certificate: %v", err)
				panic(str)
			}
			systemPool.AddCert(crt)
			SERVER_TLS_CONF.ClientAuth = tls.RequireAndVerifyClientCert
		}
		SERVER_TLS_CONF.RootCAs = systemPool
	}
}
