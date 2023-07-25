package openai

import (
	"bytes"
	json "github.com/bytedance/sonic"
	"io"
	"openai-forward-fiber/entity"
)

// GetStreamRes 获取流的响应内容
func GetStreamRes(reader io.Reader) string {
	sb := &bytes.Buffer{}
	decoder := json.ConfigDefault.NewDecoder(reader)
	for {
		var vo entity.OpenaiStreamVO
		if err := decoder.Decode(&vo); err != nil {
			break
		}
		sb.WriteString(vo.Choices[0].Delta.Content)
	}
	return sb.String()
}
