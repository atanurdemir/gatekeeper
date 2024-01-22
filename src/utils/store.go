package utils

import (
	"errors"
	"time"

	"github.com/atanurdemir/gatekeeper/src/store"
)

// RequestInfo holds the count and the timestamp for a request.
type RequestInfo struct {
	Count     int
	Timestamp time.Time
}

func UpdateRequestCount(memoryStore store.MemoryStore, actor, key string, limit int, duration time.Duration) (int, error) {
	now := time.Now()
	item, exists := memoryStore.Get(actor, key)
	info := &RequestInfo{}

	if exists {
		var ok bool
		if info, ok = item.(*RequestInfo); !ok || now.Sub(info.Timestamp) > duration {
			info = &RequestInfo{Count: 1, Timestamp: now}
		} else {
			info.Count++
		}
	} else {
		info = &RequestInfo{Count: 1, Timestamp: now}
	}

	memoryStore.Set(actor, key, info)

	if info.Count > limit {
		return info.Count, errors.New("rate limit exceeded")
	}
	return info.Count, nil
}
