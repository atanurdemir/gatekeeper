package store

import (
	"sync"
	"time"
)

type DeniedIPInfo struct {
	IP        string
	Reason    string
	ExpiresAt time.Time
}

type DeniedIPStore struct {
	sync.RWMutex
	IPs map[string]DeniedIPInfo
}

var deniedIPStore = &DeniedIPStore{
	IPs: make(map[string]DeniedIPInfo),
}

func AddDeniedIP(ip, reason string, banDuration time.Duration) {
	expirationTime := time.Now().Add(banDuration)

	deniedIPStore.Lock()
	defer deniedIPStore.Unlock()

	key := ip + "|" + reason
	info, exists := deniedIPStore.IPs[key]

	// If the entry does not exist or has expired, add or update it
	if !exists || info.ExpiresAt.Before(time.Now()) {
		deniedIPStore.IPs[key] = DeniedIPInfo{IP: ip, Reason: reason, ExpiresAt: expirationTime}
	}
}
func GetDeniedIPs() []DeniedIPInfo {
	deniedIPStore.RLock()
	defer deniedIPStore.RUnlock()

	ips := make([]DeniedIPInfo, 0, len(deniedIPStore.IPs))
	for _, info := range deniedIPStore.IPs {
		if info.ExpiresAt.After(time.Now()) {
			ips = append(ips, info)
		}
	}
	return ips
}

func IsIPBanned(ip string) bool {
	deniedIPStore.RLock()
	defer deniedIPStore.RUnlock()

	for _, info := range deniedIPStore.IPs {
		if info.IP == ip && info.ExpiresAt.After(time.Now()) {
			return true
		}
	}
	return false
}
