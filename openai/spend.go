package openai

import (
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"openai-forward-fiber/db"
	"openai-forward-fiber/entity"
	"strings"
)

var (
	// 模型对应的 token 消耗
	ratioMap = map[string]map[bool]float64{
		"gpt-3.5": {
			true:  0.0015,
			false: 0.002,
		},
		"gpt-4": {
			true:  0.03,
			false: 0.06,
		},
		"gpt-4-32k": {
			true:  0.06,
			false: 0.12,
		},
		"text-embedding": {
			true:  0.0004,
			false: 0.0004,
		},
	}
)

func SpendHandler(auth, model string, tokens int, isPrompt bool) {
	// completion_tokens
	log.Info().Str("AUTH", auth).Str("MODEL", model).Int("TOKENS", tokens).Bool("[IS COMPLETION]", !isPrompt).Msg("[Spend]")
	// 如果 auth 为 sk- 开头，则不计算消耗
	if strings.HasPrefix(auth, "Bearer sk-") {
		return
	}
	var info entity.UserInfo
	var err error
	infoByte, err := db.Db().Get("info:" + auth)
	err = json.Unmarshal(infoByte, &info)
	info.UsedTokens += getRatio(model, isPrompt) * float64(tokens) / 1000
	infoByte, err = json.Marshal(info)
	err = db.Db().Set("info:"+auth, infoByte, 0)
	if err != nil {
		log.Warn().Err(err).Msgf("记录 Tokens 损耗失败: %s", auth)
	}
}

func getRatio(model string, isPrompt bool) float64 {
	if v, ok := ratioMap[model]; ok {
		if ratio, ok := v[isPrompt]; ok {
			return ratio
		}
	}
	if isPrompt {
		return 0.003
	}
	return 0.004
}
