package openai

import "openai-forward-fiber/config"

// IsAllowedRoute router 是否在允许的路由后缀中
func IsAllowedRoute(router string) bool {
	for _, route := range config.AllowedRoutes {
		if route == router {
			return true
		}
	}
	return false
}
