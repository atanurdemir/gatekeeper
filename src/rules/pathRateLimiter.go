package rules

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/atanurdemir/gatekeeper/src/store"
	"github.com/atanurdemir/gatekeeper/src/types"
	"github.com/atanurdemir/gatekeeper/src/utils"
)

type PathRateLimitRuleHandler struct {
	mu            sync.RWMutex
	next          RuleHandler
	ruleName      string
	pathLimits    map[string]RateLimitConfig
	limitDuration time.Duration
}

type RateLimitConfig struct {
	Limit    int
	Duration time.Duration
}

var (
	pathRateLimitInstance *PathRateLimitRuleHandler
	pathRateLimitOnce     sync.Once
)

func NewPathRateLimitRuleHandler(config map[string]RateLimitConfig) *PathRateLimitRuleHandler {
	instance := &PathRateLimitRuleHandler{}
	instance.Init()
	instance.UpdateConfig(config)

	return instance
}

func (h *PathRateLimitRuleHandler) Init() {
	h.ruleName = "PathRateLimitRule"
	h.pathLimits = make(map[string]RateLimitConfig)
	h.limitDuration = 60 * time.Second
}

func (h *PathRateLimitRuleHandler) Instance() RuleHandler {
	pathRateLimitOnce.Do(func() {
		pathRateLimitInstance = h
	})

	return pathRateLimitInstance
}
func (h *PathRateLimitRuleHandler) UpdateConfig(value interface{}) {
	if configMap, ok := value.(map[string]interface{}); ok {
		h.mu.Lock()
		defer h.mu.Unlock()

		if path, ok := configMap["path"].(string); ok {
			newConfig := RateLimitConfig{}

			if limit, ok := configMap["limit"].(int); ok {
				newConfig.Limit = limit
			}
			if duration, ok := configMap["duration"].(int); ok {
				newConfig.Duration = time.Duration(duration) * time.Second
			}

			h.pathLimits[path] = newConfig
		} else {
			log.Printf("PathRateLimitRuleHandler config update: 'path' key missing or invalid")
		}
	}
}

func (h *PathRateLimitRuleHandler) EvaluateRule(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error) {
	if request.IP == "" {
		return &types.CheckResponse{Status: false, Reason: "IP address not found"}, nil
	}

	h.mu.RLock()
	limitConfig, ok := h.pathLimits[request.Path]
	h.mu.RUnlock()

	if ok {
		memoryStore := store.NewMemoryStore()
		count, err := utils.UpdateRequestCount(memoryStore, h.ruleName, request.IP, limitConfig.Limit, limitConfig.Duration)

		if err != nil || count > limitConfig.Limit {
			store.AddDeniedIP(request.IP, h.ruleName, limitConfig.Duration)
			reason := fmt.Sprintf("Path rate limit exceeded for %s. Limit: %d requests per %v.", request.Path, limitConfig.Limit, limitConfig.Duration)
			return &types.CheckResponse{Status: false, Reason: reason}, nil
		}
	}

	if h.next != nil {
		return h.next.EvaluateRule(ctx, request)
	}

	return &types.CheckResponse{Status: true}, nil
}

func (h *PathRateLimitRuleHandler) SetNext(next RuleHandler) RuleHandler {
	h.next = next
	return next
}

func (h *PathRateLimitRuleHandler) GetRuleState() interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"pathLimits": h.pathLimits,
	}
}
