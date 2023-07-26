package config

import "time"

type Config struct {
	Origins               string        `yaml:"origins"`
	ApiVersion            string        `yaml:"api_version"`
	LimiterMax            int           `yaml:"limiter_max"`
	LimiterExpiration     time.Duration `yaml:"limiter_expiration"`
	RootToken             []string      `yaml:"root_token"`
	AdminToken            string        `yaml:"admin_token"`
	StreamTimeout         time.Duration `yaml:"stream_timeout"`
	Addr                  string        `yaml:"addr"`
	RejectionOpenaiApiKey bool          `yaml:"rejection_openai_api_key"`
	AllowedRoutes         []string      `yaml:"allowed_routes"`
	ProxyAddr             string        `yaml:"proxy_addr"`
}

var (
	Origins               string
	ApiVersion            string
	LimiterMax            int
	LimiterExpiration     time.Duration
	RootToken             []string
	AdminToken            string
	StreamTimeout         time.Duration
	Addr                  string
	RejectionOpenaiApiKey bool
	AllowedRoutes         []string
	ProxyAddr             string
)
