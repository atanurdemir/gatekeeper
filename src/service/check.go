package service

import (
	"context"

	"github.com/atanurdemir/gatekeeper/src/config"
	"github.com/atanurdemir/gatekeeper/src/gates"
	"github.com/atanurdemir/gatekeeper/src/models"
	"github.com/atanurdemir/gatekeeper/src/types"
)

type s_check struct {
	appConfig *models.AppConfig
}

func NewSCheck() *s_check {
	return &s_check{appConfig: &config.GatekeeperConfig}
}

func (t *s_check) ProcessRequest(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error) {
	dynamicRules := t.findMatchingRules(request.Method, request.Path)
	strategy := gates.Factory.New(dynamicRules)
	return strategy.ProcessRequest(ctx, request)
}

func (t *s_check) findMatchingRules(method, path string) []models.Rule {
	for _, route := range t.appConfig.Gates {
		if route.Path == path && route.Method == method {
			return route.Rules
		}
	}
	return []models.Rule{}
}

var Check = NewSCheck()
