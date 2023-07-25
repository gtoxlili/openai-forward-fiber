package docs

import (
	"openai-forward-fiber/config"
)

func init() {
	SwaggerInfo.Version = config.ApiVersion
	SwaggerInfo.Host = config.Addr
	SwaggerInfo.BasePath = "/" + config.ApiVersion
}
