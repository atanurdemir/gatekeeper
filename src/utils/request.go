package utils

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/atanurdemir/gatekeeper/src/types"
)

func GetIP(r *http.Request) string {
	// Headers that may contain the real IP (usually set by proxies)
	ipHeaders := []string{"X-Forwarded-For", "X-Real-IP"}

	for _, header := range ipHeaders {
		ip := strings.TrimSpace(strings.Split(r.Header.Get(header), ",")[0])
		if ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}

	// If no headers are set or the IP is not valid, fall back to the remote address
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return remoteIP
}

func GetUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(types.UserKey).(string)
	if !ok {
		return "", errors.New("UserID not found in the request context")
	}
	return userID, nil
}
