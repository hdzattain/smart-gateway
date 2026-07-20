package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hdzattain/smart-gateway/common"
	"github.com/hdzattain/smart-gateway/logger"
	"github.com/hdzattain/smart-gateway/model"
	"github.com/hdzattain/smart-gateway/setting"
	"github.com/hdzattain/smart-gateway/setting/operation_setting"
	"github.com/hdzattain/smart-gateway/setting/system_setting"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
)

// LongyuePayRequest represents a payment request for Longyue CNP checkout.
type LongyuePayRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

// generateLongyueSign 生成龙跃CNP系统MD5签名
// 算法: MD5(appId + 按ASCII排序的key1=value1&key2=value2 + secretKey)
func generateLongyueSign(params map[string]string, appId, secretKey string) string {
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

	signStr := appId + strings.Join(pairs, "&") + secretKey
	hash := md5.Sum([]byte(signStr))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

// verifyLongyueSign 验证龙跃CNP系统回调签名
func verifyLongyueSign(params map[string]string, appId, secretKey, sign string) bool {
	if sign == "" {
		return false
	}
	expected := generateLongyueSign(params, appId, secretKey)
	return strings.EqualFold(expected, sign)
}

func getLongyuePayMoney(amount float64, group string) float64 {
	originalAmount := amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		amount = amount / common.QuotaPerUnit
	}
	topupGroupRatio := common.GetTopupGroupRatio(group)
	if topupGroupRatio == 0 {
		topupGroupRatio = 1
	}
	// apply optional preset discount by the original request amount (if configured), default 1.0
	discount := 1.0
	if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(originalAmount)]; ok {
		if ds > 0 {
			discount = ds
		}
	}
	payMoney := amount * setting.LongyueUnitPrice * topupGroupRatio * discount
	return payMoney
}

func getLongyueMinTopup() int64 {
	minTopup := setting.LongyueMinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		minTopup = minTopup * int(common.QuotaPerUnit)
	}
	return int64(minTopup)
}

func RequestLongyueAmount(c *gin.Context) {
	var req LongyuePayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}
	if req.Amount < getLongyueMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getLongyueMinTopup())})
		return
	}
	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}
	payMoney := getLongyuePayMoney(float64(req.Amount), group)
	if payMoney <= 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

func RequestLongyuePay(c *gin.Context) {
	if !isLongyueTopUpEnabled() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "龙跃支付未启用或配置不完整"})
		return
	}

	var req LongyuePayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "参数错误"})
		return
	}
	if req.PaymentMethod != model.PaymentMethodLongyue {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "不支持的支付渠道"})
		return
	}
	if req.Amount < getLongyueMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getLongyueMinTopup())})
		return
	}
	if req.Amount > 10000 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "充值数量不能大于 10000"})
		return
	}

	id := c.GetInt("id")
	user, _ := model.GetUserById(id, false)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "用户不存在"})
		return
	}
	chargedMoney := GetChargedAmount(float64(req.Amount), *user)

	reference := fmt.Sprintf("smart-gateway-ly-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
	referenceId := "ly_ref_" + common.Sha1([]byte(reference))

	// 调用龙跃CNP系统创建支付订单
	redirectUrl, err := callLongyueSubmit(referenceId, chargedMoney, user, c)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("龙跃支付 创建支付订单失败 user_id=%d trade_no=%s amount=%d error=%q", id, referenceId, req.Amount, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	topUp := &model.TopUp{
		UserId:          id,
		Amount:          req.Amount,
		Money:           chargedMoney,
		TradeNo:         referenceId,
		PaymentMethod:   model.PaymentMethodLongyue,
		PaymentProvider: model.PaymentProviderLongyue,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := topUp.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("龙跃支付 创建充值订单失败 user_id=%d trade_no=%s amount=%d error=%q", id, referenceId, req.Amount, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	logger.LogInfo(c.Request.Context(), fmt.Sprintf("龙跃支付 充值订单创建成功 user_id=%d trade_no=%s amount=%d money=%.2f", id, referenceId, req.Amount, chargedMoney))
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"pay_link": redirectUrl,
		},
	})
}

// callLongyueSubmit 调用龙跃CNP系统 /api/SaaSPay/submit 接口创建支付订单
func callLongyueSubmit(referenceId string, amount float64, user *model.User, c *gin.Context) (string, error) {
	apiBase := strings.TrimRight(setting.LongyueApiBase, "/")
	submitURL := apiBase + "/api/SaaSPay/submit"

	// 构造账单地址（使用用户邮箱 + 默认地址）
	email := user.Email
	if email == "" {
		email = "user@example.com"
	}
	billingAddress := map[string]string{
		"firstName": "Customer",
		"lastName":  "User",
		"email":     email,
		"phone":     "N/A",
		"country":   "US",
		"state":     "N/A",
		"city":      "N/A",
		"address":   "N/A",
		"zipCode":   "00000",
	}
	billingAddressJSON, _ := json.Marshal(billingAddress)
	shippingAddressJSON, _ := json.Marshal(billingAddress) // 同billingAddress结构

	// 构造商品信息
	productInfo := []map[string]interface{}{
		{
			"name":     "Smart Gateway Credits",
			"option":   "topup",
			"price":    fmt.Sprintf("%.2f", amount),
			"quantity": 1,
			"url":      system_setting.ServerAddress,
		},
	}
	productInfoJSON, _ := json.Marshal(productInfo)

	// 构造客户端信息
	userAgent := c.GetHeader("User-Agent")
	clientAgent := map[string]string{
		"os":             "Unknown",
		"browser":        "Unknown",
		"language":       c.GetHeader("Accept-Language"),
		"timeZone":       "UTC",
		"screenResolution": "1920x1080",
		"currentUrl":     c.Request.Referer(),
		"userAgent":      userAgent,
	}
	clientAgentJSON, _ := json.Marshal(clientAgent)

	// 构造请求参数
	params := map[string]string{
		"appId":           setting.LongyueAppId,
		"productName":     "Smart Gateway Credits",
		"billingAddress":  string(billingAddressJSON),
		"shippingAddress": string(shippingAddressJSON),
		"amount":          fmt.Sprintf("%.2f", amount),
		"orderSn":         referenceId,
		"notificationUrl": serviceCallbackAddress() + "/api/longyue/webhook",
		"callbackUrl":     paymentReturnPath("/console/topup"),
		"customer_ip":     c.ClientIP(),
		"productUrl":      system_setting.ServerAddress,
		"productInfo":     string(productInfoJSON),
		"client_agent":    string(clientAgentJSON),
		"currency":        setting.LongyueCurrency,
		"card":            "", // 空字符串，跳转模式
		"payType":         "1", // 3D支付
		"paySource":       "0",
	}

	// 生成签名
	params["sign"] = generateLongyueSign(params, setting.LongyueAppId, setting.LongyueSecretKey)

	// 发送 HTTP POST 请求（application/form-data）
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, value)
	}

	req, err := http.NewRequest("POST", submitURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("cnp-merchant-id", setting.LongyueAppId)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求龙跃支付接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			RedirectUrl   string `json:"redirectUrl"`
			PayCode       int    `json:"payCode"`
			PayStatus     int    `json:"payStatus"`
			TransactionId string `json:"transactionId"`
			OrderSn       string `json:"orderSn"`
			Sign          string `json:"sign"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}

	if result.Code != 1 {
		return "", fmt.Errorf("龙跃支付返回错误: code=%d, msg=%s", result.Code, result.Msg)
	}

	// payCode=301 表示需要跳转
	if result.Data.PayCode == 301 && result.Data.RedirectUrl != "" {
		return result.Data.RedirectUrl, nil
	}

	// payCode=100 表示直接成功（一般不会在跳转模式出现）
	if result.Data.PayCode == 100 {
		return result.Data.RedirectUrl, nil
	}

	// 其他情况返回错误
	return "", fmt.Errorf("龙跃支付返回未知状态: payCode=%d, msg=%s", result.Data.PayCode, result.Msg)
}

// LongyueWebhook 处理龙跃CNP系统异步回调通知
func LongyueWebhook(c *gin.Context) {
	ctx := c.Request.Context()
	if !isLongyueWebhookEnabled() {
		logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 webhook 被拒绝 reason=webhook_disabled path=%q client_ip=%s", c.Request.RequestURI, c.ClientIP()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// 读取请求体
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("龙跃支付 webhook 读取请求体失败 path=%q client_ip=%s error=%q", c.Request.RequestURI, c.ClientIP(), err.Error()))
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 webhook 收到请求 path=%q client_ip=%s body=%q", c.Request.RequestURI, c.ClientIP(), string(payload)))

	// 解析回调参数（支持 JSON 和 form-data 两种格式）
	params := make(map[string]string)
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.Unmarshal(payload, &params); err != nil {
			logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 webhook JSON解析失败 path=%q client_ip=%s error=%q", c.Request.RequestURI, c.ClientIP(), err.Error()))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else {
		// form-data 格式
		c.Request.Body = io.NopCloser(strings.NewReader(string(payload)))
		if err := c.Request.ParseForm(); err != nil {
			logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 webhook 表单解析失败 path=%q client_ip=%s error=%q", c.Request.RequestURI, c.ClientIP(), err.Error()))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		for key, values := range c.Request.Form {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}

	// 验证签名
	sign := params["sign"]
	if !verifyLongyueSign(params, setting.LongyueAppId, setting.LongyueSecretKey, sign) {
		logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 webhook 验签失败 path=%q client_ip=%s", c.Request.RequestURI, c.ClientIP()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 提取订单号
	orderSn := params["orderSn"]
	if orderSn == "" {
		logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 webhook 缺少订单号 client_ip=%s", c.ClientIP()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	payStatus := params["payStatus"]
	payCode := params["payCode"]

	logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 webhook 验签成功 trade_no=%s payStatus=%s payCode=%s client_ip=%s", orderSn, payStatus, payCode, c.ClientIP()))

	// 加锁处理
	LockOrder(orderSn)
	defer UnlockOrder(orderSn)

	// 根据 payStatus/payCode 处理
	// payCode=100 或 payStatus="1" 表示成功
	// payCode=101 或 payStatus="2" 表示失败
	if payCode == "100" || payStatus == "1" {
		// 支付成功
		// 先尝试处理订阅订单
		if err := model.CompleteSubscriptionOrder(orderSn, string(payload), model.PaymentProviderLongyue, ""); err == nil {
			logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 订阅订单处理成功 trade_no=%s client_ip=%s", orderSn, c.ClientIP()))
			c.String(http.StatusOK, "SUCCESS")
			return
		} else if !errors.Is(err, model.ErrSubscriptionOrderNotFound) {
			logger.LogError(ctx, fmt.Sprintf("龙跃支付 订阅订单处理失败 trade_no=%s client_ip=%s error=%q", orderSn, c.ClientIP(), err.Error()))
			c.String(http.StatusOK, "SUCCESS")
			return
		}

		// 降级处理充值订单
		if err := model.RechargeGeneric(orderSn, model.PaymentProviderLongyue, c.ClientIP()); err != nil {
			logger.LogError(ctx, fmt.Sprintf("龙跃支付 充值处理失败 trade_no=%s client_ip=%s error=%q", orderSn, c.ClientIP(), err.Error()))
			c.String(http.StatusOK, "SUCCESS")
			return
		}
		logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 充值成功 trade_no=%s client_ip=%s", orderSn, c.ClientIP()))
	} else if payCode == "101" || payStatus == "2" {
		// 支付失败
		if err := model.UpdatePendingTopUpStatus(orderSn, model.PaymentProviderLongyue, common.TopUpStatusFailed); err != nil {
			logger.LogWarn(ctx, fmt.Sprintf("龙跃支付 更新订单失败状态失败 trade_no=%s client_ip=%s error=%q", orderSn, c.ClientIP(), err.Error()))
		} else {
			logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 订单已标记为失败 trade_no=%s client_ip=%s", orderSn, c.ClientIP()))
		}
	} else {
		// 其他状态（如 payCode=102 支付中, payCode=103 订单关闭）
		logger.LogInfo(ctx, fmt.Sprintf("龙跃支付 webhook 忽略状态 trade_no=%s payStatus=%s payCode=%s client_ip=%s", orderSn, payStatus, payCode, c.ClientIP()))
	}

	c.String(http.StatusOK, "SUCCESS")
}
