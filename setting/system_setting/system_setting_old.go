package system_setting

var ServerAddress = "https://smart-gateway.shop"
var WorkerUrl = ""
var WorkerValidKey = ""
var WorkerAllowHttpImageRequestEnabled = false

func EnableWorker() bool {
	return WorkerUrl != ""
}
