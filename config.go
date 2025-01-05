package main

import (
	"encoding/json"
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
}
