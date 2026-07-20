package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hdzattain/smart-gateway/common"
	"github.com/hdzattain/smart-gateway/logger"
	"github.com/hdzattain/smart-gateway/model"
	"github.com/hdzattain/smart-gateway/setting"
	"github.com/thanhpk/randstr"
)

type SubscriptionLongyuePayRequest struct {
	PlanId int `json:"plan_id"`
}

func SubscriptionRequestLongyuePay(c *gin.Context) {
	if !requirePaymentCompliance(c) {
		return
	}

	var req SubscriptionLongyuePayRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PlanId <= 0 {
		common.ApiErrorMsg(c, "参数错误")
		return
	}

	plan, err := model.GetSubscriptionPlanById(req.PlanId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if !plan.Enabled {
		common.ApiErrorMsg(c, "套餐未启用")
		return
	}
	if plan.LongyueProductId == "" {
		common.ApiErrorMsg(c, "该套餐未配置龙跃支付产品ID")
		return
	}
	if strings.TrimSpace(setting.LongyueAppId) == "" || strings.TrimSpace(setting.LongyueSecretKey) == "" {
		common.ApiErrorMsg(c, "龙跃支付未配置或密钥无效")
		return
	}

	userId := c.GetInt("id")
	user, err := model.GetUserById(userId, false)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if user == nil {
		common.ApiErrorMsg(c, "用户不存在")
		return
	}

	if plan.MaxPurchasePerUser > 0 {
		count, err := model.CountUserSubscriptionsByPlan(userId, plan.Id)
		if err != nil {
			common.ApiError(c, err)
			return
		}
		if count >= int64(plan.MaxPurchasePerUser) {
			common.ApiErrorMsg(c, "已达到该套餐购买上限")
			return
		}
	}

	reference := fmt.Sprintf("sub-ly-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
	referenceId := "sub_ly_ref_" + common.Sha1([]byte(reference))

	// 调用龙跃CNP系统创建支付订单
	redirectUrl, err := callLongyueSubmit(referenceId, plan.PriceAmount, user, c)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("龙跃支付 订阅支付链接创建失败 trade_no=%s plan_id=%d error=%q", referenceId, plan.Id, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	order := &model.SubscriptionOrder{
		UserId:          userId,
		PlanId:          plan.Id,
		Money:           plan.PriceAmount,
		TradeNo:         referenceId,
		PaymentMethod:   model.PaymentMethodLongyue,
		PaymentProvider: model.PaymentProviderLongyue,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := order.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("龙跃支付 创建订阅订单失败 trade_no=%s plan_id=%d error=%q", referenceId, plan.Id, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("龙跃支付 订阅订单创建成功 user_id=%d trade_no=%s plan_id=%d money=%.2f", userId, referenceId, plan.Id, plan.PriceAmount))
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"pay_link": redirectUrl,
		},
	})
}
