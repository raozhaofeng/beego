package utils

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// GetUserRealIP 获取用户真实IP
func GetUserRealIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "0.0.0.0"
	}

	if net.ParseIP(ip) != nil {
		return ip
	}
	return "0.0.0.0"
}

// PasswordEncrypt 密码加密
func PasswordEncrypt(password string) string {
	if password == "" {
		return ""
	}

	has1 := md5.Sum([]byte(password))
	has2 := md5.Sum([]byte(fmt.Sprintf("%x", has1)))
	return fmt.Sprintf("%x", has2)
}
