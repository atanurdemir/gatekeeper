package gates

import (
	"context"

	"github.com/atanurdemir/gatekeeper/src/types"
)

type Gatekeeper interface {
	ProcessRequest(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error)
	SetNext(next Gatekeeper) Gatekeeper
}
