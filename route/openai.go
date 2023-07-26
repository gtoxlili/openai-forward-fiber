package route

import (
	"bytes"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io"
	"openai-forward-fiber/common/pool"
	"openai-forward-fiber/common/signBuffer"
	"openai-forward-fiber/config"
	"openai-forward-fiber/entity"
	"openai-forward-fiber/openai"
	"strings"
)

func Openai(r fiber.Router) {
	if config.RejectionOpenaiApiKey {
		r.Use(openaiReject)
	}
	r.Use(openaiAllow)
	r.Use(limiter.New(limiter.Config{
		SkipFailedRequests: true,
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Get("Authorization"), "Bearer sk-")
		},
		Max:               config.LimiterMax,
		Expiration:        config.LimiterExpiration,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("Authorization")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"message": "Too many requests, please try again later.",
					"code":    "too_many_requests",
				},
			})
		},
		Storage: db,
	}))
	r.Use(openaiCheckToken)
	r.Post("/", openaiForward)
}

func openaiReject(c *fiber.Ctx) error {
	apiKey := c.Get("Authorization")
	if strings.HasPrefix(apiKey, "Bearer sk-") {
		return fmt.Errorf("invalid_token: %w", fmt.Errorf("暂不支持直接使用 OpenAI 的 Secret Key"))
	}
	return c.Next()
}

// 校验 token 是否存在且还有余额的中间件
func openaiCheckToken(c *fiber.Ctx) error {
	apiKey := c.Get("Authorization")
	if apiKey == "" || strings.HasPrefix(apiKey, "Bearer sk-") {
		c.Context().SetUserValue("api_key", apiKey)
		return c.Next()
	}
	if !strings.HasPrefix(apiKey, "Bearer ck-") {
		return fmt.Errorf("invalid_token: %w", fmt.Errorf("token 格式不正确"))
	}
	// 校验 token 是否存在
	infoByte, err := db.Get("info:" + apiKey)
	if err != nil || infoByte == nil {
		if infoByte == nil {
			err = fmt.Errorf("token 不存在")
		}
		return fmt.Errorf("token_not_found: %w", err)
	}
	// 校验 token 是否还有余额
	var info entity.UserInfo
	if err := json.Unmarshal(infoByte, &info); err != nil {
		return fmt.Errorf("unmarshal_token: %w", err)
	}

	// 检测代币
	if info.UsedTokens >= info.TotalTokens {
		return fmt.Errorf("token_not_enough: %w", fmt.Errorf("代币余额不足"))
	}
	// TODO 检测模型是否可用

	c.Context().SetUserValue("api_key", "Bearer "+config.RootToken)
	return c.Next()
}

// 允许通行的请求
func openaiAllow(c *fiber.Ctx) error {
	params := c.Params("+")
	if !openai.IsAllowedRoute(params) {
		return fmt.Errorf("not_allowed_route: %w", fmt.Errorf("不支持的服务: /%s", params))
	}
	return c.Next()
}

// Openai 路由
//
//	@Summary	OpenAI 转发服务
//	@Tags		OpenAI
//	@Router		/openai/{+} [post]
//	@Param		+				path	string	true	"服务名称"
//	@Param		Authorization	header	string	true	"API Key"
//	@Param		Content-Type	header	string	true	"Content-Type"
//	@Param		dto				body	object	true	"请求体"
func openaiForward(c *fiber.Ctx) error {

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	c.Request().CopyTo(req)
	req.SetRequestURI("https://api.openai.com/v1/" + c.Params("+"))
	req.Header.Set("Authorization", c.Context().UserValue("api_key").(string))

	resp := fiber.AcquireResponse()
	// defer fiber.ReleaseResponse(resp)
	if err := agent.Parse(); err != nil {
		return fmt.Errorf("parse_request: %w", err)
	}

	// 根据请求头判断是否需要流式响应
	var dto entity.OpenaiDto
	if err := c.BodyParser(&dto); err != nil {
		return fmt.Errorf("unmarshal_request: %w", err)
	}
	if dto.Stream {
		agent.HostClient.StreamResponseBody = true
	}

	// 设置代理
	if config.ProxyAddr != "" {
		agent.HostClient.Dial = fasthttpproxy.FasthttpHTTPDialer(config.ProxyAddr)
	}

	// 发送请求
	if err := agent.HostClient.Do(req, resp); err != nil {
		return fmt.Errorf("do_request: %w", err)
	}

	resp.CopyTo(c.Response())

	if dto.Stream {
		var bodyStream io.Reader
		// 如果返回的是流 Content-Type: text/event-stream
		if utils.UnsafeString(resp.Header.ContentType()) == "text/event-stream" {

			buf := signBuffer.New(func(b []byte) bool {
				return bytes.Contains(b, []byte("[DONE]"))
			}, config.StreamTimeout, []byte("data:"))

			// 计算代币
			openai.SpendHandler(c.Get("Authorization"), dto.Model, openai.CalculateDtoTokens(dto), true)
			go func(authorization, model string, release func()) {
				openai.SpendHandler(authorization, model, openai.CalculateTokens(openai.GetStreamRes(buf, release)), false)
			}(c.Get("Authorization"), dto.Model, func() {
				fiber.ReleaseResponse(resp)
			})

			bodyStream = io.TeeReader(resp.BodyStream(), buf)
			return c.SendStream(bodyStream)
		} else {
			dst := pool.Get(resp.Header.ContentLength())
			defer pool.Put(dst)
			n, _ := resp.BodyStream().Read(dst)
			defer fiber.ReleaseResponse(resp)

			return c.Send(dst[:n])
		}
	} else {
		defer fiber.ReleaseResponse(resp)
		body, _ := resp.BodyUnbrotli()

		// 计算代币
		var vo entity.OpenaiVO
		if err := json.Unmarshal(body, &vo); err == nil {
			openai.SpendHandler(c.Get("Authorization"), dto.Model, vo.Usage.PromptTokens, true)
			openai.SpendHandler(c.Get("Authorization"), dto.Model, vo.Usage.TotalTokens-vo.Usage.PromptTokens, false)
		}
		return nil
	}
}
