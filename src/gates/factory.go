package gates

import (
	"log"

	"github.com/atanurdemir/gatekeeper/src/models"
	"github.com/atanurdemir/gatekeeper/src/rules"
)

type GateFactory struct{}

func (t *GateFactory) New(dynamicRules []models.Rule) Gatekeeper {
	if len(dynamicRules) == 0 {
		return t.defaultGate()
	}
	return t.createGateWithRules(dynamicRules)
}

func (t *GateFactory) defaultRule(key string, value interface{}, singleton bool) *rules.RuleHandler {
	handler, exists := rules.GetRuleHandler(key)
	if !exists {
		log.Printf("RuleHandler: %s is not exists", key)
	}

	if singleton {
		h := (*handler).Instance()
		handler = &h
	}

	(*handler).UpdateConfig(value)
	return handler
}

func (t *GateFactory) defaultGate() Gatekeeper {
	defaultRules := []models.Rule{
		{Name: "ip_restriction", Config: nil},
		{Name: "ip_rate", Config: rules.ConfigIPRateLimit{Limit: 30, Duration: 3600}},
	}
	return t.createGateWithRules(defaultRules)
}

func (t *GateFactory) createGateWithRules(rulesConfig []models.Rule) Gatekeeper {
	custom := &CustomGate{
		handlers: make(map[string]*rules.RuleHandler),
	}

	for _, rule := range rulesConfig {
		handler := t.defaultRule(rule.Name, rule.Config, true)
		if handler != nil {
			custom.handlers[rule.Name] = handler
		}
	}

	return custom
}

var Factory = &GateFactory{}
