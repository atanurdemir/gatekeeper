package types

type RequestContext struct {
	Path   string
	Method string
	IP     string
	UserID string
}

type CheckResponse struct {
	Status bool   `json:"status"`
	Reason string `json:"reason"`
}
