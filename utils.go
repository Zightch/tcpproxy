package main

import (
	"net"
	"regexp"
	"strconv"
	"strings"
)

func splitIPv6Port(str string) (string, uint16) {
	if !isIPv6Port(str) {
		return "", 0
	}
	re := regexp.MustCompile("^\\[(.*)]:([0-9]+)$")
	list := re.FindStringSubmatch(str)
	port, _ := strconv.Atoi(list[2])
	return list[1], uint16(port)
}

func splitIPv4Port(str string) (string, uint16) {
	if !isIPv4Port(str) {
		return "", 0
	}
	re := regexp.MustCompile("^(.*):([0-9]+)$")
	list := re.FindStringSubmatch(str)
	port, _ := strconv.Atoi(list[2])
	return list[1], uint16(port)
}

func isIPv6Port(str string) bool {
	re := regexp.MustCompile("^\\[(.*)]:([0-9]+)$")
	if !re.MatchString(str) {
		return false
	}
	list := re.FindStringSubmatch(str)
	if len(list) != 3 {
		return false
	}
	if !isIPv6(list[1]) {
		return false
	}
	num, err := strconv.Atoi(list[2])
	if err != nil {
		return false
	}
	if num <= 0 || num > 65535 {
		return false
	}
	return true
}

func isIPv4Port(str string) bool {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		return false
	}
	if !isIPv4(list[0]) {
		return false
	}
	num, err := strconv.Atoi(list[1])
	if err != nil {
		return false
	}
	if num <= 0 || num > 65535 {
		return false
	}
	return true
}

func isAnyIPPort(str string) bool {
	re := regexp.MustCompile("^:([0-9]+)$")
	list := re.FindStringSubmatch(str)
	if len(list) != 2 {
		return false
	}
	num, err := strconv.Atoi(list[1])
	if err != nil {
		return false
	}
	if num <= 0 || num > 65535 {
		return false
	}
	return true
}

func isIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To16() != nil && ip.To4() == nil
}

func isIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To16() == nil && ip.To4() != nil
}
