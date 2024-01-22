package gates

import (
	"context"
	"sync"

	"github.com/atanurdemir/gatekeeper/src/rules"
	"github.com/atanurdemir/gatekeeper/src/types"
)

type CustomGate struct {
	handlers map[string]*rules.RuleHandler
	next     Gatekeeper
}

func (t *CustomGate) evaluateResponses(responses []*types.CheckResponse) *types.CheckResponse {
	for _, response := range responses {
		if response == nil || !response.Status {
			return response
		}
	}
	return &types.CheckResponse{Status: true}
}

func (t *CustomGate) ProcessRequest(ctx context.Context, request types.RequestContext) (*types.CheckResponse, error) {
	var wg sync.WaitGroup
	wg.Add(len(t.handlers))

	chResponse := make(chan *types.CheckResponse, len(t.handlers))

	for _, handler := range t.handlers {
		go func(h *rules.RuleHandler) {
			defer wg.Done()
			if response, err := (*h).EvaluateRule(ctx, request); err == nil {
				chResponse <- response
			} else {
				// Log or handle the error from EvaluateRule
				chResponse <- &types.CheckResponse{Status: false, Reason: err.Error()}
			}
		}(handler)
	}

	wg.Wait()
	close(chResponse)

	responses := make([]*types.CheckResponse, 0, len(t.handlers))
	for response := range chResponse {
		responses = append(responses, response)
	}

	if finalResponse := t.evaluateResponses(responses); !finalResponse.Status {
		return finalResponse, nil
	}

	if t.next != nil {
		return t.next.ProcessRequest(ctx, request)
	}

	return &types.CheckResponse{Status: true}, nil
}

func (t *CustomGate) SetNext(next Gatekeeper) Gatekeeper {
	t.next = next
	return next
}
