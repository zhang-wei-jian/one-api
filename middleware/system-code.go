package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/model"
	logModel "github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/relay/controller"
	"github.com/songquanpeng/one-api/relay/meta"
)

func SystemCode() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 处理机器码

		systemCodeHeader := c.Request.Header.Get("systemCode")

		// tokenId := c.GetInt(ctxkey.TokenId)
		// channelType := c.GetInt(ctxkey.Channel)

		SystemCode := c.GetString(ctxkey.SystemCode)

		channelId := c.GetInt(ctxkey.ChannelId)
		userId := c.GetInt(ctxkey.Id)
		group := c.GetString(ctxkey.Group)
		tokenName := c.GetString(ctxkey.TokenName)

		meta := meta.GetByContext(c)

		textRequest, err := controller.GetAndValidateTextRequest(c, meta.Mode)
		if err != nil {
			fmt.Print(err)
		}

		// model.RecordConsumeLog(ctx, meta.UserId, meta.ChannelId, promptTokens, completionTokens, textRequest.Model, meta.TokenName, quota, logContent)
		// model.UpdateUserUsedQuotaAndRequestCount(meta.UserId, quota)
		// model.UpdateChannelUsedQuota(meta.ChannelId, quota)
		// 构建要返回的 JSON 数据
		// responseData := gin.H{
		// 	"tokenId":     tokenId,
		// 	"channelType": channelType,
		// 	"channelId":   channelId,
		// 	"userId":      userId,
		// 	"group":       group,
		// 	"tokenName":   tokenName,
		// }
		// system := c.GetString("system")
		key := c.GetString("key")
		token, err := model.ValidateUserToken(key)
		if err != nil {
			abortWithMessage(c, http.StatusUnauthorized, err.Error())
			return
		}

		// switch system {
		// case "set":
		// 	logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, "注册了机器码")
		// case "err":
		// 	message := fmt.Sprintf("数据库机器码是 %s ,而用户是 %s", SystemCode, systemCodeHeader)
		// 	logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, message)
		// 	abortWithMessage(c, http.StatusForbidden, "机器码错误")
		// 	return

		// }

		if token.SystemCode == "" {
			// 数据库没有存储机器码，则存储
			token.SystemCode = systemCodeHeader
			err := token.Save()
			logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, "首次注册使用了机器码")
			// 错误之后中间件处理
			if err != nil {
				abortWithMessage(c, http.StatusForbidden, "数据库存储机器码错误")
				// return
			}

		} else {
			// 数据库存储了机器码被使用了，存储在了 token.SystemCode看看是否和请求的 systemCode 一致
			if token.SystemCode != systemCodeHeader {
				message := fmt.Sprintf("数据库机器码是 %s ,而用户是 %s", SystemCode, systemCodeHeader)
				logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, message)
				abortWithMessage(c, http.StatusForbidden, "机器码错误")
				return
			}
		}

		message := fmt.Sprintf("%s 请求了一次", group)
		// 保存日志

		logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, message)

		// if SystemCode == "" {
		// 	// 数据库没有存储机器码，则存储

		// 	logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, "使用了token存储了机器码")

		// } else {
		// 	// 数据库存储了机器码被使用了，存储在了 token.SystemCode看看是否和请求的 systemCode 一致

		// 	if SystemCode != systemCodeHeader {
		// 		message := fmt.Sprintf("数据库机器码%s 而用户是 %s", SystemCode, systemCodeHeader)
		// 		logModel.RecordConsumeLog(c, userId, channelId, int(0), 0, textRequest.Model, tokenName, 0, message)

		// 		abortWithMessage(c, http.StatusForbidden, "机器码错误")
		// 		return
		// 	}
		// }

		c.Next()
	}
}
