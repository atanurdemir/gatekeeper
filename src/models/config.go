package models

type AppConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	Gates     []Gate `mapstructure:"gates"`
}

type Gate struct {
	Path   string `mapstructure:"path"`
	Method string `mapstructure:"method"`
	Rules  []Rule `mapstructure:"rules"`
}

type Rule struct {
	Name   string      `mapstructure:"name"`
	Config interface{} `mapstructure:"config"`
}
