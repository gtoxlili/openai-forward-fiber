package route

import (
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jaevor/go-nanoid"
	"openai-forward-fiber/common/validate"
	"openai-forward-fiber/config"
	dbc "openai-forward-fiber/db"
	"openai-forward-fiber/entity"
)

var (
	db           = dbc.Db()
	generator, _ = nanoid.Custom("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 48)
)

func User(r fiber.Router) {

	r.Get("/add", userAdd)
	r.Get("/info", userInfo)

	// 鉴权
	au := r.Use(userAuth)
	au.Put("/recharge", userRecharge)
	au.Delete("/delete/ck-:apikey<len(48)>", userDelete)
}

func userAuth(c *fiber.Ctx) error {
	if c.Get("Admin-Token") == "" || c.Get("Admin-Token") != config.AdminToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{
				"message": "Unauthorized",
			},
		})
	}
	return c.Next()
}

// @Summary	新增 API Key
// @Tags		User
// @Router		/user/add [GET]
// @Produce	json
func userAdd(c *fiber.Ctx) error {
	// 随机生成一个 token
	apiKey := fmt.Sprintf("ck-%s", generator())
	// 保存到数据库
	info := &entity.UserInfo{
		Models:      []string{"*"},
		TotalTokens: 0,
		UsedTokens:  0,
		// 充值记录
		RechargeIdRecords: []string{},
	}
	bys, _ := json.Marshal(info)
	if err := db.Set("info:Bearer "+apiKey, bys, 0); err != nil {
		return fmt.Errorf("db_error: %w", err)
	}
	return c.JSON(fiber.Map{
		"api_key": apiKey,
	})
}

// @Summary	充值
// @Tags		User
// @Router		/user/recharge [PUT]
// @Produce	json
// @Param		dto	body	entity.RechargeDto true "RechargeDto"
// @Security	ApiKeyAuth
// @Success 200 {object} entity.UserInfo
func userRecharge(c *fiber.Ctx) error {
	var dto entity.RechargeDto
	if err := validate.BodyParser(c, &dto); err != nil {
		return fmt.Errorf("invalid_params: %w", err)
	}
	// 校验 token 是否存在
	infoByte, err := db.Get("info:Bearer " + dto.TargetKey)
	if err != nil || infoByte == nil {
		if err == nil {
			err = fmt.Errorf("token not found")
		}
		return fmt.Errorf("db_error: %w", err)
	}
	var info entity.UserInfo
	json.Unmarshal(infoByte, &info)
	// 检测是否已经充值过
	for _, id := range info.RechargeIdRecords {
		if id == dto.RechargeId {
			return fmt.Errorf("recharge_id_exists: %w", fmt.Errorf("重复的充值 ID"))
		}
	}
	info.TotalTokens += dto.Amount
	info.RechargeIdRecords = append(info.RechargeIdRecords, dto.RechargeId)
	bys, _ := json.Marshal(info)
	if err := db.Set("info:Bearer "+dto.TargetKey, bys, 0); err != nil {
		return fmt.Errorf("db_error: %w", err)
	}
	// 记录充值记录
	if err := db.Set("recharge:"+dto.RechargeId, []byte(fmt.Sprintf("%.2f", dto.Amount)), 0); err != nil {
		return fmt.Errorf("db_error: %w", err)
	}
	return c.JSON(info)
}

// @Summary	获取用户信息
// @Tags		User
// @Router		/user/info [GET]
// @Produce	json
// @Param		Authorization	header	string	true	"API Key"
// @Success 200 {object} entity.UserInfo
func userInfo(c *fiber.Ctx) error {
	apiKey := c.Get("Authorization")
	infoByte, err := db.Get("info:" + apiKey)
	if err != nil || infoByte == nil {
		if err == nil {
			err = fmt.Errorf("token not found")
		}
		return fmt.Errorf("db_error: %w", err)
	}
	var info entity.UserInfo
	json.Unmarshal(infoByte, &info)
	return c.JSON(info)
}

// @Summary	删除用户
// @Tags		User
// @Router		/user/delete/{apiKey} [DELETE]
// @Produce	json
// @Param		apiKey	path	string	true	"API Key"
// @Security	ApiKeyAuth
func userDelete(c *fiber.Ctx) error {
	apiKey := c.Params("apikey")
	if err := db.Delete("info:Bearer ck-" + apiKey); err != nil {
		return fmt.Errorf("db_error: %w", err)
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
