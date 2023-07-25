package entity

type OpenaiDto struct {
	Model    string `json:"model"`
	Stream   bool   `json:"stream"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type OpenaiStreamVO struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type OpenaiVO struct {
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

type UserInfo struct {
	// 可用模型
	Models      []string `json:"models"`
	TotalTokens float64  `json:"total_tokens"`
	UsedTokens  float64  `json:"used_tokens"`
	// 充值记录
	RechargeIdRecords []string `json:"recharge_id_records"`
}

type RechargeDto struct {
	TargetKey  string  `json:"target_key" validate:"required,startswith=ck-,len=51"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	RechargeId string  `json:"recharge_id" validate:"required"`
}
