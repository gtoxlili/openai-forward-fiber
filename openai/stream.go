package openai

import (
	json "github.com/bytedance/sonic"
	"io"
	"openai-forward-fiber/common/pool"
	"openai-forward-fiber/entity"
)

// GetStreamRes 获取流的响应内容
func GetStreamRes(reader io.Reader, onClose ...func()) string {
	sb := pool.GetBuffer()
	defer pool.PutBuffer(sb)
	vo := &entity.OpenaiStreamVO{
		Choices: make([]struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		}, 1)}
	decoder := json.ConfigDefault.NewDecoder(reader)
	for {
		vo.Choices[0].Delta.Content = ""
		if err := decoder.Decode(vo); err != nil {
			break
		}
		sb.WriteString(vo.Choices[0].Delta.Content)
	}
	if len(onClose) > 0 {
		for _, f := range onClose {
			f()
		}
	}
	return sb.String()
}
