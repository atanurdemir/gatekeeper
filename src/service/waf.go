package service

import (
	"context"

	"github.com/atanurdemir/gatekeeper/src/store"
)

type s_waf struct{}

func NewSWaf() *s_waf {
	return &s_waf{}
}

func (s *s_waf) GetDeniedState(ctx context.Context) (interface{}, error) {
	deniedIPs := store.GetDeniedIPs()

	return deniedIPs, nil
}

var Waf = NewSWaf()
