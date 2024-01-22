package rules

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atanurdemir/gatekeeper/src/store"
	"github.com/atanurdemir/gatekeeper/src/types"
)

type ConfigIPRestrictionRuleHandler struct {
	AllowedIPs []string
}

type IPRestrictionRuleHandler struct {
	mu         sync.RWMutex
	next       RuleHandler
	ruleName   string
	allowedIPs map[string]struct{}
}

var (
	ipRestrictionInstance *IPRestrictionRuleHandler
	ipRestrictionOnce     sync.Once
)

func NewIPRestrictionRuleHandler(config interface{}) *IPRestrictionRuleHandler {
	instance := &IPRestrictionRuleHandler{}
	instance.Init()
	instance.UpdateConfig(config)

	return instance
}

func (h *IPRestrictionRuleHandler) Init() {
	h.ruleName = "IPRestrictionRule"
}

func (h *IPRestrictionRuleHandler) Instance() RuleHandler {
	ipRestrictionOnce.Do(func() {
		ipRestrictionInstance = h
	})

	return ipRestrictionInstance
}

func (h *IPRestrictionRuleHandler) UpdateConfig(value interface{}) {
	if value == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.allowedIPs = make(map[string]struct{})
	if configMap, ok := value.(map[string]interface{}); ok {
		if allowedIPs, ok := configMap["allowed_ips"].([]interface{}); ok {
			for _, ip := range allowedIPs {
				if ipStr, ok := ip.(string); ok {
					h.allowedIPs[ipStr] = struct{}{}
				}
			}
		}
	}
}

func (h *IPRestrictionRuleHandler) EvaluateRule(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error) {
	if request.IP == "" {
		return &types.CheckResponse{Status: false, Reason: "IP address not found"}, nil
	}

	if h.isAllowedIP(request.IP) {
		if h.next != nil {
			return h.next.EvaluateRule(ctx, request)
		}
		return &types.CheckResponse{Status: true}, nil
	}

	store.AddDeniedIP(request.IP, h.ruleName, time.Hour)
	return &types.CheckResponse{Status: false, Reason: fmt.Sprintf("Request from IP '%s' is not allowed.", request.IP)}, nil
}

func (h *IPRestrictionRuleHandler) isAllowedIP(ip string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, allowed := h.allowedIPs[ip]
	return allowed
}

func (h *IPRestrictionRuleHandler) SetNext(next RuleHandler) RuleHandler {
	h.next = next
	return next
}

func (h *IPRestrictionRuleHandler) GetRuleState() interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"allowedIPs": h.allowedIPs,
	}
}
