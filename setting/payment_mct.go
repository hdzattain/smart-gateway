package setting

var (
	MCTPayEnabled     bool
	MCTPayMerchantID  string
	MCTPaySecretKey   string
	MCTPayCheckoutURL string = "https://mct.com.sg/chn/mctpay/"
	MCTPayWebhookURL  string
	MCTPayUnitPrice   float64 = 1.0
	MCTPayMinTopUp    int     = 1
)
