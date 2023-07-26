package openai

import (
	"github.com/cespare/xxhash/v2"
	"openai-forward-fiber/config"
)

func LoadBalance(apiKey []byte) string {
	hash := xxhash.Sum64(apiKey)
	return config.RootToken[hash%uint64(len(config.RootToken))]
}
