package store

import (
	"time"
)

func CleanupExpiredBans() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		removeExpiredBans()
	}
}

func removeExpiredBans() {
	now := time.Now()
	deniedIPStore.Lock()
	defer deniedIPStore.Unlock()

	for key, info := range deniedIPStore.IPs {
		if now.After(info.ExpiresAt) {
			delete(deniedIPStore.IPs, key)
		}
	}
}
