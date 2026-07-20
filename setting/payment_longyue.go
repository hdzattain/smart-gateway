package setting

var (
	LongyueEnabled   bool
	LongyueAppId     string
	LongyueSecretKey string
	LongyueApiBase   string // API基础地址，如 https://xxx.com
	LongyueUnitPrice float64 = 1.0
	LongyueMinTopUp  int     = 1
	LongyueCurrency  string  = "USD"
)
