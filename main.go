package main

import (
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"openai-forward-fiber/config"
	_ "openai-forward-fiber/log"
	"openai-forward-fiber/middleware"
	"openai-forward-fiber/route"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/swagger"
	_ "openai-forward-fiber/docs"
)

//	@title			CarPaint AI
//	@description	CarPaint AI API
//	@schemes		http https

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Admin-Token

//	@license.name	GPLv3
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.en.html

func main() {
	app := fiber.New(fiber.Config{
		AppName:     "CarPaint AI",
		Prefork:     true,
		JSONDecoder: json.Unmarshal,
		JSONEncoder: json.Marshal,
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	api := app.Group(fmt.Sprintf("/%s", config.ApiVersion),
		cors.New(cors.Config{
			AllowOrigins:     config.Origins,
			AllowHeaders:     "Content-Type,Authorization,Admin-Token",
			MaxAge:           300,
			AllowCredentials: true,
		}),
		logger.New(
			logger.Config{
				TimeFormat: "2006-01-02 15:04:05",
			}),
		recover.New(),
		middleware.Error,
	)
	route.Openai(api.Group("/openai/+"))
	route.User(api.Group("/user"))
	app.All("/forward", proxy.Forward("http://127.0.0.1:3000/v1/openai/chat/completions"))
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sign
		log.Info().Msg("application gracefully shutdown")
		_ = app.Shutdown()
	}()

	if err := app.Listen(config.Addr); err != nil {
		log.Fatal().Msgf("app error: %s", err.Error())
	}
}
