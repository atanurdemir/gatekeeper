package rules

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atanurdemir/gatekeeper/src/store"
	"github.com/atanurdemir/gatekeeper/src/types"
	"github.com/atanurdemir/gatekeeper/src/utils"
)

type ConfigIPRateLimit struct {
	Limit    int
	Duration int
}

type IPRateLimitRuleHandler struct {
	mu            sync.RWMutex
	next          RuleHandler
	ruleName      string
	ipLimitRate   int
	limitDuration time.Duration
}

var (
	ipRateLimitInstance *IPRateLimitRuleHandler
	ipRateLimitOnce     sync.Once
)

func NewIPRateLimitRuleHandler(config interface{}) *IPRateLimitRuleHandler {
	instance := &IPRateLimitRuleHandler{}
	instance.Init()
	instance.UpdateConfig(config)

	return instance
}

func (h *IPRateLimitRuleHandler) Init() {
	h.ruleName = "IPRateLimitRule"
	h.ipLimitRate = 10
	h.limitDuration = 3600 * time.Second
}

func (h *IPRateLimitRuleHandler) Instance() RuleHandler {
	ipRateLimitOnce.Do(func() {
		ipRateLimitInstance = h
	})

	return ipRateLimitInstance
}

func (h *IPRateLimitRuleHandler) UpdateConfig(value interface{}) {
	if value == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if configMap, ok := value.(map[string]interface{}); ok {
		if limit, ok := configMap["limit"].(int); ok {
			h.ipLimitRate = limit
		}

		if duration, ok := configMap["duration"].(int); ok {
			h.limitDuration = time.Duration(duration) * time.Second
		}
	}
}

func (h *IPRateLimitRuleHandler) EvaluateRule(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error) {
	requestIP := request.IP
	if requestIP == "" {
		return &types.CheckResponse{Status: false, Reason: "IP address not found"}, nil
	}

	h.mu.RLock()
	rate := h.ipLimitRate
	duration := h.limitDuration
	h.mu.RUnlock()

	memoryStore := store.NewMemoryStore()
	_, err := utils.UpdateRequestCount(memoryStore, h.ruleName, requestIP, rate, duration)
	if err != nil {
		store.AddDeniedIP(requestIP, h.ruleName, duration)

		reason := fmt.Sprintf("IP rate limit exceeded. Limit: %d requests per %v seconds.", rate, duration.Seconds())

		return &types.CheckResponse{Status: false, Reason: reason}, nil
	}

	if h.next != nil {
		return h.next.EvaluateRule(ctx, request)
	}

	return &types.CheckResponse{Status: true}, nil
}

func (h *IPRateLimitRuleHandler) SetNext(next RuleHandler) RuleHandler {
	h.next = next
	return next
}

func (h *IPRateLimitRuleHandler) GetRuleState() interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"limit":    h.ipLimitRate,
		"duration": h.limitDuration,
	}
}
