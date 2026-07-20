package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hdzattain/smart-gateway/common"
	"github.com/hdzattain/smart-gateway/logger"
	"github.com/hdzattain/smart-gateway/model"
	"github.com/hdzattain/smart-gateway/service"
	"github.com/hdzattain/smart-gateway/setting"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MCTPayRequest struct {
	Amount int64 `json:"amount"`
}

func isMCTPayTopUpEnabled() bool {
	return setting.MCTPayEnabled && strings.TrimSpace(setting.MCTPayMerchantID) != "" && strings.TrimSpace(setting.MCTPaySecretKey) != ""
}

func getMCTPayMinTopup() int64 {
	minTopup := int64(setting.MCTPayMinTopUp)
	if minTopup < 1 {
		return 1
	}
	return minTopup
}

func getMCTPayMoney(amount int64) float64 {
	unitPrice := setting.MCTPayUnitPrice
	if unitPrice <= 0 {
		unitPrice = 1.0
	}
	return decimal.NewFromInt(amount).Mul(decimal.NewFromFloat(unitPrice)).InexactFloat64()
}

func RequestMCTPayAmount(c *gin.Context) {
	var req MCTPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}
	if req.Amount < getMCTPayMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getMCTPayMinTopup())})
		return
	}
	payMoney := getMCTPayMoney(req.Amount)
	if payMoney < 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

func RequestMCTPay(c *gin.Context) {
	if !isMCTPayTopUpEnabled() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "MCT Pay 未启用或配置不完整"})
		return
	}

	var req MCTPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}
	if req.Amount < getMCTPayMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getMCTPayMinTopup())})
		return
	}

	id := c.GetInt("id")
	payMoney := getMCTPayMoney(req.Amount)
	if payMoney < 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	tradeNo := fmt.Sprintf("MCT%dNO%s%d", id, common.GetRandomString(6), time.Now().Unix())
	topUp := &model.TopUp{
		UserId:          id,
		Amount:          req.Amount,
		Money:           payMoney,
		TradeNo:         tradeNo,
		PaymentMethod:   model.PaymentMethodMCTPay,
		PaymentProvider: model.PaymentProviderMCTPay,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := topUp.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("MCT Pay 创建充值订单失败 user_id=%d trade_no=%s amount=%d error=%q", id, tradeNo, req.Amount, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	checkoutURL := strings.TrimSpace(setting.MCTPayCheckoutURL)
	if checkoutURL == "" {
		checkoutURL = "https://mct.com.sg/chn/mctpay/"
	}
	notifyURL := strings.TrimSpace(setting.MCTPayWebhookURL)
	if notifyURL == "" {
		notifyURL = serviceCallbackAddress() + "/api/mct-pay/webhook"
	}
	returnURL := paymentReturnPath("/console/topup")
	params := map[string]string{
		"merchant_id": setting.MCTPayMerchantID,
		"trade_no":    tradeNo,
		"amount":      strconv.FormatFloat(payMoney, 'f', 2, 64),
		"currency":    "SGD",
		"subject":     fmt.Sprintf("Smart Gateway top-up %d", req.Amount),
		"notify_url":  notifyURL,
		"return_url":  returnURL,
		"timestamp":   strconv.FormatInt(time.Now().Unix(), 10),
	}
	params["sign"] = signMCTPayParams(params, setting.MCTPaySecretKey)

	payURL, err := appendQueryParams(checkoutURL, params)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("MCT Pay 支付地址生成失败 user_id=%d trade_no=%s error=%q", id, tradeNo, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("MCT Pay 充值订单创建成功 user_id=%d trade_no=%s amount=%d money=%.2f", id, tradeNo, req.Amount, payMoney))
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": gin.H{"pay_link": payURL}})
}

func MCTPayWebhook(c *gin.Context) {
	ctx := c.Request.Context()
	if !isMCTPayTopUpEnabled() {
		logger.LogWarn(ctx, fmt.Sprintf("MCT Pay webhook 被拒绝 reason=webhook_disabled path=%q client_ip=%s", c.Request.RequestURI, c.ClientIP()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		logger.LogWarn(ctx, fmt.Sprintf("MCT Pay webhook 表单解析失败 path=%q client_ip=%s error=%q", c.Request.RequestURI, c.ClientIP(), err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	params := make(map[string]string, len(c.Request.Form))
	for key, values := range c.Request.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	if !verifyMCTPaySignature(params, setting.MCTPaySecretKey) {
		logger.LogWarn(ctx, fmt.Sprintf("MCT Pay webhook 验签失败 path=%q client_ip=%s", c.Request.RequestURI, c.ClientIP()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tradeNo := firstNonEmpty(params["trade_no"], params["out_trade_no"], params["order_no"], params["order_id"])
	status := strings.ToLower(firstNonEmpty(params["status"], params["trade_status"], params["payment_status"]))
	if tradeNo == "" {
		logger.LogWarn(ctx, fmt.Sprintf("MCT Pay webhook 缺少订单号 client_ip=%s", c.ClientIP()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if status != "paid" && status != "success" && status != "succeeded" && status != "completed" {
		logger.LogInfo(ctx, fmt.Sprintf("MCT Pay webhook 忽略未完成状态 trade_no=%s status=%s client_ip=%s", tradeNo, status, c.ClientIP()))
		c.Status(http.StatusOK)
		return
	}

	if err := model.RechargeGeneric(tradeNo, model.PaymentProviderMCTPay, c.ClientIP()); err != nil {
		logger.LogWarn(ctx, fmt.Sprintf("MCT Pay webhook 充值处理失败 trade_no=%s client_ip=%s error=%q", tradeNo, c.ClientIP(), err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	logger.LogInfo(ctx, fmt.Sprintf("MCT Pay webhook 充值成功 trade_no=%s client_ip=%s", tradeNo, c.ClientIP()))
	c.String(http.StatusOK, "success")
}

func signMCTPayParams(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for key, value := range params {
		if key == "sign" || value == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, key+"="+params[key])
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strings.Join(pairs, "&")))
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyMCTPaySignature(params map[string]string, secret string) bool {
	provided := strings.ToLower(strings.TrimSpace(params["sign"]))
	if provided == "" {
		provided = strings.ToLower(strings.TrimSpace(params["signature"]))
	}
	if provided == "" {
		return false
	}
	expected := signMCTPayParams(params, secret)
	return hmac.Equal([]byte(provided), []byte(expected))
}

func appendQueryParams(rawURL string, params map[string]string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func serviceCallbackAddress() string {
	return strings.TrimRight(service.GetCallbackAddress(), "/")
}
