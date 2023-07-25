package openai

import (
	"bytes"
	"github.com/tiktoken-go/tokenizer"
	"openai-forward-fiber/entity"
)

var (
	enc, _ = tokenizer.Get(tokenizer.Cl100kBase)
)

func convertPrompt(dto entity.OpenaiDto) string {
	messages := dto.Messages
	result := &bytes.Buffer{}
	for _, message := range messages {
		result.WriteString("<|im_start|>")
		result.WriteString(message.Role)
		result.WriteString("\n")
		result.WriteString(message.Content)
		result.WriteString("<|im_end|>\n")
	}
	result.WriteString("<|im_start|>assistant")
	return result.String()
}

func CalculateDtoTokens(dto entity.OpenaiDto) int {
	prompt := convertPrompt(dto)
	ids, _, _ := enc.Encode(prompt)
	return len(ids)
}

func CalculateTokens(str string) int {
	ids, _, _ := enc.Encode(str)
	return len(ids)
}
