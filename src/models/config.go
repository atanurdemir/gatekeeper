package models

type AppConfig struct {
	Gates []Gate `yaml:"gates"`
}

type Gate struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Rules  []Rule `yaml:"rules"`
}

type Rule struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
}
