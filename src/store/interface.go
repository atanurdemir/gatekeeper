package store

type MemoryStore interface {
	Set(actor, key string, value interface{})
	Get(actor, key string) (interface{}, bool)
	GetAllForActor(actor string) map[string]interface{}
}
