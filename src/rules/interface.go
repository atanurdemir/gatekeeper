package rules

import (
	"context"

	"github.com/atanurdemir/gatekeeper/src/types"
)

type RuleHandler interface {
	Init()
	Instance() RuleHandler
	EvaluateRule(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error)
	UpdateConfig(value interface{})
	SetNext(handler RuleHandler) RuleHandler
	GetRuleState() interface{}
}
