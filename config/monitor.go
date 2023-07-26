package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	_ "openai-forward-fiber/log"
	"os"
)

func init() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("启动失败,请检查配置文件")
	}
	defer file.Close()
	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal().Err(err).Msg("解析配置文件失败")
	}

	Origins = config.Origins
	ApiVersion = config.ApiVersion
	LimiterMax = config.LimiterMax
	LimiterExpiration = config.LimiterExpiration
	RootToken = config.RootToken
	AdminToken = config.AdminToken
	StreamTimeout = config.StreamTimeout
	Addr = config.Addr
	RejectionOpenaiApiKey = config.RejectionOpenaiApiKey
	AllowedRoutes = config.AllowedRoutes
	ProxyAddr = config.ProxyAddr

	// 打印配置
	log.Info().Msg("已加载配置文件")
}
