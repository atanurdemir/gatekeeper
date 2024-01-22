package rules

import (
	"fmt"
	"reflect"
)

var RuleRegistry = make(map[string]RuleHandler)

var keyToType = map[string]reflect.Type{
	"ip_rate":        reflect.TypeOf((*IPRateLimitRuleHandler)(nil)),
	"ip_restriction": reflect.TypeOf((*IPRestrictionRuleHandler)(nil)),
	"path_rate":      reflect.TypeOf((*PathRateLimitRuleHandler)(nil)),
}

func GetRuleHandler(keyName string) (*RuleHandler, bool) {
	t, exists := keyToType[keyName]
	if !exists {
		fmt.Println("RuleHandler: ", keyName, "is not exists")
		return nil, false
	}

	handlerPtr := reflect.New(t.Elem()).Interface()
	handler, ok := handlerPtr.(RuleHandler)
	if !ok {
		fmt.Println("RuleHandler: ", keyName, "cannot be converted to RuleHandler", t)
		return nil, false
	}

	handler.Init()

	return &handler, true
}
