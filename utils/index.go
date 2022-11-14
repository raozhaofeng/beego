package utils

import (
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
